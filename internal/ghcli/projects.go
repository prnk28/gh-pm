package ghcli

import (
	"fmt"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/prnk28/gh-pm/internal/models"
)

// ProjectResult contains a project and its organization
type ProjectResult struct {
	OrgLogin string
	Project  models.Project
}

// GetAllUserProjects retrieves all projects across all organizations the authenticated user has access to
func GetAllUserProjects() ([]ProjectResult, error) {
	// Create a GraphQL client
	client, err := gh.GQLClient(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GraphQL client: %w", err)
	}

	var results []ProjectResult

	// Variables for pagination
	var query models.OrganizationProjectsResponse
	orgCursor := (*graphql.String)(nil)

	// Loop through all organizations
	for {
		// Query organizations
		variables := map[string]any{
			"orgCursor": orgCursor,
		}

		err := client.Query("OrgProjects", &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error querying GraphQL API: %w", err)
		}

		// Process each organization
		for _, org := range query.Viewer.Organizations.Nodes {
			// Paginate through all projects in this organization
			projectCursor := (*graphql.String)(nil)
			hasMoreProjects := true

			for hasMoreProjects {
				variables := map[string]any{
					"orgLogin":      graphql.String(org.Login),
					"projectCursor": projectCursor,
				}

				var orgQuery models.SingleOrgProjectsResponse

				err := client.Query("OrgProjects", &orgQuery, variables)
				if err != nil {
					// Log error but continue with other orgs
					fmt.Printf("Error querying projects for %s: %v\n", org.Login, err)
					break
				}

				// Add projects to results
				for _, project := range orgQuery.Organization.ProjectsV2.Nodes {
					results = append(results, ProjectResult{
						OrgLogin: org.Login,
						Project:  project,
					})
				}

				// Check if we need to paginate
				hasMoreProjects = orgQuery.Organization.ProjectsV2.PageInfo.HasNextPage
				if hasMoreProjects {
					projectCursor = (*graphql.String)(&orgQuery.Organization.ProjectsV2.PageInfo.EndCursor)
				}
			}
		}

		// Check if we need to paginate organizations
		if !query.Viewer.Organizations.PageInfo.HasNextPage {
			break
		}
		orgCursor = (*graphql.String)(&query.Viewer.Organizations.PageInfo.EndCursor)
	}

	return results, nil
}
