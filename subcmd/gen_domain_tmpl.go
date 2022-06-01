package subcmd

const serviceTmpl = `package {{.domain}}

import (
	"context"
	"{{.projectName}}/internal/domain"
)

type Service struct {
	ctx  *domain.ServiceContext
	repo *repo
}

type QueryParam struct {
}

func NewService(ctx *domain.ServiceContext) *Service {
	return &Service{
		ctx:  ctx,
		repo: newRepo(),
	}
}

func (s *Service) Create(ctx context.Context, {{.domain}} *{{.entity}}) error {
	if err := s.repo.Create(ctx, {{.domain}}); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, {{.domain}} *{{.entity}}) error {
	if err := s.repo.Update(ctx, {{.domain}}); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, {{.domain}} *{{.entity}}) error {
	if err := s.repo.Delete(ctx, {{.domain}}); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id int) ({{.entity}}, error) {
	result, err := s.repo.Get(ctx, id)
	if err != nil {
		return {{.entity}}{}, err
	}

	return result, nil
}

func (s *Service) Query(ctx context.Context, param QueryParam) ([]{{.entity}}, int64, error) {
	result, total, err := s.repo.Query(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}
`

const entityTmpl = `package {{.domain}}

import "time"

type {{.entity}} struct {
	Id         int
	CreateTime time.Time
	UpdateTime time.Time
}
`

const initTmpl = `package {{.domain}}

func init() {
}
`

const repoTmpl = `package {{.domain}}

import (
	"context"
)

type repo struct {
}

func newRepo() *repo {
	return &repo{}
}

func (r *repo) Create(ctx context.Context, {{.domain}} *{{.entity}}) error {
	return nil
}

func (r *repo) Update(ctx context.Context, {{.domain}} *{{.entity}}) error {
	return nil
}

func (r *repo) Delete(ctx context.Context, {{.domain}} *{{.entity}}) error {
	return nil
}

func (r *repo) Get(ctx context.Context, id int) ({{.entity}}, error) {
	return {{.entity}}{}, nil
}

func (r *repo) Query(ctx context.Context, param QueryParam) ([]{{.entity}}, int64, error) {
	result := make([]{{.entity}}, 0)

	total := int64(0)

	return result, total, nil
}
`
