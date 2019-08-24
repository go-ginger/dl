package dl

import "github.com/kulichak/models"


func (base *BaseDbHandler) BeforeInsert(request models.IRequest) {
	req := request.GetBaseRequest()
	req.Body.HandleCreateDefaultValues()
}

func (base *BaseDbHandler) Insert(request models.IRequest) (*models.IBaseModel, error) {
	return nil, nil
}

func (base *BaseDbHandler) AfterInsert(request models.IRequest) {
}
