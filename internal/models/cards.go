package models

// ProjectCard represents a card in a project
type ProjectCard struct {
	ID          string
	Title       string
	Status      string
	URL         string
	ContentType string // "Issue", "PullRequest", or "DraftIssue"
	Content     interface{}
}

// ProjectCardIssue represents an issue linked to a project card
type ProjectCardIssue struct {
	ID         string
	Title      string
	Number     int
	URL        string
	State      string
	Repository struct {
		Name  string
		Owner string
	}
	Labels      []Label
	Assignees   []User
	CreatedAt   string
	UpdatedAt   string
	ClosedAt    string
	Author      User
	Body        string
	CommentCount int
}

// ProjectCardPR represents a pull request linked to a project card
type ProjectCardPR struct {
	ID         string
	Title      string
	Number     int
	URL        string
	State      string
	Repository struct {
		Name  string
		Owner string
	}
	Labels       []Label
	Assignees    []User
	CreatedAt    string
	UpdatedAt    string
	ClosedAt     string
	MergedAt     string
	Author       User
	Body         string
	CommentCount int
	Merged       bool
	Draft        bool
	BaseRef      string
	HeadRef      string
}

// ProjectColumn represents a column in a project board
type ProjectColumn struct {
	ID    string
	Name  string
	Cards []ProjectCard
}

// Label represents a GitHub issue or PR label
type Label struct {
	ID          string
	Name        string
	Description string
	Color       string
}
