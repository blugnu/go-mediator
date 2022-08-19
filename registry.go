package mediator

import "reflect"

var receivers = map[reflect.Type]interface{}{}
var handlers = map[reflect.Type]interface{}{}

// reg captures a registered type and a reference to the
// map in which the registration for that type was recorded
type reg struct {
	registry       map[reflect.Type]interface{}
	registeredtype reflect.Type
}

// Remove removes the registration entry for the recorded type
// from the registry where it was registered
func (r *reg) Remove() {
	delete(r.registry, r.registeredtype)
}
