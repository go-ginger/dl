package dl

import (
	"github.com/kulichak/models"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) BeforeQuery(request *models.Request) {
	if request.Sort == nil || len(*request.Sort) == 0 {
		// default sort with id desc
		request.Sort = &[]models.SortItem{
			{
				Name:      "id",
				Ascending: false,
			},
		}
	}
}

func (base *BaseDbHandler) handleModelAfterQuery(model interface{}, isValue bool) {
	var s reflect.Value
	if !isValue {
		s = reflect.ValueOf(model).Elem()
	} else {
		s = model.(reflect.Value)
	}
	typeOfT := s.Type()

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

func (base *BaseDbHandler) AfterQuery(request *models.Request) {
	if request.Model != nil {
		base.handleModelAfterQuery(request.Model, false)
	} else if request.Models != nil {
		s := reflect.ValueOf(request.Models).Elem()
		for i := 0; i < s.Len(); i++ {
			base.handleModelAfterQuery(s.Index(i), true)
		}
	}
}

func (base *BaseDbHandler) Paginate(request *models.Request) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request *models.Request) (*models.IBaseModel, error) {
	return nil, nil
}
