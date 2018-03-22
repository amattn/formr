package struct2webform

import "testing"

type DummyStruct struct {
	number int64
}

func TestNonStruct(t *testing.T) {
	not_a_struct := map[string]interface{}{"hi": "what am i doing?"}

	_, err := ProcessStruct(not_a_struct)

	if err != nil {
		// good.  we should get an error
	} else {
		// uh oh.  why didn't we get an error here?
		t.Errorf("4045963788 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", not_a_struct, not_a_struct)
	}
}

func TestStruct(t *testing.T) {
	is_a_struct := DummyStruct{}
	_, err := ProcessStruct(is_a_struct)
	if err != nil {
		t.Errorf("4045963785 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", is_a_struct, is_a_struct)
	}
}
