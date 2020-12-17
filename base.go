package dl

import (
	"github.com/go-ginger/dl/query"
	h "github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"reflect"
)

type IBaseDbHandler interface {
	Initialize(handler IBaseDbHandler, model interface{})
	GetModelsInstance() interface{}
	GetModelInstance() interface{}
	GetBaseDbHandler() IBaseDbHandler

	DoInsert(request models.IRequest) (result models.IBaseModel, err error)
	DoPaginate(request models.IRequest) (*models.PaginateResult, error)
	DoGet(request models.IRequest) (models.IBaseModel, error)
	DoGetFirst(request models.IRequest) (models.IBaseModel, error)
	DoUpdate(request models.IRequest) error
	DoUpsert(request models.IRequest) error
	DoDelete(request models.IRequest) error

	HandleCreateDefaultValues(request models.IRequest)
	HandleUpdateDefaultValues(request models.IRequest)

	BeforeInsert(request models.IRequest) (err error)
	Insert(request models.IRequest) (result models.IBaseModel, err error)
	AfterInsert(request models.IRequest) (err error)

	BeforeQuery(request models.IRequest) (err error)
	Paginate(request models.IRequest) (*models.PaginateResult, error)
	Get(request models.IRequest) (models.IBaseModel, error)
	First(request models.IRequest) (result models.IBaseModel, err error)
	Exists(request models.IRequest) (exists bool, err error)
	Count(request models.IRequest) (count uint64, err error)
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

	StartTransaction(request models.IRequest) (err error)
	CommitTransaction(request models.IRequest) (err error)
	RollbackTransaction(request models.IRequest) (err error)

	InsertInBackgroundEnabled() bool
	UpdateInBackgroundEnabled() bool
	DeleteInBackgroundEnabled() bool
	IsFullObjOnUpdateRequired() bool

	IdEquals(id1 interface{}, id2 interface{}) bool

	EnsureDenormalizeByID(id interface{})
	NewDenormalizeConfig(configs ...*DenormalizeConfig)
	NewDenormalizeReferenceConfig(configs ...*DenormalizeConfig)
	DenormalizeNew(id interface{})
	DenormalizeDelete(id interface{})
	DenormalizeNewEntity(entityValue reflect.Value, newEntityValue reflect.Value, info *DenormalizeConfig)
	DenormalizeDeleteEntity(entityValue reflect.Value, deletedEntityID interface{}, info *DenormalizeConfig)
	EnsureDenormalizeInterface(id, entity interface{})

	GetSecondaryDBs() []IBaseDbHandler
}

type BaseDbHandler struct {
	IBaseDbHandler

	QueryParser query.IParser

	Model        reflect.Value
	ModelType    reflect.Type
	SecondaryDBs []IBaseDbHandler

	SetFlagOnDelete              *bool
	InsertInBackground           bool
	UpdateInBackground           bool
	DeleteInBackground           bool
	IsFullObjectOnUpdateRequired bool

	HasAnyDenormalizeConfig    bool
	DenormalizeConfigs         []*DenormalizeConfig
	DenormalizeFieldRefConfigs []*DenormalizeConfig // map for referenced fields
}

func (base *BaseDbHandler) Initialize(handler IBaseDbHandler, model interface{}) {
	if model != nil {
		base.Model = reflect.ValueOf(model)
		base.ModelType = base.Model.Type()
	}
	if base.SetFlagOnDelete == nil {
		setFlag := true
		base.SetFlagOnDelete = &setFlag
	}
	base.IBaseDbHandler = handler
}

func (base *BaseDbHandler) GetBaseDbHandler() IBaseDbHandler {
	return base
}

func (base *BaseDbHandler) GetModelInstance() interface{} {
	return h.NewInstanceOfType(base.ModelType)
}

func (base *BaseDbHandler) GetModelsInstance() interface{} {
	return h.NewSliceInstanceOfType(base.ModelType)
}

func (base *BaseDbHandler) GetModelsInstancePtr() interface{} {
	return h.NewSliceInstanceOfTypePtr(base.ModelType)
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

func (base *BaseDbHandler) IdEquals(id1 interface{}, id2 interface{}) bool {
	return id1 == id2
}

func (base *BaseDbHandler) GetSecondaryDBs() []IBaseDbHandler {
	return base.SecondaryDBs
}
