package dl

type Event struct {
	Name     string
	Callback func(fieldName, value string) interface{}
}

var events map[string]*Event

func AddNewEvent(eventName string, callback func(fieldName, value string) interface{}) {
	events[eventName] = &Event{
		Name:     eventName,
		Callback: callback,
	}
}

func TryEvent(eventName, fieldName, value string) (result interface{}, ok bool) {
	if event, ok := events[eventName]; ok {
		result = event.Callback(fieldName, value)
	}
	return
}
