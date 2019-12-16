package dl

import (
	"github.com/go-ginger/models"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) BeforeQuery(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	if req.Sort == nil || len(*req.Sort) == 0 {
		// default sort with id desc
		req.Sort = &[]models.SortItem{
			{
				Name:      "id",
				Ascending: false,
			},
		}
	}
	return
}

func (base *BaseDbHandler) handleModelAfterQuery(request models.IRequest, model interface{}, isValue bool) {
	var s reflect.Value
	if !isValue {
		s = reflect.ValueOf(model).Elem()
	} else {
		s = model.(reflect.Value)
	}
	typeOfT := s.Type()

	req := request.GetBaseRequest()
	if doLoad, ok := req.Tags["load"]; !ok || doLoad {
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fType := typeOfT.Field(i)
			tag, ok := fType.Tag.Lookup("load")
			if ok {
				tagParts := strings.Split(tag, ",")
				eventName, targetFieldName := tagParts[0], tagParts[1]
				val := f.Interface()
				result, handled := TryEvent(eventName, fType.Name, val)
				if handled {
					targetField := s.FieldByName(targetFieldName)
					if targetField.IsValid() {
						if f.CanSet() {
							targetField.Set(reflect.ValueOf(result))
						}
					}
				}
			}
		}
	}
}

func (base *BaseDbHandler) AfterQuery(request models.IRequest, result interface{}) (err error) {
	if pr, ok := result.(*models.PaginateResult); ok {
		items, ok := pr.Items.([]interface{})
		if ok {
			for _, item := range items {
				base.handleModelAfterQuery(request, item, false)
			}
		}
	} else {
		base.handleModelAfterQuery(request, result, false)
	}
	return
}

func (base *BaseDbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request models.IRequest) (models.IBaseModel, error) {
	return nil, nil
}
