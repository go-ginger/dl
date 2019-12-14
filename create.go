package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeInsert(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	req.Body.HandleCreateDefaultValues()
	return
}

func (base *BaseDbHandler) Insert(request models.IRequest) (*models.IBaseModel, error) {
	return nil, nil
}

func (base *BaseDbHandler) AfterInsert(request models.IRequest) (err error) {
	return
}
