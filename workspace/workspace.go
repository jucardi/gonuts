package workspace

import (
	"errors"
	"fmt"
	"github.com/jucardi/go-beans/beans"
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/gonuts/deps"
	"github.com/jucardi/gonuts/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultWorkspaceManager = "default-deps-manager"
	workspaceFilename       = "gonuts.yml"
)

var (
	// To validate the interface implementation at compile time.
	_ IWorkspaceManager = (*workspaceManager)(nil)

	wsMngrInstance IWorkspaceManager
)

type IWorkspaceManager interface {
	Set() error
	Reset(restoreInfo ...*RestoreInfo) error
	Load() (*RestoreInfo, error)
}

type workspaceManager struct {
}

type RestoreInfo struct {
	Name        string
	RestoreInfo []*PkgInfo
}

type PkgInfo struct {
	Name     string
	Revision string
}

func Manager() IWorkspaceManager {
	return beans.Resolve((*IWorkspaceManager)(nil), DefaultWorkspaceManager).(IWorkspaceManager)
}

func init() {
	// Registering the bean implementation.
	beans.RegisterFunc((*IWorkspaceManager)(nil), DefaultWorkspaceManager, func() interface{} {
		if wsMngrInstance != nil {
			return wsMngrInstance
		}

		wsMngrInstance = &workspaceManager{}
		return wsMngrInstance
	})
}

func (d *RestoreInfo) Save() error {
	data, err := yaml.Marshal(d)

	if err != nil {
		return fmt.Errorf("unable to save workspace restore info, %s", err.Error())
	}

	return ioutil.WriteFile(getWsFilePath(), data, 0644)
}

func (w *workspaceManager) Set() error {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("unable to set workspace, %s", err.Error())
	}
	if rootLvl, err := utils.GetRootDir(""); err == nil {
		workDir = rootLvl
	}

	if _, err := os.Stat(paths.Combine(os.Getenv("GOPATH"), "src", workDir)); os.IsNotExist(err) {
		return errors.New("must be in a project inside GOPATH")
	}

	d, err := deps.Manager().LoadOrGenerate()
	if err != nil {
		return err
	}

	split := strings.Split(workDir, "/")
	pkgName := split[len(split)-1]
	info := &RestoreInfo{
		Name: pkgName,
	}
	for k := range d.Dependencies {
		branch, err := utils.GetBranch(k)
		if err != nil {
			return fmt.Errorf("unable to read current HEAD in %s, %s\naborting.", k, err.Error())
		}
		info.RestoreInfo = append(info.RestoreInfo, &PkgInfo{
			Name:     k,
			Revision: branch,
		})
	}
	if err := info.Save(); err != nil {
		return err
	}
	for k, v := range d.Dependencies {
		if err := utils.Checkout(k, v.Revision); err != nil {
			w.Reset(info)
			return fmt.Errorf("setting workspace failed, unable to checkout revision %s for package %s", k, v.Revision)
		}
	}
	return nil
}

func (w *workspaceManager) Reset(restoreInfo ...*RestoreInfo) error {
	var info *RestoreInfo

	if len(restoreInfo) > 0 && restoreInfo[0] != nil {
		info = restoreInfo[0]
	} else {
		if _, err := os.Stat(getWsFilePath()); os.IsNotExist(err) {
			return nil
		}

		inf, err := w.Load()
		if err != nil {
			return err
		}
		info = inf
	}

	errBuilder := stringx.Builder()
	for _, v := range info.RestoreInfo {
		if err := utils.Checkout(v.Name, v.Revision); err != nil {
			errBuilder.AppendLine(err.Error())
		}
	}
	if errBuilder.IsEmpty() {
		if err := os.Remove(getWsFilePath()); err != nil {
			return fmt.Errorf("workspace restored, unable to delete restore info file, %s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("errors occurred while restoring workspace:\n\n%s", errBuilder.Build())
}

// Load loads the dependencies file
func (w *workspaceManager) Load() (*RestoreInfo, error) {
	if _, err := os.Stat(getWsFilePath()); os.IsNotExist(err) {
		return nil, nil
	}
	data, err := ioutil.ReadFile(getWsFilePath())
	if err != nil {
		return nil, fmt.Errorf("unable to read restore workspace file, %s", err.Error())
	}
	ret := &RestoreInfo{}
	if err := yaml.Unmarshal(data, ret); err != nil {
		return nil, fmt.Errorf("unable to deserialize deps file, %s", err.Error())
	}
	return ret, nil
}

func getWsFilePath() string {
	return paths.Combine(os.Getenv("GOPATH"), "src", workspaceFilename)
}
