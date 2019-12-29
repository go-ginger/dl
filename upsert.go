package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeUpsert(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleUpsertDefaultValues()
	return
}

func (base *BaseDbHandler) AfterUpsert(request models.IRequest) (err error) {
	if base.SecondaryDB != nil {
		err = base.SecondaryDB.Upsert(request)
	}
	return
}

func (base *BaseDbHandler) DoUpsert(request models.IRequest) (err error) {
	err = base.IBaseDbHandler.BeforeUpsert(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.Upsert(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.AfterUpsert(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Upsert(request models.IRequest) error {
	return nil
}
