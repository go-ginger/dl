package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeUpdate(request models.IRequest) (err error) {
	base.handleReadOnlyFields(request)
	req := request.GetBaseRequest()
	if req.Body != nil {
		req.Body.HandleUpdateDefaultValues()
	}
	return
}

func (base *BaseDbHandler) handleSecondaryUpdate(request models.IRequest, secondaryDB IBaseDbHandler) (err error) {
	if secondaryDB.IsFullObjOnUpdateRequired() {
		objID := request.GetID()
		req := request.GetBaseRequest()
		req.Filters = &models.Filters{
			"id": objID,
		}
		item, e := base.IBaseDbHandler.DoGet(request)
		if e != nil {
			err = e
			return
		}
		request.SetBody(item)
	}
	err = secondaryDB.DoUpdate(request)
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
