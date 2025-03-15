package mod

import (
	"os"
	"os/exec"
	"strings"
)

// FindModuleRoot finds the module root directory (where the go.mod file exists).
func FindModuleRoot() (string, error) {
	defaultRoot, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return defaultRoot, nil
	}
	goMod := string(out)
	goModDir := strings.TrimSuffix(goMod, "/go.mod\n")
	return goModDir, nil
}
