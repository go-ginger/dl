package dl

import "github.com/kulichak/models"

func (base *BaseDbHandler) BeforeQuery(request *models.Request) {
}

func (base *BaseDbHandler) AfterQuery(request *models.Request) {
}

func (base *BaseDbHandler) Paginate(request *models.Request) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request *models.Request) (*models.IBaseModel, error) {
	return nil, nil
}
