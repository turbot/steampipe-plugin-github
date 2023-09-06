package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubMyTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_team",
		Description: "GitHub Teams that you belong to in your organization. GitHub Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyTeamList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "slug"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTeamGet,
		},
		Columns: gitHubTeamColumns(),
	}
}

func tableGitHubMyTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			Organizations struct {
				PageInfo models.PageInfo
				Nodes    []struct {
					Login string
					Teams struct {
						PageInfo models.PageInfo
						Nodes    []models.TeamWithCounts
					} `graphql:"teams(first: $pageSize, after: $cursor)"`
				}
			} `graphql:"organizations(first: $orgPageSize, after: $orgCursor)"`
		}
	}

	orgPageSize := 10 // Note: most users will be in <10 orgs, so this keeps node count down
	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"orgPageSize": githubv4.Int(orgPageSize),
		"orgCursor":   (*githubv4.String)(nil),
		"pageSize":    githubv4.Int(pageSize),
		"cursor":      (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	var teams []models.TeamWithCounts
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_team", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_team", "api_error", err)
			return nil, err
		}

		for _, org := range query.Viewer.Organizations.Nodes {
			plugin.Logger(ctx).Debug("github_my_team", "org", org.Login)
			teams = append(teams, org.Teams.Nodes...)
			if len(teams) >= int(d.RowsRemaining(ctx)) {
				break
			}
			if org.Teams.PageInfo.HasNextPage {
				ts, err := getAdditionalTeams(ctx, client, org.Login, org.Teams.PageInfo.EndCursor)
				if err != nil {
					plugin.Logger(ctx).Error("github_my_team", "api_error", err)
					return nil, err
				}
				teams = append(teams, ts...)
			}
			if len(teams) >= int(d.RowsRemaining(ctx)) {
				break
			}
		}

		for _, team := range teams {
			d.StreamListItem(ctx, team)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.Organizations.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.Organizations.PageInfo.EndCursor)
	}

	return nil, nil
}

func getAdditionalTeams(ctx context.Context, client *githubv4.Client, org string, initialCursor githubv4.String) ([]models.TeamWithCounts, error) {
	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Teams struct {
				PageInfo models.PageInfo
				Nodes    []models.TeamWithCounts
			} `graphql:"teams(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"pageSize": githubv4.Int(100),
		"cursor":   githubv4.NewString(initialCursor),
		"login":    githubv4.String(org),
	}

	var ts []models.TeamWithCounts
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_team", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_team", "api_error", err)
			return nil, err
		}

		ts = append(ts, query.Organization.Teams.Nodes...)

		if !query.Organization.Teams.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Teams.PageInfo.EndCursor)
	}

	return ts, nil
}
