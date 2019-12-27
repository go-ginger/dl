package dl

import (
	h "github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"reflect"
)

type IBaseDbHandler interface {
	Initialize(handler IBaseDbHandler, model interface{})
	GetModelsInstance() interface{}
	GetModelInstance() interface{}

	BeforeInsert(request models.IRequest) (err error)
	Insert(request models.IRequest) (interface{}, error)
	AfterInsert(request models.IRequest) (err error)

	BeforeQuery(request models.IRequest) (err error)
	Paginate(request models.IRequest) (*models.PaginateResult, error)
	Get(request models.IRequest) (models.IBaseModel, error)
	AfterQuery(request models.IRequest, result interface{}) (err error)

	BeforeUpdate(request models.IRequest) (err error)
	Update(request models.IRequest) error
	AfterUpdate(request models.IRequest) (err error)

	BeforeUpsert(request models.IRequest) (err error)
	Upsert(request models.IRequest) error
	AfterUpsert(request models.IRequest) (err error)

	BeforeDelete(request models.IRequest) (err error)
	Delete(request models.IRequest) error
	AfterDelete(request models.IRequest) (err error)

	InsertInBackgroundEnabled() bool
	UpdateInBackgroundEnabled() bool
	DeleteInBackgroundEnabled() bool
	IsFullObjOnUpdateRequired() bool
}

type BaseDbHandler struct {
	IBaseDbHandler

	Model       reflect.Value
	ModelType   reflect.Type
	SecondaryDB IBaseDbHandler

	InsertInBackground           bool
	UpdateInBackground           bool
	DeleteInBackground           bool
	IsFullObjectOnUpdateRequired bool
}

func (base *BaseDbHandler) Initialize(handler IBaseDbHandler, model interface{}) {
	if model != nil {
		base.Model = reflect.ValueOf(model)
		base.ModelType = base.Model.Type()
	}
	base.IBaseDbHandler = handler
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

func (base *BaseDbHandler) InsertInBackgroundEnabled() bool {
	return base.InsertInBackground
}

func (base *BaseDbHandler) UpdateInBackgroundEnabled() bool {
	return base.UpdateInBackground
}

func (base *BaseDbHandler) DeleteInBackgroundEnabled() bool {
	return base.DeleteInBackground
}

func (base *BaseDbHandler) IsFullObjOnUpdateRequired() bool {
	return base.IsFullObjectOnUpdateRequired
}
