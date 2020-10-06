package dl

import (
	"github.com/go-ginger/models"
)

func (base *BaseDbHandler) BeforeUpdate(request models.IRequest) (err error) {
	base.HandleReadOnlyFields(request)
	base.IBaseDbHandler.HandleUpdateDefaultValues(request)
	return
}

func (base *BaseDbHandler) HandleUpdateDefaultValues(request models.IRequest) {
	body := request.GetBody()
	if body != nil {
		body.HandleUpdateDefaultValues()
	}
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) (err error) {
	if base.HasAnyDenormalizeConfig {
		go base.IBaseDbHandler.DenormalizeNew(request.GetID())
	}
	if base.DenormalizeFieldRefConfigs != nil {
		go base.IBaseDbHandler.EnsureDenormalizeByID(request.GetID())
	}
	return
}

func (base *BaseDbHandler) DoUpdate(request models.IRequest) (err error) {
	err = base.IBaseDbHandler.BeforeUpdate(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.Update(request)
	if err != nil {
		return
	}
	request.SetTag("system", true)
	err = base.IBaseDbHandler.AfterUpdate(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Update(request models.IRequest) error {
	return nil
}
