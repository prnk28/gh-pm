package ghc

var (
	// QueryProjectList is a command to query the GitHub API for a list of projects
	QueryProjectList = newCommand("project list --limit 100 --format json -L 100 --jq .items")

	// QueryProjectItemList is a command to query the GitHub API for a list of project items
	QueryProjectItemList = newCommand("project item-list 4 --owner coindotfi --format json -L 100 --jq .items")

	// QueryUserWhoami is a command to query the GitHub API for the current user
	QueryUserWhoami = newCommand("api user")
)
