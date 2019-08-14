package dl

import "github.com/kulichak/models"

type IBaseData interface {
	Paginate(request *models.IRequest) *models.PaginateResult
}

type BaseData struct {
}

func (base *BaseData) Paginate(request *models.IRequest) *models.PaginateResult {
	(*request).AddNewFilter("deleted", map[string]bool{
		"$ne": true,
	})
	return nil
}
