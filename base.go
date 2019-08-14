package dl

import "github.com/kulichak/models"

type IBaseDbHandler interface {
	Paginate(request *models.IRequest) (*models.PaginateResult, error)
}

type BaseDbHandler struct {
}

func (base *BaseDbHandler) Paginate(request *models.IRequest) (*models.PaginateResult, error) {
	(*request).AddNewFilter("deleted", map[string]bool{
		"$ne": true,
	})
	return nil, nil
}
