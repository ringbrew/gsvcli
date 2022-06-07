package subcmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type GenDomain struct {
	module string
	domain string
	tmpl   map[string]string
}

func NewGenDomain(domain string) *GenDomain {
	module := ""
	if file, err := os.Open("go.mod"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, "module ") {
				module = strings.TrimPrefix(text, "module ")
				break
			}
		}
	} else {
		log.Fatal("[FATAL] " + err.Error())
	}

	if module == "" {
		log.Fatal("[FATAL] " + "please init your go module with go mod init")
	}

	return &GenDomain{
		module: module,
		domain: domain,
		tmpl: map[string]string{
			"init.go":    initTmpl,
			"entity.go":  entityTmpl,
			"usecase.go": useCaseTmpl,
			"repo.go":    repoTmpl,
		},
	}
}

func (gd *GenDomain) Process() error {
	domainPath := filepath.Join("internal", "domain", gd.domain)

	entity := strings.Title(strings.ReplaceAll(gd.domain, " ", ""))

	if err := os.MkdirAll(domainPath, os.ModePerm); err != nil {
		return err
	}

	fileList := make([]*os.File, 0, len(gd.tmpl))
	defer func() {
		for i := range fileList {
			fileList[i].Close()
		}
	}()

	for fileName, tmpl := range gd.tmpl {
		fileFullPath := filepath.Join(domainPath, fileName)

		if _, err := os.Stat(fileFullPath); err != nil && os.IsNotExist(err) {
			file, err := os.Create(fileFullPath)
			if err != nil {
				return err
			}
			fileList = append(fileList, file)

			t, err := template.New(fileName).Parse(tmpl)
			if err != nil {
				return err
			}

			if err := t.Execute(file, map[string]interface{}{
				"projectName": gd.module,
				"domain":      gd.domain,
				"entity":      entity,
			}); err != nil {
				return err
			}
		}
	}

	c := exec.Command("go", "fmt")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = domainPath
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
