package models

// User represents a GitHub user profile
type User struct {
	Login                   string
	Name                    string
	Bio                     string
	Company                 string
	Location                string
	Email                   string
	WebsiteURL              string
	TwitterUsername         string
	AvatarURL               string
	FollowerCount           int
	FollowingCount          int
	Repositories            RepositoryConnection
	StarredRepositories     RepositoryConnection
	Organizations           OrganizationConnection
	ContributionsCollection ContributionsCollection
	CreatedAt               string
}

// RepositoryConnection represents a paginated list of repositories
type RepositoryConnection struct {
	TotalCount int
	Nodes      []Repository
	PageInfo   PageInfo
}

// Repository represents a GitHub repository
type Repository struct {
	ID              string
	Name            string
	NameWithOwner   string
	Description     string
	URL             string
	StargazerCount  int
	ForkCount       int
	IsPrivate       bool
	IsArchived      bool
	PrimaryLanguage *Language
	UpdatedAt       string
}

// Language represents a programming language
type Language struct {
	Name  string
	Color string
}

// OrganizationConnection represents a paginated list of organizations
type OrganizationConnection struct {
	TotalCount int
	Nodes      []OrganizationBasic
	PageInfo   PageInfo
}

// OrganizationBasic represents basic information about a GitHub organization
type OrganizationBasic struct {
	Login     string
	Name      string
	AvatarURL string
	URL       string
}

// ContributionsCollection represents a user's contributions
type ContributionsCollection struct {
	TotalCommitContributions            int
	TotalIssueContributions             int
	TotalPullRequestContributions       int
	TotalPullRequestReviewContributions int
	ContributionCalendar                ContributionCalendar
}

// ContributionCalendar represents a user's contribution calendar
type ContributionCalendar struct {
	TotalContributions int
	Weeks              []ContributionWeek
}

// ContributionWeek represents a week of contributions
type ContributionWeek struct {
	ContributionDays []ContributionDay
}

// ContributionDay represents a day of contributions
type ContributionDay struct {
	Date              string
	ContributionCount int
}

// PageInfo represents pagination information
type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

// UserInfoResponse represents the GraphQL response for user information
type UserInfoResponse struct {
	Viewer User
}
