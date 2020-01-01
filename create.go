package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeInsert(request models.IRequest) (err error) {
	base.handleReadOnlyFields(request)
	req := request.GetBaseRequest()
	req.Body.HandleCreateDefaultValues()
	return
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
	if base.SecondaryDB != nil {
		if base.SecondaryDB.InsertInBackgroundEnabled() {
			go func() {
				_, err := base.SecondaryDB.Insert(request)
				if err != nil {
					log.Println(fmt.Sprintf("Insert error on secondary dbHandler, err: %v", err))
					return
				}
			}()
		} else {
			_, err = base.SecondaryDB.Insert(request)
		}
	}
	return
}
