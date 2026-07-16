package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubOrganizationRuleset() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_ruleset",
		Description: "Retrieve the rulesets of a specified GitHub organization.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubOrganizationRulesetList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
			},
		},
		Columns: commonColumns(gitHubOrganizationRulesetColumns()),
	}
}

func gitHubOrganizationRulesetColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The organization login name."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the ruleset."},
		{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the ruleset."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertRulesetTimestamp), Description: "The date and time when the ruleset was created."},
		{Name: "database_id", Type: proto.ColumnType_INT, Description: "The database ID of the ruleset."},
		{Name: "enforcement", Type: proto.ColumnType_STRING, Description: "The enforcement level of the ruleset."},
		{Name: "rules", Type: proto.ColumnType_JSON, Description: "The list of rules in the ruleset."},
		{Name: "bypass_actors", Type: proto.ColumnType_JSON, Description: "The list of actors who can bypass the ruleset."},
		{Name: "conditions", Type: proto.ColumnType_JSON, Description: "The conditions under which the ruleset applies."},
	}
}

func tableGitHubOrganizationRulesetList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Rulesets struct {
				PageInfo struct {
					HasNextPage bool
					EndCursor   githubv4.String
				}
				Edges []struct {
					Node struct {
						CreatedAt   githubv4.DateTime
						DatabaseID  int
						Enforcement string
						Name        string
						ID          string
						Rules       struct {
							PageInfo models.PageInfo
							Edges    []struct {
								Node models.Rule
							}
						} `graphql:"rules(first: $rulePageSize, after: $ruleCursor)"`
						BypassActors struct {
							PageInfo models.PageInfo
							Edges    []struct {
								Node models.BypassActor
							}
						} `graphql:"bypassActors(first: $bypassActorPageSize, after: $bypassActorCursor)"`
						Conditions models.Conditions
					}
				}
			} `graphql:"rulesets(first: $rulesetPageSize, after: $rulesetCursor)"`
		} `graphql:"organization(login: $org)"`
	}

	rulesetPageSize := adjustPageSize(100, d.QueryContext.Limit)
	rulePageSize := 100
	bypassActorPageSize := 100
	org := d.EqualsQuals["organization"].GetStringValue()

	variables := map[string]interface{}{
		"org":                 githubv4.String(org),
		"rulesetPageSize":     githubv4.Int(rulesetPageSize),
		"rulesetCursor":       (*githubv4.String)(nil),
		"rulePageSize":        githubv4.Int(rulePageSize),
		"ruleCursor":          (*githubv4.String)(nil),
		"bypassActorPageSize": githubv4.Int(bypassActorPageSize),
		"bypassActorCursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_organization_ruleset", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_ruleset", "api_error", err)
			return nil, err
		}

		for _, edge := range query.Organization.Rulesets.Edges {
			var rules []models.Rule
			for _, rule := range edge.Node.Rules.Edges {
				rules = append(rules, rule.Node)
			}
			if edge.Node.Rules.PageInfo.HasNextPage {
				additionalRules := getAdditionalOrgRules(ctx, client, org, edge.Node.DatabaseID, "")
				rules = append(rules, additionalRules...)
			}

			var bypassActors []models.BypassActor
			for _, actor := range edge.Node.BypassActors.Edges {
				bypassActors = append(bypassActors, actor.Node)
			}
			if edge.Node.BypassActors.PageInfo.HasNextPage {
				additionalBypassActors := getAdditionalOrgBypassActors(ctx, client, org, edge.Node.DatabaseID, "")
				bypassActors = append(bypassActors, additionalBypassActors...)
			}

			ruleset := models.Ruleset{
				CreatedAt:    edge.Node.CreatedAt.String(),
				DatabaseID:   edge.Node.DatabaseID,
				Enforcement:  edge.Node.Enforcement,
				Name:         edge.Node.Name,
				ID:           edge.Node.ID,
				Rules:        rules,
				BypassActors: bypassActors,
				Conditions:   edge.Node.Conditions,
			}

			d.StreamListItem(ctx, ruleset)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Rulesets.PageInfo.HasNextPage {
			break
		}
		variables["rulesetCursor"] = githubv4.NewString(query.Organization.Rulesets.PageInfo.EndCursor)
	}

	return nil, nil
}

func getAdditionalOrgRules(ctx context.Context, client *githubv4.Client, org string, databaseID int, initialCursor githubv4.String) []models.Rule {
	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Ruleset struct {
				Rules struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Edges []struct {
						Node models.Rule
					}
				} `graphql:"rules(first: $pageSize, after: $cursor)"`
			} `graphql:"ruleset(databaseId: $databaseID)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"org":        githubv4.String(org),
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
	}

	var rules []models.Rule
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_organization_ruleset.getAdditionalOrgRules", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_ruleset.getAdditionalOrgRules", "api_error", err)
			return nil
		}

		for _, edge := range query.Organization.Ruleset.Rules.Edges {
			rules = append(rules, edge.Node)
		}

		if !query.Organization.Ruleset.Rules.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Ruleset.Rules.PageInfo.EndCursor)
	}

	return rules
}

func getAdditionalOrgBypassActors(ctx context.Context, client *githubv4.Client, org string, databaseID int, initialCursor githubv4.String) []models.BypassActor {
	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Ruleset struct {
				BypassActors struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
					Edges []struct {
						Node models.BypassActor
					}
				} `graphql:"bypassActors(first: $pageSize, after: $cursor)"`
			} `graphql:"ruleset(databaseId: $databaseID)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"org":        githubv4.String(org),
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
	}

	var bypassActors []models.BypassActor
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_organization_ruleset.getAdditionalOrgBypassActors", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_ruleset.getAdditionalOrgBypassActors", "api_error", err)
			return nil
		}

		for _, edge := range query.Organization.Ruleset.BypassActors.Edges {
			bypassActors = append(bypassActors, edge.Node)
		}

		if !query.Organization.Ruleset.BypassActors.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Ruleset.BypassActors.PageInfo.EndCursor)
	}

	return bypassActors
}
