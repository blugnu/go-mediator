package mediator

import "reflect"

var commandHandlers = map[reflect.Type]interface{}{}
var queryHandlers = map[reflect.Type]interface{}{}

// reg captures a registered request type and a reference to the handler
// registry in which the registration for that type was recorded
type reg struct {
	handlers map[reflect.Type]interface{}
	rqt      reflect.Type
}

// Remove removes the handler from the handler registry
func (r *reg) Remove() {
	delete(r.handlers, r.rqt)
}
