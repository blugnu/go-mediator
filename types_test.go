package mediator

import "testing"

func TestThatHandlerTypeNameReturnsCorrectNames(t *testing.T) {
	tests := map[string]struct {
		handlerType handlerType
		wanted      string
	}{
		"command":      {handlerType: command, wanted: "command"},
		"query":        {handlerType: query, wanted: "query"},
		"42 (invalid)": {handlerType: 42, wanted: "<undefined>"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.handlerType.Name()
			if got != test.wanted {
				t.Errorf("wanted %q, got %q", test.wanted, got)
			}
		})
	}
}
