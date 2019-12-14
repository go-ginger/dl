package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeUpsert(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleUpsertDefaultValues()
	return
}

func (base *BaseDbHandler) AfterUpsert(request models.IRequest) (err error) {
	return
}

func (base *BaseDbHandler) Upsert(request models.IRequest) error {
	return nil
}
