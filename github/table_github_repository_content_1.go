package github

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubRepositoryContent1() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_content_1",
		Description: "List the content in a repository (list directory, or get file content",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubRepositoryContentList1,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "repository_content_path", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Description: "The full name of the repository (login/repo-name).", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name")},
			{Name: "type", Description: "The file type (directory or file).", Type: proto.ColumnType_STRING},
			{Name: "name", Description: "The file name.", Type: proto.ColumnType_STRING},
			{Name: "oid", Description: "The Git object ID.", Type: proto.ColumnType_STRING},
			{Name: "abbreviated_oid", Description: "An abbreviated version of the Git object ID.", Type: proto.ColumnType_STRING},
			{Name: "repository_content_path", Description: "The requested path in repository search.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_content_path")},
			{Name: "path", Description: "The path of the file.", Type: proto.ColumnType_STRING},
			{Name: "path_raw", Description: "A Base64-encoded representation of the file's path.", Type: proto.ColumnType_STRING},
			{Name: "mode", Description: "The mode of the file.", Type: proto.ColumnType_INT},
			{Name: "size", Description: "The size of the file (in KB).", Type: proto.ColumnType_INT},
			{Name: "line_count", Description: "The number of lines available in the file.", Type: proto.ColumnType_INT},
			{Name: "content", Description: "The decoded file content (if the element is a file).", Type: proto.ColumnType_STRING},
			{Name: "is_generated", Description: "Whether or not this tree entry is generated.", Type: proto.ColumnType_BOOL},
			{Name: "is_binary", Description: "Indicates whether the Blob is binary or text.", Type: proto.ColumnType_BOOL},
			{Name: "commit_url", Description: "Git URL (with SHA) of the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type ContentInfo1 struct {
	Oid            string
	AbbreviatedOid string
	Name           string
	Mode           int
	PathRaw        string
	IsGenerated    bool
	Path           string
	Size           int
	LineCount      int
	Type           string
	Content        string
	CommitUrl      string
	IsBinary       bool
}

