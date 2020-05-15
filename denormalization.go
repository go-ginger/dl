package dl

import (
	"fmt"
	"github.com/go-ginger/models"
	"log"
	"reflect"
)

type DenormalizeConfig struct {
	TargetHandler      IBaseDbHandler
	TargetFieldName    string
	TargetIdFilter     string
	ReferenceFieldName string
}

func (base *BaseDbHandler) NewDenormalizeConfig(configs ...*DenormalizeConfig) {
	if base.DenormalizeConfigs == nil {
		base.DenormalizeConfigs = []*DenormalizeConfig{}
	}
	base.IBaseDbHandler.NewDenormalizeReferenceConfig(configs...)
	base.DenormalizeConfigs = append(base.DenormalizeConfigs, configs...)
	base.HasAnyDenormalizeConfig = true
}

func (base *BaseDbHandler) NewDenormalizeReferenceConfig(configs ...*DenormalizeConfig) {
	for _, target := range configs {
		if target.ReferenceFieldName != "" {
			if base.DenormalizeFieldRefConfigs == nil {
				base.DenormalizeFieldRefConfigs = []*DenormalizeConfig{}
			}
			base.DenormalizeFieldRefConfigs = append(base.DenormalizeFieldRefConfigs, target)
		}
	}
}

func (base *BaseDbHandler) DenormalizeNewEntity(entityValue reflect.Value, newEntityValue reflect.Value,
	info *DenormalizeConfig) {
	newEntityPtrValue := newEntityValue.Addr()
	newEntityID := newEntityPtrValue.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
	targetField := entityValue.FieldByName(info.TargetFieldName)
	if targetField.Kind() == reflect.Ptr {
		targetField = targetField.Elem()
	}
	if targetField.Kind() == reflect.Slice {
		found := false
		for i := 0; i < targetField.Len(); i++ {
			item := targetField.Index(i)
			if !item.CanSet() {
				log.Println("can not set target field with entity")
				continue
			}
			itemAddr := item
			if item.Kind() != reflect.Ptr {
				itemAddr = item.Addr()
			}
			itemID := itemAddr.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
			if base.IBaseDbHandler.IdEquals(itemID, newEntityID) {
				found = true
				item.Set(newEntityValue)
				break
			}
		}
		if !found {
			targetField.Set(reflect.Append(targetField, newEntityValue))
		}
	} else {
		targetField.Set(newEntityValue)
	}
}

// DenormalizeNew upserts current entity in all entities which is referenced to
func (base *BaseDbHandler) DenormalizeNew(id interface{}) {
	entity, err := base.IBaseDbHandler.Get(&models.Request{
		ID: id,
	})
	if err != nil {
		log.Println(fmt.Sprintf("error on denormalize entity of %v. error: %v", entity, err))
		return
	}
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}
	for _, targetInfo := range base.DenormalizeConfigs {
		var page uint64 = 0
		for {
			page++
			result, err := targetInfo.TargetHandler.Paginate(&models.Request{
				Filters: &models.Filters{
					targetInfo.TargetIdFilter: id,
				},
				Page:    page,
				PerPage: 30,
			})
			if err != nil {
				log.Println(fmt.Sprintf("denormalize: error on entity of %v on paginate. error: %v",
					entity, err))
				break
			}
			if result.ReflectItems == nil {
				v := reflect.ValueOf(result.Items)
				result.ReflectItems = &v
			}
			for i := 0; i < result.ReflectItems.Len(); i++ {
				item := result.ReflectItems.Index(i)
				base.IBaseDbHandler.DenormalizeNewEntity(item, entityValue, targetInfo)
				if item.Kind() != reflect.Ptr {
					item = item.Addr()
				}
				entityID := item.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
				denormalized := item.Interface().(models.IBaseModel)
				err = targetInfo.TargetHandler.Update(&models.Request{
					ID:   entityID,
					Body: denormalized,
				})
				if err != nil {
					log.Println(fmt.Sprintf("denormalized: error on update entity. error: %v", err))
				}
			}
			if !result.Pagination.HasNext {
				break
			}
		}
	}
}

func (base *BaseDbHandler) DenormalizeDeleteEntity(entityValue reflect.Value, deletedEntityID interface{},
	info *DenormalizeConfig) {
	targetField := entityValue.FieldByName(info.TargetFieldName)
	if targetField.Kind() == reflect.Slice {
		for i := 0; i < targetField.Len(); i++ {
			item := targetField.Index(i)
			itemAddr := item
			if item.Kind() != reflect.Ptr {
				itemAddr = item.Addr()
			}
			itemID := itemAddr.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
			if base.IBaseDbHandler.IdEquals(itemID, deletedEntityID) {
				var items []interface{}
				for tfi := 0; tfi < targetField.Len(); tfi++ {
					if tfi == i {
						continue
					}
					items = append(items, targetField.Index(tfi))
				}
				itemsValue := reflect.ValueOf(items)
				item.Set(itemsValue)
			}
		}
	}
}

