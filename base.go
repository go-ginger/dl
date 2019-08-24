package dl

import "github.com/kulichak/models"

type IBaseDbHandler interface {
	BeforeInsert(request models.IRequest)
	Insert(request models.IRequest) (interface{}, error)
	AfterInsert(request models.IRequest)

	BeforeQuery(request models.IRequest)
	Paginate(request models.IRequest) (*models.PaginateResult, error)
	Get(request models.IRequest) (models.IBaseModel, error)
	AfterQuery(request models.IRequest)

	BeforeUpdate(request models.IRequest)
	Update(request models.IRequest) error
	AfterUpdate(request models.IRequest)

	BeforeDelete(request models.IRequest)
	Delete(request models.IRequest) error
	AfterDelete(request models.IRequest)
}

type BaseDbHandler struct {
	IBaseDbHandler
}

func (base *BaseDbHandler) handleFilters(request models.IRequest) {
	request.AddNewFilter("deleted", map[string]bool{
		"$ne": true,
	})
}
