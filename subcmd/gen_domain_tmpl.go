package subcmd

const useCaseTmpl = `package {{.domain}}

import (
	"context"
	"{{.projectName}}/internal/domain"
)

type UseCase struct {
	ctx  *domain.UseCaseContext
	repo *repo
}

type QueryParam struct {
}

func NewUseCase(ctx *domain.UseCaseContext) *UseCase {
	return &UseCase{
		ctx:  ctx,
		repo: newRepo(),
	}
}

func (uc *UseCase) Create(ctx context.Context, {{.domain}} *{{.entity}}) error {
	if err := uc.repo.Create(ctx, {{.domain}}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, {{.domain}} *{{.entity}}) error {
	if err := uc.repo.Update(ctx, {{.domain}}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, {{.domain}}); err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) Get(ctx context.Context, id string) ({{.entity}}, error) {
	result, err := uc.repo.Get(ctx, string)
	if err != nil {
		return {{.entity}}{}, err
	}

	return result, nil
}

func (uc *UseCase) Query(ctx context.Context, param QueryParam) ([]{{.entity}}, int64, error) {
	result, total, err := uc.repo.Query(ctx, param)
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

func (r *repo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *repo) Get(ctx context.Context, id string) ({{.entity}}, error) {
	return {{.entity}}{}, nil
}

func (r *repo) Query(ctx context.Context, param QueryParam) ([]{{.entity}}, int64, error) {
	result := make([]{{.entity}}, 0)

	total := int64(0)

	return result, total, nil
}
`
