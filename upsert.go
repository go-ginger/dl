package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeUpsert(request models.IRequest) {
	req := request.GetBaseRequest()
	req.Body.HandleUpdateDefaultValues()
}

func (base *BaseDbHandler) AfterUpsert(request models.IRequest) {
}

func (base *BaseDbHandler) Upsert(request models.IRequest) error {
	return nil
}