// File content query
var query struct {
	RateLimit  models.RateLimit
	Repository struct {
		Object struct {
			Tree struct {
				Oid            githubv4.String
				AbbreviatedOid githubv4.String
				Entries        []struct {
					Name        githubv4.String
					Path        githubv4.String
					Size        githubv4.Int
					LineCount   githubv4.Int
					Mode        githubv4.Int
					PathRaw     githubv4.String
					IsGenerated githubv4.Boolean
					Type        githubv4.String
					Object      struct {
						Blob struct {
							Oid            githubv4.String
							AbbreviatedOid githubv4.String
							Text           githubv4.String
							IsBinary       githubv4.Boolean
							CommitUrl      githubv4.String
						} `graphql:"... on Blob"`
						SubTree struct {
							Entries []struct {
								Name        githubv4.String
								Path        githubv4.String
								Size        githubv4.Int
								LineCount   githubv4.Int
								Mode        githubv4.Int
								PathRaw     githubv4.String
								IsGenerated githubv4.Boolean
								Type        githubv4.String
								Object      struct {
									Blob struct {
										Oid            githubv4.String
										AbbreviatedOid githubv4.String
										Text           githubv4.String
										IsBinary       githubv4.Boolean
										CommitUrl      githubv4.String
									} `graphql:"... on Blob"`
								}
							}
						} `graphql:"... on Tree"`
					}
				}
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// Query to list out only trees
var treeQuery struct {
	oneLevel struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Tree struct {
					Entries []struct {
						Name githubv4.String
						Path githubv4.String
						Type githubv4.String
					}
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	twoLevel struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Tree struct {
					Entries []struct {
						Name   githubv4.String
						Path   githubv4.String
						Type   githubv4.String
						Object struct {
							SubTree struct {
								Entries []struct {
									Name githubv4.String
									Path githubv4.String
									Type githubv4.String
								}
							} `graphql:"... on Tree"`
						}
					}
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	threeLevel struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Tree struct {
					Entries []struct {
						Name   githubv4.String
						Path   githubv4.String
						Type   githubv4.String
						Object struct {
							SubTree struct {
								Entries []struct {
									Name   githubv4.String
									Path   githubv4.String
									Type   githubv4.String
									Object struct {
										HydraTree struct {
											Entries []struct {
												Name githubv4.String
												Path githubv4.String
												Type githubv4.String
											}
										} `graphql:"... on Tree"`
									}
								}
							} `graphql:"... on Tree"`
						}
					}
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
}

//// LIST FUNCTION

func tableGitHubRepositoryContentList1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner, repo := parseRepoFullName(d.EqualsQualString("repository_full_name"))
	var filterPath string
	if d.EqualsQualString("repository_content_path") != "" {
		filterPath = d.EqualsQualString("repository_content_path")
	}
	plugin.Logger(ctx).Debug("tableGitHubRepositoryContentList", "owner", owner, "repo", repo, "path", filterPath)

	client := connectV4(ctx, d)

	var allTrees []string

	plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "EXECUTION START", "1111111111")
	// Get root level trees
	treesPathRoot, err := getTreesPathOnly(ctx, client, repo, owner, []string{filterPath})
	if err != nil {
		plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getTreesPathOnly", "root_level", err)
		return nil, err
	}

	allTrees = append(allTrees, treesPathRoot...)
	// Remove duplicate tree paths, if there are sub directories under a directory then we should keep the dirs where the files are available. Tot the parent directory.
	allTrees = removeDuplicates(allTrees)

	plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "EXECUTION IN STEP ", "2222222222")

	// Get all the trees under the repository up to nth level
	for len(treesPathRoot) > 0 {
		allTrees = append(allTrees, treesPathRoot...)
		treesPathRoot, err = getTreesPathOnly(ctx, client, repo, owner, treesPathRoot)
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getTreesPathOnly", "recursive_api_call_for_trees", err)
		}
		allTrees = append(allTrees, treesPathRoot...)
		allTrees = removeDuplicates(allTrees)
	}

	plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "EXECUTION IN STEP ", "3333333333")

	//// Root level files contents
	trees, err := getFileContentByPath(ctx, d, client, repo, owner, []string{""})
	if err != nil {
		plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getFileContentByPath", "error_in_get_file_content_root_level", err)
	}

	allTrees = append(allTrees, trees...)
	allTrees = removeDuplicates(allTrees)
	plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "EXECUTION IN STEP ", "4444444444")

	// Fetch all the file content for the trees under a repository
	// for len(allTrees) > 0 {
	// 	innerTrees := trees
	// 	trees, err = getFileContentByPath(ctx, d, client, repo, owner, innerTrees)
	// 	if err != nil {
	// 		plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getFileContentByPath", "recursive_file_content_fetch", err)
	// 	}

	// 	allTrees = append(allTrees, trees...)
	// 	allTrees = removeDuplicates(allTrees)
	// }

	plugin.Logger(ctx).Debug("github_repository_content_1.tableGitHubRepositoryContentList1", "ALL TREES:", allTrees)

	for _, tr := range allTrees {
		trees, err = getFileContentByPath(ctx, d, client, repo, owner, []string{tr})
		if len(trees) > 0 {
			plugin.Logger(ctx).Debug("github_repository_content_1.tableGitHubRepositoryContentList1", "TREES ARE MISSING:", trees)
		}
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getFileContentByPath", "recursive_file_content_fetch", err)
		}
	}

	return nil, nil
}

