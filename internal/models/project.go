package models

// Organization represents a GitHub organization
type Organization struct {
	Login      string
	ProjectsV2 struct {
		Nodes    []Project
		PageInfo struct {
			HasNextPage bool
			EndCursor   string
		}
	}
}

// Project represents a GitHub ProjectV2
type Project struct {
	ID        string
	Number    int
	Title     string
	URL       string
	Closed    bool
	CreatedAt string
}

// OrganizationProjectsResponse represents the GraphQL response for organizations and their projects
type OrganizationProjectsResponse struct {
	Viewer struct {
		Organizations struct {
			Nodes    []Organization
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		}
	}
}

// SingleOrgProjectsResponse represents the GraphQL response for a single organization's projects
type SingleOrgProjectsResponse struct {
	Organization struct {
		ProjectsV2 struct {
			Nodes    []Project
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		}
	}
}
