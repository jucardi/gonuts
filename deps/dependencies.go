package deps

import (
	"errors"
	"fmt"
	"github.com/jucardi/gonuts/utils"
	"github.com/jucardi/infuse/util/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const depsFile = "gonuts.yml"

var (
	ErrNoDepsExists = errors.New("deps file not found")
)

type Dependencies struct {
	Dependencies DepMap
}

type DepMap map[string]*DependencyInfo

// DependencyInfo
type DependencyInfo struct {
	Revision string
}

// Add adds a dependency to the dependency list, returns false if the dependency already existed or if it is unable to fetch the git HEAD id, otherwise true
func (d *Dependencies) Add(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	rootPkg, err := utils.GetRootFromPkg(name)
	if err != nil {
		log.Errorf("unable to get git root from package %s, %s", name, err.Error())
		return "", false
	}
	if d.Dependencies == nil {
		d.Dependencies = DepMap{}
	}
	if _, ok := d.Dependencies[rootPkg]; ok {
		return "", false
	}

	hash, err := utils.GetRevisionHash(rootPkg)
	if err != nil {
		log.Errorf("unable to resolve git HEAD id for package %s, %s", rootPkg, err.Error())
		return "", false
	}

	d.Dependencies[rootPkg] = &DependencyInfo{
		Revision: hash,
	}
	return rootPkg, true
}

// List returns a list of the registered dependencies
func (d *Dependencies) List() []string {
	var ret []string

	for k := range d.Dependencies {
		ret = append(ret, k)
	}

	return ret
}

// Save saves the dependencies info
func (d *Dependencies) Save() error {
	data, err := yaml.Marshal(d)

	if err != nil {
		return fmt.Errorf("unable to serialize dependencies, %s", err.Error())
	}

	return ioutil.WriteFile(depsFile, data, 0644)
}
