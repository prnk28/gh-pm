package ghcli

import (
	"fmt"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

func Exec(args ...string) (string, error) {
	gh.RESTClient(nil)
	out, _, err := gh.Exec(args...)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func CurrentRepo() (string, error) {
	r, err := gh.CurrentRepository()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", r.Owner(), r.Name()), nil
}

func Whoami() (string, error) {
	ghclicmd := []string{"api", "user", "--jq", `"You are @\(.login) (\(.name))"`}
	out, err := Exec(ghclicmd...)
	if err != nil {
		return "", err
	}
	return out, nil
}

func RESTClient() (api.RESTClient, error) {
	return gh.RESTClient(nil)
}

func GQLClient() (api.GQLClient, error) {
	return gh.GQLClient(nil)
}
