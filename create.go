package dl

import (
	"fmt"
	"github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeInsert(request models.IRequest) (err error) {
	base.HandleReadOnlyFields(request)
	base.IBaseDbHandler.HandleCreateDefaultValues(request)
	return
}

func (base *BaseDbHandler) HandleCreateDefaultValues(request models.IRequest) {
	request.GetBody().HandleCreateDefaultValues()
}

func (base *BaseDbHandler) DoInsert(request models.IRequest) (result interface{}, err error) {
	err = base.IBaseDbHandler.BeforeInsert(request)
	if err != nil {
		return
	}
	result, err = base.IBaseDbHandler.Insert(request)
	if err != nil {
		return
	}
	request.SetTag("system", true)
	err = base.IBaseDbHandler.AfterInsert(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Insert(request models.IRequest) (result interface{}, err error) {
	return
}

func (base *BaseDbHandler) AfterInsert(request models.IRequest) (err error) {
	if base.SecondaryDBs != nil {
		for _, secondaryDB := range base.SecondaryDBs {
			if secondaryDB.InsertInBackgroundEnabled() {
				go func(db IBaseDbHandler) {
					secondaryRequest := helpers.Clone(request).(models.IRequest)
					_, err := db.DoInsert(secondaryRequest)
					if err != nil {
						log.Println(fmt.Sprintf("Insert error on secondary dbHandler, err: %v", err))
						return
					}
				}(secondaryDB)
			} else {
				secondaryRequest := request.Populate(&models.Request{
					Body: request.GetBody(),
				})
				_, err = secondaryDB.DoInsert(secondaryRequest)
			}
		}
		if base.HasAnyDenormalizeConfig {
			go base.IBaseDbHandler.DenormalizeNew(request.GetID())
		}
		if base.DenormalizeFieldRefConfigs != nil {
			//go base.IBaseDbHandler.EnsureDenormalizeByID(request.GetID())
		}
	}
	return
}
