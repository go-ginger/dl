package dl

import (
	"fmt"
	"github.com/kulichak/models"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) BeforeQuery(request *models.Request) {
}

func (base *BaseDbHandler) handleModelAfterQuery(model interface{}) {
	s := reflect.ValueOf(model).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fType := typeOfT.Field(i)
		tag, ok := fType.Tag.Lookup("load")
		if ok {
			tagParts := strings.Split(tag, ",")
			eventName, targetFieldName := tagParts[0], tagParts[1]
			fmt.Println(eventName)
			fmt.Println(targetFieldName)
			val := f.String()
			result, ok := TryEvent(eventName, fType.Name, val)
			if ok {
				targetField := s.FieldByName(targetFieldName)
				if targetField.IsValid() {
					if f.CanSet() {
						f.Set(reflect.ValueOf(result))
					}
				}
			}
		}
	}
}

func (base *BaseDbHandler) AfterQuery(request *models.Request) {
	if request.Model != nil {
		base.handleModelAfterQuery(request.Model)
	} else if request.Models != nil {
		for _, model := range request.Models.([]interface{}) {
			base.handleModelAfterQuery(model)
		}
	}
}

func (base *BaseDbHandler) Paginate(request *models.Request) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request *models.Request) (*models.IBaseModel, error) {
	return nil, nil
}
