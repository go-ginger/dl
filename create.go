package dl

import "github.com/kulichak/models"


func (base *BaseDbHandler) BeforeInsert(request *models.Request) {
}

func (base *BaseDbHandler) Insert(request *models.Request) (*models.IBaseModel, error) {
	return nil, nil
}

func (base *BaseDbHandler) AfterInsert(request *models.Request) {
}
