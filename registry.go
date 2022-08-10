package mediator

import "reflect"

var commandHandlers = map[reflect.Type]interface{}{}
var queryHandlers = map[reflect.Type]interface{}{}

// reg is returned from RegisterHandler calls.  It is typically ignored
// when registering production handlers but may be used in test handlers
// to remove a handler registration once the test is complete:
//
//   reg := mediator.RegisterCommandHandler[*some.Request](&mock.SomeFake{})
//   defer reg.Remove()
type reg struct {
	handlers map[reflect.Type]interface{}
	rqt      reflect.Type
}

// Remove removes the handler from the handler registry
func (r *reg) Remove() {
	delete(r.handlers, r.rqt)
}
