package infra

import "testing"

func TestShouldCloneObjects(t *testing.T) {

	type Test struct {
		ID string
	}
	from := new(Test)
	from.ID = "A"

	to := new(Test)

	err := Clone(from, to)

	if err != nil {
		t.Fail()
	}
	if to.ID != from.ID {
		t.Fail()
	}
}
