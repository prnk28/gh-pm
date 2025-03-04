package app

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	contextKeyUser   contextKey = "gh-pm.user"
	contextKeySystem contextKey = "gh-pm.system"
)

func ctxKeyUser(name string) string {
	return fmt.Sprintf("%s.%s", contextKeyUser.String(), name)
}

func ctxKeySystem(name string) string {
	return fmt.Sprintf("%s.%s", contextKeySystem.String(), name)
}

func hasCmdInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func jsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func DepStatus() map[string]bool {
	installStatus := map[string]bool{
		"gh":  hasCmdInstalled("gh"),
		"fzf": hasCmdInstalled("fzf"),
		"jq":  hasCmdInstalled("jq"),
		"gum": hasCmdInstalled("gum"),
	}
	return installStatus
}
