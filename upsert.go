package dl

import (
	"fmt"
	"github.com/go-ginger/helpers"
	"github.com/go-ginger/models"
	"log"
	"reflect"
	"strings"
)

func (base *BaseDbHandler) HandleReadOnlyFields(request models.IRequest) {
	req := request.GetBaseRequest()
	if req.Body == nil {
		return
	}
	if isSystem, ok := req.Tags["system"]; ok && isSystem {
		return
	}
	if checkReadOnly, ok := req.Tags["check_edit_roles"]; ok && !checkReadOnly {
		return
	}
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
						if auth.HasRole(role) || (role == "id" && auth.GetCurrentAccountId(request) == req.ID) {
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
	base.HandleReadOnlyFields(request)
	req := request.GetBaseRequest()
	req.Body.HandleUpsertDefaultValues()
	return
}

func (base *BaseDbHandler) handleSecondaryUpsert(request models.IRequest, secondaryDB IBaseDbHandler) (err error) {
	secondaryRequest := helpers.Clone(request).(models.IRequest)
	if secondaryDB.IsFullObjOnUpdateRequired() {
		item, e := base.IBaseDbHandler.DoGet(secondaryRequest)
		if e != nil {
			err = e
			return
		}
		secondaryRequest.SetID(item.GetID())
		secondaryRequest.SetBody(item)
	}
	err = secondaryDB.DoUpsert(secondaryRequest)
	return
}

func (base *BaseDbHandler) AfterUpsert(request models.IRequest) (err error) {
	if base.SecondaryDBs != nil {
		for _, secondaryDB := range base.SecondaryDBs {
			if base.InsertInBackgroundEnabled() && base.UpdateInBackgroundEnabled() {
				go func(db IBaseDbHandler) {
					err := base.handleSecondaryUpsert(request, db)
					if err != nil {
						log.Println(fmt.Sprintf("error upsert secondary db %v, err: %v", secondaryDB, err))
					}
				}(secondaryDB)
			} else {
				e := base.handleSecondaryUpsert(request, secondaryDB)
				if e != nil {
					err = e
				}
			}
		}
		if base.HasAnyDenormalizeConfig {
			go base.IBaseDbHandler.DenormalizeNew(request.GetID())
		}
		if base.DenormalizeFieldRefConfigs != nil {
			go base.IBaseDbHandler.EnsureDenormalizeByID(request.GetID())
		}
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
	request.SetTag("system", true)
	err = base.IBaseDbHandler.AfterUpsert(request)
	if err != nil {
		return
	}
	return
}

func (base *BaseDbHandler) Upsert(request models.IRequest) error {
	return nil
}
