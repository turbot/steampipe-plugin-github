package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRepositoryRuleset() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_ruleset",
		Description: "Retrieve the rulesets of a specified GitHub repository.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubRepositoryRulesetList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
			},
		},
		Columns: gitHubRulesetColumns(),
	}
}

func gitHubRulesetColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the ruleset."},
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

func tableGitHubRepositoryRulesetList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
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
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	rulesetPageSize := adjustPageSize(100, d.QueryContext.Limit)
	rulePageSize := 100
	bypassActorPageSize := 100
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	variables := map[string]interface{}{
		"owner":               githubv4.String(owner),
		"name":                githubv4.String(repo),
		"rulesetPageSize":     githubv4.Int(rulesetPageSize),
		"rulesetCursor":       (*githubv4.String)(nil),
		"rulePageSize":        githubv4.Int(rulePageSize),
		"ruleCursor":          (*githubv4.String)(nil),
		"bypassActorPageSize": githubv4.Int(bypassActorPageSize),
		"bypassActorCursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	var rulesets []models.Ruleset
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_ruleset", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_ruleset", "api_error", err)
			return nil, err
		}

		for _, edge := range query.Repository.Rulesets.Edges {

			// Fetch additional Rules.
			var rules []models.Rule
			for _, rule := range edge.Node.Rules.Edges {
				rules = append(rules, rule.Node)
			}
			if edge.Node.Rules.PageInfo.HasNextPage {
				additionalRules := getAdditionalRules(ctx, d, client, edge.Node.DatabaseID, owner, repo, "")
				rules = append(rules, additionalRules...)
			}

			// Fetch additional ByPassActors.
			var bypassActors []models.BypassActor
			for _, actor := range edge.Node.BypassActors.Edges {
				bypassActors = append(bypassActors, actor.Node)
			}
			if edge.Node.BypassActors.PageInfo.HasNextPage {

				additionalBypassActors := getAdditionalBypassActors(ctx, d, client, owner, repo, edge.Node.DatabaseID, "")
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
			rulesets = append(rulesets, ruleset)
		}

		for _, ruleset := range rulesets {
			d.StreamListItem(ctx, ruleset)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Rulesets.PageInfo.HasNextPage {
			break
		}
		variables["rulesetCursor"] = githubv4.NewString(query.Repository.Rulesets.PageInfo.EndCursor)
	}

	return nil, nil
}

func getAdditionalRules(ctx context.Context, d *plugin.QueryData, client *githubv4.Client, databaseID int, owner string, repo string, initialCursor githubv4.String) []models.Rule {

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
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
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repo),
	}

	var rules []models.Rule
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_ruleset.getAdditionalRules", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_ruleset.getAdditionalRules", "api_error", err)
			return nil
		}

		for _, edge := range query.Repository.Ruleset.Rules.Edges {
			rules = append(rules, edge.Node)
		}

		if !query.Repository.Ruleset.Rules.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Ruleset.Rules.PageInfo.EndCursor)
	}

	return rules
}

func getAdditionalBypassActors(ctx context.Context, d *plugin.QueryData, client *githubv4.Client, owner string, repo string, databaseID int, initialCursor githubv4.String) []models.BypassActor {

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
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
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repo),
		"pageSize":   githubv4.Int(100),
		"cursor":     githubv4.NewString(initialCursor),
		"databaseID": githubv4.Int(databaseID),
	}

	var bypassActors []models.BypassActor
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_ruleset.getAdditionalBypassActors", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_ruleset.getAdditionalBypassActors", "api_error", err)
			return nil
		}

		for _, edge := range query.Repository.Ruleset.BypassActors.Edges {
			bypassActors = append(bypassActors, edge.Node)
		}

		if !query.Repository.Ruleset.BypassActors.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Ruleset.BypassActors.PageInfo.EndCursor)
	}

	return bypassActors
}

//// TRANSFORM FUNCTION

// The timestamp value we are receiving has the layout "2024-06-11 13:18:48 +0000 UTC".
// Our generic timestamp function does not support converting this specific layout to the desired format.
// Additionally, it is not feasible to create a generic function that handles all possible timestamp layouts.
// Therefore, we have opted to implement a specific timestamp conversion function for this table only.
func convertRulesetTimestamp(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	t := d.Value.(string)

	// Parse the timestamp into a time.Time object
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", t)
	if err != nil {
		plugin.Logger(ctx).Error("Error parsing time:", err)
		return nil, err
	}
	// Format the time.Time object to RFC 3339 format
	rfc3339Time := parsedTime.Format(time.RFC3339)

	return rfc3339Time, nil
}
