package ghcli

import (
	"fmt"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/prnk28/gh-pm/internal/models"
)

// GetProjectCards retrieves all cards for a given project ID
func GetProjectCards(projectID string) ([]models.ProjectColumn, error) {
	// Create a GraphQL client
	client, err := gh.GQLClient(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GraphQL client: %w", err)
	}

	// Define the GraphQL query
	var query struct {
		Node struct {
			ProjectV2 struct {
				ID     graphql.String
				Title  graphql.String
				Closed graphql.Boolean
				Fields struct {
					Nodes []struct {
						// For status fields
						SingleSelectField struct {
							ID      graphql.String
							Name    graphql.String
							Options []struct {
								ID   graphql.String
								Name graphql.String
							}
						} `graphql:"... on ProjectV2SingleSelectField"`
					}
				} `graphql:"fields(first: 20)"`
				Items struct {
					Nodes []struct {
						ID          graphql.String
						Type        graphql.String
						FieldValues struct {
							Nodes []struct {
								// For titles
								Title struct {
									Title graphql.String
								} `graphql:"... on ProjectV2ItemFieldTextValue"`
								// For status
								SingleSelectOption struct {
									Field struct {
										ID   graphql.String
										Name graphql.String
									}
									Name graphql.String
								} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
							}
						} `graphql:"fieldValues(first: 20)"`
						Content struct {
							TypeName graphql.String `graphql:"__typename"`
							Issue    struct {
								Title graphql.String
								URL   graphql.String
							} `graphql:"... on Issue"`
							PullRequest struct {
								Title graphql.String
								URL   graphql.String
							} `graphql:"... on PullRequest"`
						} `graphql:"content"`
					}
					PageInfo struct {
						HasNextPage graphql.Boolean
						EndCursor   graphql.String
					}
				} `graphql:"items(first: 100)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $id)"`
	}

	// Set the variables for the query
	variables := map[string]any{
		"id": graphql.ID(projectID),
	}

	// Execute the query
	err = client.Query("ProjectCards", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("error querying GraphQL API: %w", err)
	}

	// Process the response
	project := query.Node.ProjectV2

	// Map status field IDs to their names
	statusFields := make(map[string]bool)
	for _, field := range project.Fields.Nodes {
		if string(field.SingleSelectField.Name) == "Status" {
			statusFields[string(field.SingleSelectField.ID)] = true
		}
	}

	// Create a map to organize cards by status
	cardsByStatus := make(map[string][]models.ProjectCard)

	// Process each item in the project
	for _, item := range project.Items.Nodes {
		var card models.ProjectCard
		card.ID = string(item.ID)

		// Set default status
		card.Status = "To Do" // Default status

		// Determine the title and URL from the content
		if string(item.Content.TypeName) == "Issue" {
			card.Title = string(item.Content.Issue.Title)
			card.URL = string(item.Content.Issue.URL)
			card.ContentType = "Issue"
			
			// We could populate a full ProjectCardIssue here if needed
			// This would require extending the GraphQL query
		} else if string(item.Content.TypeName) == "PullRequest" {
			card.Title = string(item.Content.PullRequest.Title)
			card.URL = string(item.Content.PullRequest.URL)
			card.ContentType = "PullRequest"
			
			// We could populate a full ProjectCardPR here if needed
			// This would require extending the GraphQL query
		} else {
			// For draft issues or other content types, get title from field values
			for _, fieldValue := range item.FieldValues.Nodes {
				if string(fieldValue.Title.Title) != "" {
					card.Title = string(fieldValue.Title.Title)
					card.ContentType = "DraftIssue"
				}
			}
		}

		// Get the status from field values
		for _, fieldValue := range item.FieldValues.Nodes {
			// Check for status field values
			if string(fieldValue.SingleSelectOption.Name) != "" &&
				string(fieldValue.SingleSelectOption.Field.Name) == "Status" {
				card.Status = string(fieldValue.SingleSelectOption.Name)
				break
			}
		}

		// Skip cards with empty titles (likely invalid items)
		if card.Title == "" {
			continue
		}

		// Add card to the appropriate status group
		cardsByStatus[card.Status] = append(cardsByStatus[card.Status], card)
	}

	// Convert the map to a slice of columns
	columns := make([]models.ProjectColumn, 0, len(cardsByStatus))

	// Define the standard column order
	standardColumns := []string{"To Do", "In Progress", "Done"}

	// First add standard columns in the preferred order
	for _, colName := range standardColumns {
		if cards, ok := cardsByStatus[colName]; ok {
			columns = append(columns, models.ProjectColumn{
				Name:  colName,
				Cards: cards,
			})
			delete(cardsByStatus, colName)
		}
	}

	// Then add any remaining custom columns
	for status, cards := range cardsByStatus {
		columns = append(columns, models.ProjectColumn{
			Name:  status,
			Cards: cards,
		})
	}

	return columns, nil
}
