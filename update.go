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
		item, e := base.IBaseDbHandler.Get(request)
		if e != nil {
			err = e
			return
		}
		request.SetBody(item)
	}
	err = secondaryDB.Update(request)
	return
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) (err error) {
	if base.SecondaryDBs != nil {
		for _, secondaryDB := range base.SecondaryDBs {
			if secondaryDB.UpdateInBackgroundEnabled() {
				go func() {
					err = base.handleSecondaryUpdate(request, secondaryDB)
					if err != nil {
						log.Println(fmt.Sprintf("error on handleSecondaryUpdate, err: %v", err))
						return
					}
				}()
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
	req := request.GetBaseRequest()
	req.AddTag("system", true)
	err = base.IBaseDbHandler.AfterUpdate(req)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Update(request models.IRequest) error {
	return nil
}
