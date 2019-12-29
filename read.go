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

func (base *BaseDbHandler) handleModelAfterQuery(request models.IRequest, model interface{},
	isValue bool, remainingDepth int) {
	if remainingDepth == 0 {
		return
	}
	var s reflect.Value
	if !isValue {
		s = reflect.ValueOf(model).Elem()
	} else {
		s = model.(reflect.Value)
	}
	typeOfT := s.Type()

	req := request.GetBaseRequest()
	if doLoad, ok := req.Tags["load"]; !ok || doLoad {
		switch s.Kind() {
		case reflect.Struct:
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				fType := typeOfT.Field(i)
				switch f.Type().Kind() {
				case reflect.Ptr:
					if f.IsNil() {
						continue
					}
					i := f.Elem()
					base.handleModelAfterQuery(request, i, true, remainingDepth-1)
					continue
				case reflect.Struct:
					base.handleModelAfterQuery(request, f, true, remainingDepth-1)
					continue
				}
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
			break
		}
	}
}

func (base *BaseDbHandler) AfterQuery(request models.IRequest, result interface{}) (err error) {
	if pr, ok := result.(*models.PaginateResult); ok {
		switch reflect.TypeOf(pr.Items).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(pr.Items)
			for i := 0; i < s.Len(); i++ {
				base.handleModelAfterQuery(request, s.Index(i), true, 3)
			}
		}
	} else {
		base.handleModelAfterQuery(request, result, false, 3)
	}
	return
}

func (base *BaseDbHandler) DoPaginate(request models.IRequest) (result *models.PaginateResult, err error) {
	err = base.IBaseDbHandler.BeforeQuery(request)
	if err != nil {
		return
	}
	result, err = base.IBaseDbHandler.Paginate(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.AfterQuery(request, result)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) DoGet(request models.IRequest) (result models.IBaseModel, err error) {
	err = base.IBaseDbHandler.BeforeQuery(request)
	if err != nil {
		return
	}
	result, err = base.IBaseDbHandler.Get(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.AfterQuery(request, result)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	return nil, nil
}

func (base *BaseDbHandler) Get(request models.IRequest) (models.IBaseModel, error) {
	return nil, nil
}
