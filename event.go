package dl

type Event struct {
	Name     string
	Callback func(fieldName string, value interface{}) interface{}
}

var events = map[string]*Event{}

func AddNewEvent(eventName string, callback func(fieldName string, value interface{}) interface{}) {
	events[eventName] = &Event{
		Name:     eventName,
		Callback: callback,
	}
}

func TryEvent(eventName, fieldName string, value interface{}) (result interface{}, ok bool) {
	if event, exists := events[eventName]; exists {
		ok = true
		result = event.Callback(fieldName, value)
	}
	return
}
