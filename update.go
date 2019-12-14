package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeUpdate(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleUpdateDefaultValues()
	return
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) (err error) {
	return
}

func (base *BaseDbHandler) Update(request models.IRequest) error {
	return nil
}
