package dl

import (
	"fmt"
	"github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"log"
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

func (base *BaseDbHandler) handleSecondaryUpdate(request models.IRequest, secondaryDB IBaseDbHandler) (err error) {
	secondaryRequest := helpers.Clone(request).(models.IRequest)
	if secondaryDB.IsFullObjOnUpdateRequired() {
		objID := request.GetID()
		req := secondaryRequest.GetBaseRequest()
		req.Filters = &models.Filters{
			"id": objID,
		}
		item, e := base.IBaseDbHandler.DoGet(secondaryRequest)
		if e != nil {
			err = e
			return
		}
		secondaryRequest.SetBody(item)
	}
	err = secondaryDB.DoUpdate(secondaryRequest)
	return
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) (err error) {
	if base.SecondaryDBs != nil {
		for _, secondaryDB := range base.SecondaryDBs {
			if secondaryDB.UpdateInBackgroundEnabled() {
				go func(db IBaseDbHandler) {
					err = base.handleSecondaryUpdate(request, db)
					if err != nil {
						log.Println(fmt.Sprintf("error on handleSecondaryUpdate, err: %v", err))
						return
					}
				}(secondaryDB)
			} else {
				err = base.handleSecondaryUpdate(request, secondaryDB)
			}
		}
	}
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
