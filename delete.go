package dl

import "github.com/go-ginger/models"

func (base *BaseDbHandler) BeforeDelete(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	if req.Body != nil {
		req.Body.HandleDeleteDefaultValues()
	}
	return
}

func (base *BaseDbHandler) AfterDelete(request models.IRequest) (err error) {
	return
}

func (base *BaseDbHandler) Delete(request models.IRequest) error {
	return nil
}
