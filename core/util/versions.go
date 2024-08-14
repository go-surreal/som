package util

import (
	"fmt"
	"os/exec"
	"strings"
)

func SOMVersion() (string, error) {
	return checkVersion(pkgSOM)
}

func checkVersion(pkg string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-u", pkg+"@latest")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to check version of %s: %w", pkg, err)
	}

	parts := strings.Split(string(out), " ")

	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected output: %s", out)
	}

	if parts[0] != pkg {
		return "", fmt.Errorf("unexpected package in output: %s", parts[0])
	}

	version, err := versionOrdinal(parts[1])
	if err != nil {
		return "", fmt.Errorf("could not parse version: %w", err)
	}

	return version, nil
}
