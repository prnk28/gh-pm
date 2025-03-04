package exc

import (
	"github.com/cli/go-gh"
)

func Gh(args ...string) (string, error) {
	out, _, err := gh.Exec(args...)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
