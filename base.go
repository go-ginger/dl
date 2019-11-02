package dl

import (
	h "github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"reflect"
)

type IBaseDbHandler interface {
	Initialize(model interface{})
	GetModelsInstance() interface{}
	GetModelInstance() interface{}

	BeforeInsert(request models.IRequest)
	Insert(request models.IRequest) (interface{}, error)
	AfterInsert(request models.IRequest)

	BeforeQuery(request models.IRequest)
	Paginate(request models.IRequest) (*models.PaginateResult, error)
	Get(request models.IRequest) (models.IBaseModel, error)
	AfterQuery(request models.IRequest, result interface{})

	BeforeUpdate(request models.IRequest)
	Update(request models.IRequest) error
	AfterUpdate(request models.IRequest)

	BeforeDelete(request models.IRequest)
	Delete(request models.IRequest) error
	AfterDelete(request models.IRequest)
}

type BaseDbHandler struct {
	IBaseDbHandler

	Model     reflect.Value
	ModelType reflect.Type
}

func (base *BaseDbHandler) Initialize(model interface{}) {
	base.Model = reflect.ValueOf(model)
	base.ModelType = base.Model.Type()
}

func (base *BaseDbHandler) GetModelInstance() interface{} {
	return h.NewInstanceOfType(base.ModelType)
}

func (base *BaseDbHandler) GetModelsInstance() interface{} {
	return h.NewSliceInstanceOfType(base.ModelType)
}

func (base *BaseDbHandler) handleFilters(request models.IRequest) {
	request.AddNewFilter("deleted", map[string]bool{
		"$ne": true,
	})
}
