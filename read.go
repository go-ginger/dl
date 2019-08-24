package dl

import (
	"github.com/kulichak/models"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) BeforeQuery(request models.IRequest) {
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

func (base *BaseDbHandler) AfterQuery(request models.IRequest) {
	req := request.GetBaseRequest()
	if req.Model != nil {
		base.handleModelAfterQuery(request, req.Model, false)
	} else if req.Models != nil {
		s := reflect.ValueOf(req.Models).Elem()
		for i := 0; i < s.Len(); i++ {
			base.handleModelAfterQuery(request, s.Index(i), true)
		}
	}
}

func (base *BaseDbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request models.IRequest) (models.IBaseModel, error) {
	return nil, nil
}
