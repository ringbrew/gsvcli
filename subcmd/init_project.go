package subcmd

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type InitProject struct {
	projectName  string
	forkEndpoint string
	tempPath     string
}

func NewInitProject(projectName string) *InitProject {
	return &InitProject{
		projectName:  projectName,
		forkEndpoint: "https://github.com/ringbrew/gsvtmpl.git",
		tempPath:     ".tmp",
	}
}

func (p InitProject) Check() error {
	if Exists(filepath.Base(p.projectName)) {
		return errors.New("project path already exist")
	}

	return nil
}

func (p InitProject) SetGoEnv() error {
	c := exec.Command("go", "env", `GOPROXY="https://goproxy.cn"`)
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("go", "env", `GO111MODULE="on"`)
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

func (p InitProject) GetTemplate() error {
	if !Exists(p.tempPath) {
		if err := os.Mkdir(p.tempPath, os.ModePerm); err != nil {
			return err
		}
	} else {
		if err := os.RemoveAll(p.tempPath); err != nil {
			return err
		}

		if err := os.Mkdir(p.tempPath, os.ModePerm); err != nil {
			return err
		}
	}

	c := exec.Command("git", "clone", p.forkEndpoint)
	c.Dir = p.tempPath
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func (p InitProject) Render() error {
	projectName := p.projectName
	baseName := filepath.Base(projectName)

	err := filepath.Walk(p.tempPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ft, err := template.ParseFiles(path); err == nil {
				if file, err := os.Create(path); err == nil {
					if err := ft.Execute(file, map[string]interface{}{
						"projectName": projectName,
						"baseName":    baseName,
					}); err != nil {
						return err
					}
					if err := file.Close(); err != nil {
						log.Println(err.Error())
					}
				}
			}
			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func (p InitProject) GetTmpProjectPath() string {
	b := filepath.Base(p.forkEndpoint)

	return filepath.Join(p.tempPath, strings.TrimSuffix(b, filepath.Ext(b)))
}

func (p InitProject) WriteDefaultConfig() error {
	file, err := os.Create(filepath.Join(p.GetTmpProjectPath(), "config.yaml"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(ConfigTmpl)
	if err != nil {
		return err
	}

	return nil
}

func (p InitProject) Complete() error {
	trn := p.GetTmpProjectPath()

	if err := os.RemoveAll(filepath.Join(trn, ".git")); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(trn, "cmd.go")); err == nil {
		if err := os.Rename(filepath.Join(trn, "cmd.go"), filepath.Join(trn, "cmd", strings.ReplaceAll(filepath.Base(p.projectName), "-", "_")+".go")); err != nil {
			return err
		}
	}

	c := exec.Command("go", "mod", "init", p.projectName)
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("go", "mod", "tidy")
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("go", "fmt", "./...")
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "init")
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "add", ".")
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("git", "commit", "-m", "initial project")
	c.Dir = trn
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	if err := p.WriteDefaultConfig(); err != nil {
		return err
	}

	if err := os.Rename(trn, filepath.Base(p.projectName)); err != nil {
		return err
	}

	if err := os.RemoveAll(p.tempPath); err != nil {
		return err
	}

	return nil
}
