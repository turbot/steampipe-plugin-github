package github

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v32/github"
	"time"

	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubConnectivity() *plugin.Table {
	return &plugin.Table{
		Name:        "github_connectivity_test",
		Description: "Github connectivity test suite.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubConnectivityTestList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    tableGitHubConnectivityTestGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: pb.ColumnType_STRING, Description: "The name of the connectivity test."},
			{Name: "description", Type: pb.ColumnType_STRING, Description: "The description of the connectivity test."},
			{Name: "status", Type: pb.ColumnType_STRING, Description: "The result of the connectivity test."},
			{Name: "reason", Type: pb.ColumnType_STRING, Description: "The reason for the status of the connectivity test."},
		},
	}
}

//// list ////

func tableGitHubConnectivityTestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	tests := getConnectivityTests(ctx, client)

	for _, test := range tests {
		row, err := timeout(2*time.Second, test.Action)
		if err != nil {
			errorRow := ConnectivityTestResultRow{
				Name:        test.Name,
				Description: test.Description,
				Status:      "error",
				Reason:      err.Error(),
			}
			d.StreamListItem(ctx, errorRow)
		} else {
			d.StreamListItem(ctx, row)
		}
	}

	return nil, nil
}

//// hydrate functions ////

func tableGitHubConnectivityTestGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	tests := getConnectivityTests(ctx, client)
	name := d.KeyColumnQuals["name"].GetStringValue()
	test := tests[name]
	row, err := timeout(2*time.Second, test.Action)
	if err != nil {
		timeoutRow := ConnectivityTestResultRow{
			Name:        test.Name,
			Description: test.Description,
			Status:      "error",
			Reason:      err.Error(),
		}
		return timeoutRow, nil
	} else {
		return row, nil
	}
}

type ConnectivityTestAction func() ConnectivityTestResultRow

type ConnectivityTestInfo struct {
	Name        string
	Description string
	Action      ConnectivityTestAction
}

type ConnectivityTestResultRow struct {
	Name        string
	Description string
	Status      string
	Reason      string
}

// Wrap tests in a timeout
func timeout(duration time.Duration, action ConnectivityTestAction) (*ConnectivityTestResultRow, error) {
	ch := make(chan ConnectivityTestResultRow)
	timeout := time.After(duration)

	go func() {
		row := action()
		ch <- row
	}()

	select {
	case <-timeout:
		return nil, errors.New(fmt.Sprintf("Timed out after %d seconds", duration.Milliseconds()/1000))
	case row := <-ch:
		return &row, nil
	}
}

func getConnectivityTests(ctx context.Context, client *github.Client) map[string]ConnectivityTestInfo {
	// TODO handle timeouts
	connectivityTests := make(map[string]ConnectivityTestInfo)
	connectivityTests["get_authenticated_user"] = ConnectivityTestInfo{
		Name:        "get_authenticated_user",
		Description: "Ensure the provided token can retrieve the authenticated user",
		Action: func() ConnectivityTestResultRow {
			user, _, err := client.Users.Get(ctx, "")
			var reason string
			var status string
			if err != nil {
				status = "error"
				reason = err.Error()
			} else {
				status = "ok"
				reason = fmt.Sprintf("Got user [%s]", *user.Name)
			}
			return ConnectivityTestResultRow{
				Name:        "get_authenticated_user",
				Description: "Ensure the provided token can retrieve the authenticated user",
				Status:      status,
				Reason:      reason,
			}
		},
	}
	connectivityTests["list_organizations"] = ConnectivityTestInfo{
		Name:        "list_organizations",
		Description: "Ensure the provided token can list organizations",
		Action: func() ConnectivityTestResultRow {
			orgs, _, err := client.Organizations.List(ctx, "", nil)
			var reason string
			var status string
			if err != nil {
				status = "error"
				reason = err.Error()
			} else {
				status = "ok"
				reason = fmt.Sprintf("Got first page of [%d] organizations", len(orgs))
			}
			return ConnectivityTestResultRow{
				Name:        "list_organizations",
				Description: "Ensure the provided token can list organizations",
				Status:      status,
				Reason:      reason,
			}
		},
	}
	connectivityTests["list_repositories"] = ConnectivityTestInfo{
		Name:        "list_repositories",
		Description: "Ensure the provided token can list repositories",
		Action: func() ConnectivityTestResultRow {
			orgs, _, err := client.Repositories.List(ctx, "", nil)
			var reason string
			var status string
			if err != nil {
				status = "error"
				reason = err.Error()
			} else {
				status = "ok"
				reason = fmt.Sprintf("Got first page of [%d] repositories", len(orgs))
			}
			return ConnectivityTestResultRow{
				Name:        "list_repositories",
				Description: "Ensure the provided token can list repositories",
				Status:      status,
				Reason:      reason,
			}
		},
	}
	return connectivityTests
}
