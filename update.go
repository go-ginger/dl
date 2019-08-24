package dl

import "github.com/kulichak/models"

func (base *BaseDbHandler) BeforeUpdate(request models.IRequest) {
	req := request.GetBaseRequest()
	req.Body.HandleUpdateDefaultValues()
}

func (base *BaseDbHandler) AfterUpdate(request models.IRequest) {
}

func (base *BaseDbHandler) Update(request models.IRequest) error {
	return nil
}
