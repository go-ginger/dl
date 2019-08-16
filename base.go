package dl

import "github.com/kulichak/models"

type IBaseDbHandler interface {
	BeforeInsert(request *models.Request)
	Insert(request *models.Request) (interface{}, error)
	AfterInsert(request *models.Request)

	BeforeQuery(request *models.Request)
	Paginate(request *models.Request) (*models.PaginateResult, error)
	Get(request *models.Request) (*models.IBaseModel, error)
	AfterQuery(request *models.Request)

	BeforeUpdate(request *models.Request)
	Update(request *models.Request) error
	AfterUpdate(request *models.Request)

	BeforeDelete(request *models.Request)
	Delete(request *models.Request) error
	AfterDelete(request *models.Request)
}

type BaseDbHandler struct {
	IBaseDbHandler
}

func (base *BaseDbHandler) handleFilters(request *models.Request) {
	request.AddNewFilter("deleted", map[string]bool{
		"$ne": true,
	})
}
