package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
)

func (base *BaseDbHandler) BeforeInsert(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleCreateDefaultValues()
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
