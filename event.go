package dl

import "reflect"

type Event struct {
	Name     string
	Callback func(fieldName string, value interface{}) (interface{}, error)
}

var events = map[string]*Event{}

func AddNewEvent(eventName string, callback func(fieldName string, value interface{}) (interface{}, error)) {
	events[eventName] = &Event{
		Name:     eventName,
		Callback: callback,
	}
}

func TryEvent(eventName, fieldName string, value interface{}) (result interface{}, ok bool) {
	if event, exists := events[eventName]; exists {
		v := reflect.ValueOf(value)
		result = nil
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				return
			}
			value = v.Elem().Interface()
			break
		case reflect.Slice:
			if v.Len() == 0 {
				return
			}
		}
		ok = true
		eventResult, err := event.Callback(fieldName, value)
		if err != nil {
			return nil, false
		}
		result = eventResult
	}
	return
}
