package dl

import (
	"github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"reflect"
	"strings"
	"sync"
)

func (base *BaseDbHandler) BeforeQuery(request models.IRequest) (err error) {
	return
}

const (
	CallerGet      = "GET"
	CallerPaginate = "PAGINATE"
)

func (base *BaseDbHandler) performEvent(request models.IRequest, eventName, fieldName string, value interface{},
	field *reflect.Value, wg *sync.WaitGroup) {
	defer wg.Done()
	result, handled := TryEvent(request, eventName, fieldName, value)
	if handled && result != nil {
		if field.CanSet() {
			field.Set(reflect.ValueOf(result))
		}
	}
}

func (base *BaseDbHandler) beginHandleModelAfterQuery(request models.IRequest, model interface{},
	isValue bool, remainingDepth int, wg *sync.WaitGroup) {
	defer wg.Done()
	base.handleModelAfterQuery(request, model, isValue, remainingDepth)
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
	sType := s.Type()

	req := request.GetBaseRequest()
	if doLoad, ok := req.Tags["load"]; !ok || doLoad {
		callerType := request.GetTemp("caller_type")
		switch s.Kind() {
		case reflect.Struct:
			wg := sync.WaitGroup{}
			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)
				ff := sType.Field(i)
				tag, ok := ff.Tag.Lookup("load_from")
				if !ok {
					if callerType == CallerGet {
						tag, ok = ff.Tag.Lookup("get_load_from")
					}
					if !ok {
						if callerType == CallerPaginate {
							tag, ok = ff.Tag.Lookup("paginate_load_from")
						}
					}
				}
				if ok {
					tagParts := strings.Split(tag, ",")
					sourceFieldName, eventName := tagParts[0], tagParts[1]
					sourceField := s.FieldByName(sourceFieldName)
					val := sourceField.Interface()
					wg.Add(1)
					go base.performEvent(request, eventName, ff.Name, val, &f, &wg)
				}
				if helpers.IsEmptyValue(f) {
					continue
				}
				tag, ok = ff.Tag.Lookup("load")
				if !ok {
					if callerType == CallerGet {
						tag, ok = ff.Tag.Lookup("get_load")
					}
					if !ok {
						if callerType == CallerPaginate {
							tag, ok = ff.Tag.Lookup("paginate_load")
						}
					}
				}
				if ok {
					tagParts := strings.Split(tag, ",")
					eventName, targetFieldName := tagParts[0], tagParts[1]
					val := f.Interface()
					targetField := s.FieldByName(targetFieldName)
					if f.CanSet() {
						wg.Add(1)
						go base.performEvent(request, eventName, ff.Name, val, &targetField, &wg)
					}
				}
				switch f.Type().Kind() {
				case reflect.Ptr:
					if f.IsNil() {
						break
					}
					i := f.Elem()
					wg.Add(1)
					go base.beginHandleModelAfterQuery(request, i, true, remainingDepth-1, &wg)
					break
				case reflect.Struct:
					wg.Add(1)
					go base.beginHandleModelAfterQuery(request, f, true, remainingDepth-1, &wg)
					break
				}
			}
			addr := s.Addr()
			if addr.IsValid() && addr.CanInterface() {
				mv := addr.Interface()
				if baseModel, ok := mv.(models.IBaseModel); ok {
					baseModel.Populate(request)
				}
			}
			wg.Wait()
			break
		case reflect.Slice:
			wg := sync.WaitGroup{}
			for i := 0; i < s.Len(); i++ {
				wg.Add(1)
				go func(item reflect.Value) {
					defer wg.Done()
					base.handleModelAfterQuery(request, item, true, 3)
				}(s.Index(i))
			}
			wg.Wait()
			break
		}
	}
}

func (base *BaseDbHandler) AfterQuery(request models.IRequest, result interface{}) (err error) {
	if pr, ok := result.(*models.PaginateResult); ok {
		switch reflect.TypeOf(pr.Items).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(pr.Items)
			wg := sync.WaitGroup{}
			for i := 0; i < s.Len(); i++ {
				wg.Add(1)
				go func(item reflect.Value) {
					defer wg.Done()
					base.handleModelAfterQuery(request, item, true, 3)
				}(s.Index(i))
			}
			wg.Wait()
		}
	} else {
		base.handleModelAfterQuery(request, result, false, 3)
	}
	return
}

func (base *BaseDbHandler) DoPaginate(request models.IRequest) (result *models.PaginateResult, err error) {
	request.SetTemp("caller_type", CallerPaginate)
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
	request.SetTemp("caller_type", CallerGet)
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

func (base *BaseDbHandler) DoGetFirst(request models.IRequest) (result models.IBaseModel, err error) {
	err = base.IBaseDbHandler.BeforeQuery(request)
	if err != nil {
		return
	}
	result, err = base.IBaseDbHandler.First(request)
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

func (base *BaseDbHandler) First(request models.IRequest) (result models.IBaseModel, err error) {
	return
}
