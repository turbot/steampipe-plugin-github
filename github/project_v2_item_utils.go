package github

import (
	"slices"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// projectV2ItemRow flattens a project item and its content into one row. Content
// columns that don't apply to the item's type are left nil.
type projectV2ItemRow struct {
	Id                  string
	NodeId              string
	Type                githubv4.ProjectV2ItemType
	IsArchived          bool
	CreatedAt           models.NullableTime
	UpdatedAt           models.NullableTime
	Creator             models.Actor
	FieldValues         map[string]interface{}
	ContentId           *int
	ContentNodeId       string
	RepositoryFullName  string
	Number              *int
	Title               string
	Body                *string
	Url                 string
	State               string
	Closed              *bool
	ClosedAt            models.NullableTime
	Author              models.Actor
	Assignees           []models.BasicUser
	AssigneesTotalCount *int
	Labels              map[string]bool
	Milestone           models.Milestone
	ContentCreatedAt    models.NullableTime
	ContentUpdatedAt    models.NullableTime
}

var projectV2ItemContentColumns = []string{
	"content_id", "content_node_id", "repository_full_name", "number", "title", "body", "url", "state", "closed",
	"closed_at", "author", "assignees", "assignees_total_count", "labels", "milestone", "content_created_at", "content_updated_at",
}

func appendProjectV2ItemColumnIncludes(m *map[string]interface{}, cols []string) {
	includeContent := slices.ContainsFunc(cols, func(col string) bool {
		return slices.Contains(projectV2ItemContentColumns, col)
	})
	(*m)["includeItemCreator"] = githubv4.Boolean(slices.Contains(cols, "creator"))
	(*m)["includeItemFieldValues"] = githubv4.Boolean(slices.Contains(cols, "field_values"))
	(*m)["includeItemContent"] = githubv4.Boolean(includeContent)
	(*m)["includeItemBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeItemAuthor"] = githubv4.Boolean(slices.Contains(cols, "author"))
	(*m)["includeItemAssignees"] = githubv4.Boolean(slices.Contains(cols, "assignees") || slices.Contains(cols, "assignees_total_count"))
	(*m)["includeItemLabels"] = githubv4.Boolean(slices.Contains(cols, "labels"))
	(*m)["includeItemMilestone"] = githubv4.Boolean(slices.Contains(cols, "milestone"))
}

// projectV2ItemArchivedStates returns both states unless the query filters on
// is_archived, since the API only returns unarchived items by default.
func projectV2ItemArchivedStates(d *plugin.QueryData) []models.ProjectV2ItemArchivedState {
	states := []models.ProjectV2ItemArchivedState{
		models.ProjectV2ItemArchivedStateArchived,
		models.ProjectV2ItemArchivedStateNotArchived,
	}
	if d.Quals["is_archived"] == nil || len(d.Quals["is_archived"].Quals) == 0 {
		return states
	}
	// Postgres still applies every qual to the results, so the first one is
	// enough to narrow the request.
	q := d.Quals["is_archived"].Quals[0]
	archived := q.Value.GetBoolValue()
	if q.Operator == "<>" {
		archived = !archived
	}
	if archived {
		return states[:1]
	}
	return states[1:]
}

// projectV2ItemQuery builds the items filter for the type and
// repository_full_name quals. REDACTED items have no filter, so they're
// filtered after listing.
func projectV2ItemQuery(d *plugin.QueryData) string {
	var filters []string
	switch d.EqualsQualString("type") {
	case string(githubv4.ProjectV2ItemTypeIssue):
		// is:issue also matches draft issues
		filters = append(filters, "is:issue", "-is:draft")
	case string(githubv4.ProjectV2ItemTypePullRequest):
		filters = append(filters, "is:pr")
	case string(githubv4.ProjectV2ItemTypeDraftIssue):
		filters = append(filters, "is:draft")
	}
	if repo := d.EqualsQualString("repository_full_name"); repo != "" {
		filters = append(filters, "repo:"+repo)
	}
	return strings.Join(filters, " ")
}

func newProjectV2ItemRow(item models.ProjectV2Item) projectV2ItemRow {
	row := projectV2ItemRow{
		Id:          item.Id,
		NodeId:      item.NodeId,
		Type:        item.Type,
		IsArchived:  item.IsArchived,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Creator:     item.Creator,
		FieldValues: projectV2ItemFieldValues(item.FieldValues.Nodes),
	}

	switch item.Type {
	case githubv4.ProjectV2ItemTypeIssue:
		c := item.Content.Issue
		row.ContentId = &c.Id
		row.ContentNodeId = c.NodeId
		row.RepositoryFullName = c.Repository.NameWithOwner
		row.Number = &c.Number
		row.Title = c.Title
		row.Body = &c.Body
		row.Url = c.Url
		row.State = string(c.State)
		row.Closed = &c.Closed
		row.ClosedAt = c.ClosedAt
		row.Author = c.Author
		row.Assignees = c.Assignees.Nodes
		row.AssigneesTotalCount = &c.Assignees.TotalCount
		row.Labels = projectV2ItemLabels(c.Labels)
		row.Milestone = c.Milestone
		row.ContentCreatedAt = c.CreatedAt
		row.ContentUpdatedAt = c.UpdatedAt
	case githubv4.ProjectV2ItemTypePullRequest:
		c := item.Content.PullRequest
		row.ContentId = &c.Id
		row.ContentNodeId = c.NodeId
		row.RepositoryFullName = c.Repository.NameWithOwner
		row.Number = &c.Number
		row.Title = c.Title
		row.Body = &c.Body
		row.Url = c.Url
		row.State = string(c.State)
		row.Closed = &c.Closed
		row.ClosedAt = c.ClosedAt
		row.Author = c.Author
		row.Assignees = c.Assignees.Nodes
		row.AssigneesTotalCount = &c.Assignees.TotalCount
		row.Labels = projectV2ItemLabels(c.Labels)
		row.Milestone = c.Milestone
		row.ContentCreatedAt = c.CreatedAt
		row.ContentUpdatedAt = c.UpdatedAt
	case githubv4.ProjectV2ItemTypeDraftIssue:
		c := item.Content.DraftIssue
		row.ContentNodeId = c.NodeId
		row.Title = c.Title
		row.Body = &c.Body
		row.Author = c.Creator
		row.Assignees = c.Assignees.Nodes
		row.AssigneesTotalCount = &c.Assignees.TotalCount
		row.ContentCreatedAt = c.CreatedAt
		row.ContentUpdatedAt = c.UpdatedAt
	}

	return row
}

func projectV2ItemLabels(labels models.ProjectV2ItemLabels) map[string]bool {
	m := make(map[string]bool, len(labels.Nodes))
	for _, l := range labels.Nodes {
		m[l.Name] = true
	}
	return m
}

// projectV2ItemFieldValues maps custom field names to their values. Built-in
// fields (Title, Assignees, Labels, etc.) are skipped since they have their own
// columns.
func projectV2ItemFieldValues(nodes []models.ProjectV2ItemFieldValue) map[string]interface{} {
	values := map[string]interface{}{}
	for _, v := range nodes {
		switch v.TypeName {
		case "ProjectV2ItemFieldTextValue":
			if v.Text.Field.Common.DataType == "TITLE" {
				continue
			}
			values[v.Text.Field.Common.Name] = v.Text.Text
		case "ProjectV2ItemFieldNumberValue":
			values[v.Number.Field.Common.Name] = v.Number.Number
		case "ProjectV2ItemFieldDateValue":
			values[v.Date.Field.Common.Name] = v.Date.Date
		case "ProjectV2ItemFieldSingleSelectValue":
			values[v.SingleSelect.Field.Common.Name] = v.SingleSelect.Name
		case "ProjectV2ItemFieldMultiSelectValue":
			options := make([]string, 0, len(v.MultiSelect.Options))
			for _, o := range v.MultiSelect.Options {
				options = append(options, o.Name)
			}
			values[v.MultiSelect.Field.Common.Name] = options
		case "ProjectV2ItemFieldIterationValue":
			values[v.Iteration.Field.Common.Name] = map[string]interface{}{
				"title":        v.Iteration.Title,
				"start_date":   v.Iteration.StartDate,
				"duration":     v.Iteration.Duration,
				"iteration_id": v.Iteration.IterationId,
			}
		case "ProjectV2ItemIssueFieldValue":
			ifv := v.IssueField.IssueFieldValue
			if ifv.Number.NumberValue != nil {
				values[v.IssueField.Field.Common.Name] = *ifv.Number.NumberValue
				continue
			}
			for _, s := range []string{ifv.Text.StringValue, ifv.Date.StringValue, ifv.SingleSelect.StringValue, ifv.MultiSelect.StringValue} {
				if s != "" {
					values[v.IssueField.Field.Common.Name] = s
					break
				}
			}
		}
	}
	if len(values) == 0 {
		return nil
	}
	return values
}
