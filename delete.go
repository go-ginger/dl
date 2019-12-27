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
	if base.SecondaryDB != nil {
		if base.SecondaryDB.DeleteInBackgroundEnabled() {
			go func() {
				err := base.SecondaryDB.Delete(request)
				if err != nil {
					log.Println(fmt.Sprintf("error on delete secondary dbHandler, err: %v", err))
					return
				}
			}()
		} else {
			err = base.SecondaryDB.Delete(request)
		}
	}
	return
}

func (base *BaseDbHandler) Delete(request models.IRequest) (err error) {
	return
}