func getFileContentByPath(ctx context.Context, d *plugin.QueryData, client *githubv4.Client, repo string, owner string, filePaths []string) ([]string, error) {

	// Store the results
	var contents []ContentInfo1
	var treeTypePaths []string

	for _, filePath := range filePaths {
		variables := map[string]interface{}{
			"owner":      githubv4.String(owner),
			"repo":       githubv4.String(repo),
			"expression": githubv4.String("HEAD:" + filePath),
		}
		plugin.Logger(ctx).Error("github_repository_content_1.tableGitHubRepositoryContentList1", "getFileContentByPath", "RUNNING FOR FILE PATH ", filePath)

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return []string{}, fmt.Errorf("github_repository_content1.getFileContentByPath file path ==>> %s : %s", filePath, err)
		}

		for _, data := range query.Repository.Object.Tree.Entries {
			fmt.Printf("Got file content =>> %+v\n", data.Name)
			if data.Type == "tree" {
				treeTypePaths = append(treeTypePaths, string(data.Path))
			} else {
				contents = append(contents, ContentInfo1{
					Oid:            string(data.Object.Blob.Oid),
					AbbreviatedOid: string(data.Object.Blob.AbbreviatedOid),
					Name:           string(data.Name),
					Mode:           int(data.Mode),
					PathRaw:        string(data.PathRaw),
					IsGenerated:    bool(data.IsGenerated),
					Path:           string(data.Path),
					Size:           int(data.Size),
					LineCount:      int(data.LineCount),
					Type:           string(data.Type),
					Content:        string(data.Object.Blob.Text),
					IsBinary:       bool(data.Object.Blob.IsBinary),
					CommitUrl:      string(data.Object.Blob.CommitUrl),
				})
			}
			if len(data.Object.SubTree.Entries) > 0 {
				for _, subData := range data.Object.SubTree.Entries {
					if subData.Type == "tree" {
						treeTypePaths = append(treeTypePaths, string(subData.Path))
					} else {
						fmt.Printf("Got file content =>> %+v\n", subData.Name)
						contents = append(contents, ContentInfo1{
							Oid:            string(subData.Object.Blob.Oid),
							AbbreviatedOid: string(subData.Object.Blob.AbbreviatedOid),
							Name:           string(subData.Name),
							Mode:           int(subData.Mode),
							PathRaw:        string(subData.PathRaw),
							IsGenerated:    bool(subData.IsGenerated),
							Path:           string(subData.Path),
							Size:           int(subData.Size),
							LineCount:      int(subData.LineCount),
							Type:           string(subData.Type),
							Content:        string(subData.Object.Blob.Text),
							IsBinary:       bool(subData.Object.Blob.IsBinary),
							CommitUrl:      string(subData.Object.Blob.CommitUrl),
						})
					}
				}
			}
		}

		for _, c := range contents {
			d.StreamListItem(ctx, c)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return treeTypePaths, nil
}

func getTreesPathOnly(ctx context.Context, client *githubv4.Client, repo string, owner string, filePaths []string) ([]string, error) {

	var treeTypePaths []string

	for _, filePath := range filePaths {
		// fmt.Printf("getFileContentByPath File Path ===>>>> %s\n", filePath)
		variables := map[string]interface{}{
			"owner":      githubv4.String(owner),
			"repo":       githubv4.String(repo),
			"expression": githubv4.String("HEAD:" + filePath),
		}

		// Define the query levels
		setLevelValue := "LEVEL_3"
		level3 := &treeQuery.threeLevel
		level2 := &treeQuery.twoLevel
		level1 := &treeQuery.oneLevel

		// Adapting the query depth based on repository data volume:
		// The GraphQL endpoint enforces a 10-second maximum execution time, which cannot be extended.
		// In cases of substantial data volumes, errors such as "Something went wrong while executing your query. This may be the result of a timeout" may occur.
		// The strategy involves initially querying down to 3 levels. If a timeout error occurs, the query is retried for up to 2 levels deep, and then 1 level deep as a final fallback.
		// Through testing, it's observed that this approach effectively mitigates timeout errors, even with extensive tree structures. However, surpassing 5000 subdirectories might trigger primary rate limit errors.

		plugin.Logger(ctx).Debug("github_repository_content1.getTreesPathOnly", "file path", filePath, "query_level", setLevelValue)
		err := client.Query(ctx, level3, variables)

		if err != nil {
			setLevelValue = "LEVEL_2"
			plugin.Logger(ctx).Debug("github_repository_content1.getTreesPathOnly", "file path", filePath, "query_level", setLevelValue)
			if strings.Contains(err.Error(), "Something went wrong while executing your query. This may be the result of a timeout, or it could be a GitHub bug.") || strings.Contains(err.Error(), "504 Gateway Timeout") {

				err = client.Query(ctx, level2, variables)
				if err != nil {
					plugin.Logger(ctx).Debug("github_repository_content1.getTreesPathOnly", "file path", filePath, "query_level", setLevelValue)
					if strings.Contains(err.Error(), "Something went wrong while executing your query. This may be the result of a timeout, or it could be a GitHub bug.") || strings.Contains(err.Error(), "504 Gateway Timeout") {
						setLevelValue = "LEVEL_1"
						fmt.Printf("getTreesPathOnly Level Set to ->> %s\n", setLevelValue)

						if err != nil {
							plugin.Logger(ctx).Error("github_repository_content1.getTreesPathOnly", "file path", filePath, "query_level", setLevelValue, "error", err)
						}
					}
				}
			}
		}

		// Get the results based on query level execution
		switch setLevelValue {
		case "LEVEL_3":
			for _, data := range level3.Repository.Object.Tree.Entries {
				if data.Type == "tree" && data.Path != githubv4.String(filePath) {
					treeTypePaths = append(treeTypePaths, string(data.Path))
				}
				if len(data.Object.SubTree.Entries) > 0 {
					for _, subData := range data.Object.SubTree.Entries {
						if subData.Type == "tree" && subData.Path != githubv4.String(filePath) {
							treeTypePaths = append(treeTypePaths, string(subData.Path))
						}
						if len(subData.Object.HydraTree.Entries) > 0 {
							for _, hData := range subData.Object.HydraTree.Entries {
								if hData.Type == "tree" && hData.Path != githubv4.String(filePath) {
									treeTypePaths = append(treeTypePaths, string(hData.Path))
								}
							}
						}
					}
				}
			}
		case "LEVEL_2":
			for _, data := range level2.Repository.Object.Tree.Entries {
				if data.Type == "tree" && data.Path != githubv4.String(filePath) {
					treeTypePaths = append(treeTypePaths, string(data.Path))
				}
				if len(data.Object.SubTree.Entries) > 0 {
					for _, subData := range data.Object.SubTree.Entries {
						if subData.Type == "tree" && subData.Path != githubv4.String(filePath) {
							treeTypePaths = append(treeTypePaths, string(subData.Path))
						}
					}
				}
			}
		case "LEVEL_1":
			for _, data := range level1.Repository.Object.Tree.Entries {
				if data.Type == "tree" && data.Path != githubv4.String(filePath) {
					treeTypePaths = append(treeTypePaths, string(data.Path))
				}
			}
		}

	}

	plugin.Logger(ctx).Debug("github_repository_content1.getTreesPathOnly", "numbers of trees returned:", len(treeTypePaths))

	return treeTypePaths, nil
}

// removeDuplicates function processes an input array of string paths and outputs an array with duplicates eliminated.
// Specifically, if the array includes paths such as "src/abc" and "src/abc/def", and the directory "src/abc" contains no files but has subdirectories, then "src/abc" will be removed from the output array. This is because the presence of subdirectories under "src/abc" implies it's a parent directory without direct file contents.
func removeDuplicates(inputStrings []string) []string {
	// Eliminate duplicates
	uniqueStrings := make(map[string]bool)
	for _, str := range inputStrings {
		uniqueStrings[str] = true
	}

	// Convert map keys to a slice for processing
	filteredList := make([]string, 0, len(uniqueStrings))
	for str := range uniqueStrings {
		filteredList = append(filteredList, str)
	}

	// Determines if one string is a non-exact prefix of another
	isNonExactPrefix := func(s string, list []string) bool {
		for _, item := range list {
			if s != item && strings.HasPrefix(item, s) {
				return true
			}
		}
		return false
	}

	// Remove strings that are prefixes of others
	finalList := []string{}
	for _, str := range filteredList {
		if !isNonExactPrefix(str, filteredList) {
			finalList = append(finalList, str)
		}
	}

	// Sort the final list for consistent output
	sort.Strings(finalList)

	return finalList
}
