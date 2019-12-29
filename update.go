package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeUpdate(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleUpdateDefaultValues()
	return
}

func (base *BaseDbHandler) handleSecondaryUpdate(request models.IRequest) (err error) {
	if base.SecondaryDB.IsFullObjOnUpdateRequired() {
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
	err = base.SecondaryDB.Update(request)
	return
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) (err error) {
	if base.SecondaryDB != nil {
		if base.SecondaryDB.UpdateInBackgroundEnabled() {
			go func() {
				err = base.handleSecondaryUpdate(request)
				if err != nil {
					log.Println(fmt.Sprintf("error on handleSecondaryUpdate, err: %v", err))
					return
				}
			}()
		} else {
			err = base.handleSecondaryUpdate(request)
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
	err = base.IBaseDbHandler.AfterUpdate(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Update(request models.IRequest) error {
	return nil
}
