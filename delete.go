package dl

import "github.com/kulichak/models"


func (base *BaseDbHandler) BeforeDelete(request models.IRequest) {
	req := request.GetBaseRequest()
	req.Body.HandleDeleteDefaultValues()
}

func (base *BaseDbHandler) AfterDelete(request models.IRequest) {
}

func (base *BaseDbHandler) Delete(request models.IRequest) error {
	return nil
}
