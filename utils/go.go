package utils

import (
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/go-strings/stringx"
	"os"
)

func GetPkgNameFromPath(dir string) string {
	return stringx.New(dir).
		Replace(paths.Combine(os.Getenv("GOPATH"), "src")+"/", "", -1).
		Trim("\n").
		Trim("\r").
		TrimSpace().
		S()
}

func GetPathFromPkgName(pkg string) string {
	return paths.Combine(os.Getenv("GOPATH"), "src", pkg)
}

func GetRootFromPkg(pkg string) (string, error) {
	if dir, err := GetRootDir(pkg); err != nil {
		return "", err
	} else {
		return GetPkgNameFromPath(dir), nil
	}
}
