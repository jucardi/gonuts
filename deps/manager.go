package deps

import (
	"bytes"
	"fmt"
	"github.com/jucardi/go-beans/beans"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/gonuts/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	DefaultDepsManager = "default-deps-manager"
	goListCmd          = `go list -f '{{ if .Imports }}{{ join .Imports "\n" }}{{ end }}' ./... | xargs -L1 go list -f '{{ if not .Standard }}{{ .ImportPath  }}{{ end }}'`
)

var (
	// To validate the interface implementation at compile time.
	_ IDepsManager = (*depsManager)(nil)

	depsMngrInstance IDepsManager

	packageExceptions = []string{
		"golang.org",
	}
)

type IDepsManager interface {
	Generate() (*Dependencies, error)
	Load() (*Dependencies, error)
	LoadOrGenerate() (*Dependencies, error)
}

type depsManager struct {
}

func init() {
	// Registering the bean implementation.
	beans.RegisterFunc((*IDepsManager)(nil), DefaultDepsManager, func() interface{} {
		if depsMngrInstance != nil {
			return depsMngrInstance
		}

		depsMngrInstance = &depsManager{}
		return depsMngrInstance
	})
}

func Manager() IDepsManager {
	return beans.Resolve((*IDepsManager)(nil), DefaultDepsManager).(IDepsManager)
}

func (d *depsManager) Generate() (*Dependencies, error) {
	ret := &Dependencies{}
	if err := d.generate(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *depsManager) generate(ret *Dependencies, pkg ...string) (error) {
	var pkgName string
	c := exec.Command("bash", "-c", goListCmd)

	if len(pkg) > 0 {
		c.Dir = utils.GetPathFromPkgName(pkg[0])
		pkgName = pkg[0]
	} else if dir, err := os.Getwd(); err != nil {
		return fmt.Errorf("unable to get working directory, %s", err.Error())
	} else {
		pkgName = utils.GetPkgNameFromPath(dir)
	}

	log.Infof("analyzing dependencies for %s", pkgName)
	stderr := &bytes.Buffer{}
	c.Stderr = stderr
	out, err := c.Output()
	if err != nil {
		log.Warnf("not all dependencies found for %s, this is may be normal for unused packages within that project.", pkgName)
		log.Debugf("details for %s\n%s", pkgName, string(stderr.Bytes()))
	}

	outStr := string(out)

	depList := stringx.New(outStr).Replace("\r", "", -1).Split("\n")
	depList = streams.From(depList).Filter(func(i interface{}) bool {
		x := i.(string)
		return !strings.HasPrefix(x, pkgName) && !strings.Contains(x, "golang.org")
	}).ToArray().([]string)
	for _, v := range depList {
		rootPkg, added := ret.Add(v)
		if added {
			if err := d.generate(ret, rootPkg); err != nil {
				log.Error(v, ": ", err.Error())
			}
		}
	}

	return nil
}

// Load loads the dependencies file
func (d *depsManager) Load() (*Dependencies, error) {
	if _, err := os.Stat(depsFile); os.IsNotExist(err) {
		return nil, ErrNoDepsExists
	}
	data, err := ioutil.ReadFile(depsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read deps file, %s", err.Error())
	}
	ret := &Dependencies{}
	if err := yaml.Unmarshal(data, ret); err != nil {
		return nil, fmt.Errorf("unable to deserialize deps file, %s", err.Error())
	}
	return ret, nil
}

func (d *depsManager) LoadOrGenerate() (*Dependencies, error) {
	if _, err := os.Stat(depsFile); os.IsNotExist(err) {
		return d.Generate()
	}
	return d.Load()
}
