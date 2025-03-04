package main

import (
	"fmt"

	"github.com/cli/go-gh"
)

func main() {
	args := []string{"api", "user", "--jq", `"You are @\(.login) (\(.name))"`}
	stdout, _, err := gh.Exec(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stdout.String())
}
