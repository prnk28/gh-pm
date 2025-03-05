package ghcli

import (
	"fmt"
	"time"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/prnk28/gh-pm/internal/models"
)

// UserInfo contains comprehensive information about a GitHub user
type UserInfo struct {
	Profile        models.User
	Organizations  []models.OrganizationBasic
	StarredRepos   []models.Repository
	RecentActivity ActivitySummary
}

// ActivitySummary contains summarized activity information
type ActivitySummary struct {
	TotalContributions int
	CommitsLastMonth   int
	PRsLastMonth       int
	IssuesLastMonth    int
}

// GetUserInfo retrieves comprehensive information about the authenticated user
func GetUserInfo() (*UserInfo, error) {
	// Create a GraphQL client
	client, err := gh.GQLClient(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GraphQL client: %w", err)
	}

	// Get current date for time-based queries
	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0).Format("2006-01-02T15:04:05Z")

	// Variables for the query
	variables := map[string]any{
		"orgCount":     graphql.Int(100),
		"repoCount":    graphql.Int(100),
		"starredCount": graphql.Int(100),
		"since":        graphql.String(oneMonthAgo),
	}

	// Define the query
	var query struct {
		Viewer struct {
			Login           graphql.String
			Name            graphql.String
			Bio             graphql.String
			Company         graphql.String
			Location        graphql.String
			Email           graphql.String
			WebsiteURL      graphql.String
			TwitterUsername graphql.String
			AvatarURL       graphql.String
			FollowerCount   graphql.Int
			FollowingCount  graphql.Int

			Repositories struct {
				TotalCount graphql.Int
				Nodes      []struct {
					ID              graphql.String
					Name            graphql.String
					NameWithOwner   graphql.String
					Description     graphql.String
					URL             graphql.String
					StargazerCount  graphql.Int
					ForkCount       graphql.Int
					IsPrivate       graphql.Boolean
					IsArchived      graphql.Boolean
					PrimaryLanguage *struct {
						Name  graphql.String
						Color graphql.String
					}
					UpdatedAt graphql.String
				}
				PageInfo struct {
					HasNextPage graphql.Boolean
					EndCursor   graphql.String
				}
			} `graphql:"repositories(first: $repoCount, orderBy: {field: UPDATED_AT, direction: DESC})"`

			StarredRepositories struct {
				TotalCount graphql.Int
				Nodes      []struct {
					ID              graphql.String
					Name            graphql.String
					NameWithOwner   graphql.String
					Description     graphql.String
					URL             graphql.String
					StargazerCount  graphql.Int
					ForkCount       graphql.Int
					IsPrivate       graphql.Boolean
					IsArchived      graphql.Boolean
					PrimaryLanguage *struct {
						Name  graphql.String
						Color graphql.String
					}
					UpdatedAt graphql.String
				}
				PageInfo struct {
					HasNextPage graphql.Boolean
					EndCursor   graphql.String
				}
			} `graphql:"starredRepositories(first: $starredCount, orderBy: {field: STARRED_AT, direction: DESC})"`

			Organizations struct {
				TotalCount graphql.Int
				Nodes      []struct {
					Login     graphql.String
					Name      graphql.String
					AvatarURL graphql.String
					URL       graphql.String
				}
				PageInfo struct {
					HasNextPage graphql.Boolean
					EndCursor   graphql.String
				}
			} `graphql:"organizations(first: $orgCount)"`

			ContributionsCollection struct {
				TotalCommitContributions            graphql.Int
				TotalIssueContributions             graphql.Int
				TotalPullRequestContributions       graphql.Int
				TotalPullRequestReviewContributions graphql.Int
				ContributionCalendar                struct {
					TotalContributions graphql.Int
				}
				CommitContributionsByRepository []struct {
					Contributions struct {
						TotalCount graphql.Int
					} `graphql:"contributions(first: 1)"`
				} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
			} `graphql:"contributionsCollection(from: $since)"`

			CreatedAt graphql.String
		}
	}

	// Execute the query
	err = client.Query("UserInfo", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("error querying GraphQL API: %w", err)
	}

	// Convert query results to our models
	viewer := query.Viewer

	// Convert organizations
	orgs := make([]models.OrganizationBasic, 0, len(viewer.Organizations.Nodes))
	for _, org := range viewer.Organizations.Nodes {
		orgs = append(orgs, models.OrganizationBasic{
			Login:     string(org.Login),
			Name:      string(org.Name),
			AvatarURL: string(org.AvatarURL),
			URL:       string(org.URL),
		})
	}

	// Convert starred repositories
	starredRepos := make([]models.Repository, 0, len(viewer.StarredRepositories.Nodes))
	for _, repo := range viewer.StarredRepositories.Nodes {
		var primaryLanguage *models.Language
		if repo.PrimaryLanguage != nil {
			primaryLanguage = &models.Language{
				Name:  string(repo.PrimaryLanguage.Name),
				Color: string(repo.PrimaryLanguage.Color),
			}
		}

		starredRepos = append(starredRepos, models.Repository{
			ID:              string(repo.ID),
			Name:            string(repo.Name),
			NameWithOwner:   string(repo.NameWithOwner),
			Description:     string(repo.Description),
			URL:             string(repo.URL),
			StargazerCount:  int(repo.StargazerCount),
			ForkCount:       int(repo.ForkCount),
			IsPrivate:       bool(repo.IsPrivate),
			IsArchived:      bool(repo.IsArchived),
			PrimaryLanguage: primaryLanguage,
			UpdatedAt:       string(repo.UpdatedAt),
		})
	}

	// Create activity summary
	activitySummary := ActivitySummary{
		TotalContributions: int(viewer.ContributionsCollection.ContributionCalendar.TotalContributions),
		CommitsLastMonth:   int(viewer.ContributionsCollection.TotalCommitContributions),
		PRsLastMonth:       int(viewer.ContributionsCollection.TotalPullRequestContributions),
		IssuesLastMonth:    int(viewer.ContributionsCollection.TotalIssueContributions),
	}

	// Create and return UserInfo
	userInfo := &UserInfo{
		Profile: models.User{
			Login:           string(viewer.Login),
			Name:            string(viewer.Name),
			Bio:             string(viewer.Bio),
			Company:         string(viewer.Company),
			Location:        string(viewer.Location),
			Email:           string(viewer.Email),
			WebsiteURL:      string(viewer.WebsiteURL),
			TwitterUsername: string(viewer.TwitterUsername),
			AvatarURL:       string(viewer.AvatarURL),
			FollowerCount:   int(viewer.FollowerCount),
			FollowingCount:  int(viewer.FollowingCount),
			CreatedAt:       string(viewer.CreatedAt),
		},
		Organizations:  orgs,
		StarredRepos:   starredRepos,
		RecentActivity: activitySummary,
	}

	return userInfo, nil
}