func (base *BaseDbHandler) DenormalizeDelete(id interface{}) {
	for _, targetInfo := range base.DenormalizeConfigs {
		var page uint64 = 0
		for {
			page++
			result, err := targetInfo.TargetHandler.Paginate(&models.Request{
				Filters: &models.Filters{
					targetInfo.TargetIdFilter: id,
				},
				Page:    page,
				PerPage: 30,
			})
			if err != nil {
				log.Println(fmt.Sprintf("denormalize: error on delete entity of %v on paginate. error: %v",
					id, err))
				break
			}
			if result.ReflectItems == nil {
				v := reflect.ValueOf(result.Items)
				result.ReflectItems = &v
			}
			for i := 0; i < result.ReflectItems.Len(); i++ {
				item := result.ReflectItems.Index(i)
				if item.Kind() != reflect.Ptr {
					item = item.Addr()
				}
				entityID := item.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
				base.IBaseDbHandler.DenormalizeDeleteEntity(item, entityID, targetInfo)
			}
			if !result.Pagination.HasNext {
				break
			}
		}
	}
}

// EnsureDenormalizeInterface ensures all references entities are denormalized into given entity
func (base *BaseDbHandler) EnsureDenormalizeInterface(id, entity interface{}) {
	if base.DenormalizeFieldRefConfigs == nil || len(base.DenormalizeFieldRefConfigs) == 0 {
		return
	}
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}
	for _, cfg := range base.DenormalizeFieldRefConfigs {
		field := entityValue.FieldByName(cfg.TargetFieldName)
		if !field.CanSet() {
			log.Println("ensure denormalize: can not set field")
			continue
		}
		referenceField := entityValue.FieldByName(cfg.ReferenceFieldName)
		if referenceField.IsZero() {
			typ := field.Type()
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			var el reflect.Value
			switch typ.Kind() {
			case reflect.Slice:
				el = reflect.MakeSlice(typ, 0, 0)
				break
			default:
				el = reflect.New(typ.Elem()).Elem()
				break
			}
			item := reflect.New(el.Type())
			item.Elem().Set(el)
			field.Set(item)
			continue
		}
		if referenceField.Kind() == reflect.Ptr {
			referenceField = referenceField.Elem()
		}
		fieldEl := field
		if fieldEl.Kind() == reflect.Ptr {
			fieldEl = fieldEl.Elem()
		}
		if referenceField.Kind() == reflect.Slice && (field.IsZero() || fieldEl.Kind() == reflect.Slice) {
			if !field.IsZero() {
				// remove field extra items
				for i := fieldEl.Len() - 1; i >= 0; i-- {
					item := fieldEl.Index(i)
					itemAddr := item
					if item.Kind() != reflect.Ptr {
						itemAddr = item.Addr()
					}
					itemID := itemAddr.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
					found := false
					for j := 0; j < referenceField.Len(); j++ {
						referenceID := referenceField.Index(j).Interface()
						if base.IBaseDbHandler.IdEquals(itemID, referenceID) {
							found = true
							break
						}
					}
					if !found {
						// id not found so delete denormalized entity
						fieldEl.Set(reflect.AppendSlice(fieldEl.Slice(0, i), fieldEl.Slice(i+1, fieldEl.Len())))
					}
				}
			}
			// add field missing items
			for i := referenceField.Len() - 1; i >= 0; i-- {
				referenceID := referenceField.Index(i).Interface()
				found := false
				if !field.IsZero() {
					for j := 0; j < fieldEl.Len(); j++ {
						item := fieldEl.Index(j)
						itemAddr := item
						if item.Kind() != reflect.Ptr {
							itemAddr = item.Addr()
						}
						itemID := itemAddr.MethodByName("GetID").Call([]reflect.Value{})[0].Interface()
						if base.IBaseDbHandler.IdEquals(itemID, referenceID) {
							found = true
							break
						}
					}
				}
				if !found {
					if field.IsZero() {
						typ := field.Type().Elem()
						sliceValue := reflect.MakeSlice(typ, 0, 0)
						if field.Kind() == reflect.Ptr {
							sliceValuePtr := reflect.New(sliceValue.Type())
							sliceValuePtr.Elem().Set(sliceValue)
							sliceValue = sliceValuePtr
						}
						field.Set(sliceValue)
						fieldEl = field.Elem()
					}
					// entity not found so find & add entity
					result, err := cfg.TargetHandler.Get(&models.Request{
						ID: referenceID,
					})
					if err != nil {
						continue
					}
					resultValue := reflect.ValueOf(result)
					if resultValue.Kind() == reflect.Ptr {
						resultValue = resultValue.Elem()
					}
					sliceValue := fieldEl
					if fieldEl.Kind() == reflect.Ptr {
						sliceValue = fieldEl.Elem()
					}
					sliceValue = reflect.Append(sliceValue, resultValue)
					if fieldEl.Kind() == reflect.Ptr {
						fieldEl.Elem().Set(sliceValue)
					} else {
						fieldEl.Set(sliceValue)
					}
				}
			}
		}
	}
	body := entityValue.Addr().Interface().(models.IBaseModel)
	err := base.IBaseDbHandler.Update(&models.Request{
		ID:   id,
		Body: body,
	})
	if err != nil {
		log.Println(fmt.Sprintf("denormalized: error on update entity. error: %v", err))
	}
}

func (base *BaseDbHandler) EnsureDenormalizeByID(id interface{}) {
	entity, err := base.IBaseDbHandler.Get(&models.Request{
		ID: id,
	})
	if err != nil {
		log.Println(fmt.Sprintf("error on ensure denormalize of %v. error: %v", entity, err))
		return
	}
	base.EnsureDenormalizeInterface(id, entity)
}
