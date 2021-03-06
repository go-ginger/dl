package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeDelete(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	if req.Body != nil {
		req.Body.HandleDeleteDefaultValues()
	}
	return
}

func (base *BaseDbHandler) AfterDelete(request models.IRequest) (err error) {
	if base.SecondaryDBs != nil {
		for _, secondaryDB := range base.SecondaryDBs {
			if secondaryDB.DeleteInBackgroundEnabled() {
				go func(db IBaseDbHandler) {
					err := db.DoDelete(request)
					if err != nil {
						log.Println(fmt.Sprintf("error on delete secondary dbHandler, err: %v", err))
						return
					}
				}(secondaryDB)
			} else {
				err = secondaryDB.DoDelete(request)
			}
		}
		if base.HasAnyDenormalizeConfig {
			go base.IBaseDbHandler.DenormalizeNew(request.GetID())
		}
	}
	return
}

func (base *BaseDbHandler) DoDelete(request models.IRequest) (err error) {
	err = base.IBaseDbHandler.BeforeDelete(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.Delete(request)
	if err != nil {
		return
	}
	request.SetTag("system", true)
	err = base.IBaseDbHandler.AfterDelete(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Delete(request models.IRequest) (err error) {
	return
}
