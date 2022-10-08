package subcmd

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type GenHandler struct {
	module    string
	domain    string
	subdomain string
}

const httpHandlerGenTmpl = `package [[.domain]]

import (
	"[[.module]]/internal/domain"
	"github.com/ringbrew/gsv/service"
	"net/http"
)

type [[.prefix]]Handler struct {
	ctx           *domain.UseCaseContext
}

func New[[.prefix]]Handler(ctx *domain.UseCaseContext) *[[.prefix]]Handler {
	return &[[.prefix]]Handler{
		ctx:           ctx,
	}
}

func (h *[[.prefix]]Handler) Query(w http.ResponseWriter, r *http.Request) {
	panic("not implement")
}
func (h *[[.prefix]]Handler) Get(w http.ResponseWriter, r *http.Request) {
	panic("not implement")
}
func (h *[[.prefix]]Handler) Post(w http.ResponseWriter, r *http.Request) {
	panic("not implement")
}
func (h *[[.prefix]]Handler) Put(w http.ResponseWriter, r *http.Request) {
	panic("not implement")
}
func (h *[[.prefix]]Handler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("not implement")
}

func (h *[[.prefix]]Handler) HttpRoute() []service.HttpRoute {
	result := []service.HttpRoute{
		service.NewHttpRoute(http.MethodGet, "/[[.domain]]", h.Query, service.HttpMeta{}),
		service.NewHttpRoute(http.MethodGet, "/[[.domain]]/{id}", h.Get, service.HttpMeta{}),
		service.NewHttpRoute(http.MethodPost, "/[[.domain]]", h.Post, service.HttpMeta{}),
		service.NewHttpRoute(http.MethodPut, "/[[.domain]]/{id}", h.Put, service.HttpMeta{}),
		service.NewHttpRoute(http.MethodDelete, "/[[.domain]]/{id}", h.Delete, service.HttpMeta{}),
	}
	return result
}
`

func NewGenHandler(domain string, subdomain string) *GenHandler {
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

	return &GenHandler{
		module:    module,
		domain:    domain,
		subdomain: subdomain,
	}
}

func (gh *GenHandler) Process() error {
	deliveryPath := filepath.Join("internal", "delivery", gh.domain)

	fn := fmt.Sprintf("%s_handler.go", gh.domain)
	fp := filepath.Join(deliveryPath, fn)

	sfn := "service.http.impl.go"
	sfp := filepath.Join(deliveryPath, sfn)

	prefix := cases.Title(language.Und).String(gh.subdomain)

	if _, err := os.Stat(sfp); err != nil {
		return fmt.Errorf("error[%s] please init service.http.impl.go first", err.Error())
	}

	if _, err := os.Stat(fp); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		return nil
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New(fn).Delims("[[", "]]").Parse(httpHandlerGenTmpl)
	if err != nil {
		return err
	}

	if err := t.Execute(file, map[string]interface{}{
		"module": gh.module,
		"domain": gh.domain,
		"prefix": prefix,
	}); err != nil {
		return err
	}

	serviceFileData, err := os.ReadFile(sfp)
	if err != nil {
		return err
	}

	lines := strings.Split(string(serviceFileData), "\n")

	result := make([]string, 0, len(lines)+2)

	found := false

	/*
		handler := NewChildcareHandler(ctx)
		s.desc.HttpRoute = append(s.desc.HttpRoute, handler.HttpRoute()...)
	*/
	for _, v := range lines {
		if strings.Contains(v, "return s") {
			found = true
			result = append(result, fmt.Sprintf("%sHandler := New%sHandler(ctx)", gh.subdomain, prefix))
			result = append(result, fmt.Sprintf("s.desc.HttpRoute = append(s.desc.HttpRoute, %sHandler.HttpRoute()...)", gh.subdomain))
		}
		result = append(result, v)
	}

	if !found {
		return errors.New("error sentry not found, use \"return s\" as sentry now")
	}

	serviceFile, err := os.OpenFile(fp, os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	for _, v := range result {
		if _, err := serviceFile.WriteString(v + "\n"); err != nil {
			return err
		}
	}

	c := exec.Command("go", "fmt")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = deliveryPath
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
