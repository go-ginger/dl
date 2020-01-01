package dl

import (
	"github.com/go-ginger/models"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) handleReadOnlyFields(request models.IRequest) {
	req := request.GetBaseRequest()
	s := reflect.ValueOf(req.Body).Elem()
	typeOfT := s.Type()
	switch s.Kind() {
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			ff := typeOfT.Field(i)
			tag, ok := ff.Tag.Lookup("dl")
			if ok {
				tagParts := strings.Split(tag, ",")
				for _, tagPart := range tagParts {
					switch tagPart {
					case "read_only":
						if f.IsValid() {
							if f.CanSet() {
								f.Set(reflect.Zero(ff.Type))
							}
						}
						break
					}
				}
			}
			tag, ok = ff.Tag.Lookup("edit_roles")
			if ok {
				canEdit := false
				auth := request.GetAuth()
				if auth != nil {
					tagParts := strings.Split(tag, ",")
					for _, role := range tagParts {
						if auth.HasRole(role) || (role == "id" && auth.GetCurrentAccountId() == req.ID) {
							canEdit = true
							break
						}
					}
				}
				if !canEdit {
					if f.IsValid() {
						if f.CanSet() {
							f.Set(reflect.Zero(ff.Type))
						}
					}
				}
			}
		}
		break
	}
}

func (base *BaseDbHandler) BeforeUpsert(request models.IRequest) (err error) {
	base.handleReadOnlyFields(request)
	req := request.GetBaseRequest()
	req.Body.HandleUpsertDefaultValues()
	return
}

func (base *BaseDbHandler) AfterUpsert(request models.IRequest) (err error) {
	if base.SecondaryDB != nil {
		err = base.SecondaryDB.Upsert(request)
	}
	return
}

func (base *BaseDbHandler) DoUpsert(request models.IRequest) (err error) {
	err = base.IBaseDbHandler.BeforeUpsert(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.Upsert(request)
	if err != nil {
		return
	}
	err = base.IBaseDbHandler.AfterUpsert(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Upsert(request models.IRequest) error {
	return nil
}
