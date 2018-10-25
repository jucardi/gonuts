package utils

import (
	"github.com/jucardi/go-logger-lib/log"
	"os/exec"
	"strings"
)

func GetRootDir(pkg string) (string, error) {
	result, err := gitExec("git rev-parse --show-toplevel", pkg)
	if err != nil {
		return "", err
	}
	return GetPkgNameFromPath(result), nil
}

func GetRevisionHash(pkg string) (string, error) {
	return gitExec("git rev-parse HEAD", pkg)
}

func GetBranch(pkg string) (string, error) {
	result, err := gitExec(`git branch 2>/dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/'`, pkg)
	if err != nil {
		return "", err
	}
	if strings.Contains(result, "(HEAD detached at") {
		return GetRevisionHash(pkg)
	}
	return result, nil
}

func Checkout(pkg, revision string) error {
	msg, err := gitExec("git checkout "+revision, pkg)
	log.Debug(msg)
	return err
}

func gitExec(gitCmd, pkg string) (string, error) {
	cmd := exec.Command("bash", "-c", gitCmd)
	if pkg != "" {
		cmd.Dir = GetPathFromPkgName(pkg)
	}

	if result, err := cmd.Output(); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}
