package struct2webform

import (
	"log"
	"reflect"
	"testing"
	"time"
)

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

type ComplexStruct struct {
	ShortString string  `s2w_id:"uname" s2w_name:"user_name" s2w_label:"Name:" s2w_value:"Bobber McTester"`
	LongString  string  `s2w_id:"ustory" s2w_name:"story" s2w_label:"Story:" `
	BigInt      int64   `s2w_id:"age_id" s2w_name:"age_name" s2w_label:"Age:" s2w_value:"42"`
	Float       float64 `s2w_id:"ratio_id" s2w_name:"ratio_name" s2w_label:"Ratio:" s2w_value:"1.55" `
	TrueOrFalse bool

	DateTime    time.Time
	DateTimePtr *time.Time
}

func TestComplexStruct(t *testing.T) {
	is_a_struct := ComplexStruct{}
	form_strings, err := ProcessStruct(is_a_struct)
	if err != nil {
		t.Errorf("1198170734 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", is_a_struct, is_a_struct)
	}

	log.Println(form_strings)

	some_type := reflect.TypeOf(is_a_struct)

	if len(form_strings) != some_type.NumField() {
		t.Errorf("Expected %d form elements of output, but got: %d", some_type.NumField(), len(form_strings))
		t.Log("Output Lines:")
		for i, form_element := range form_strings {
			t.Logf("%d %s", i, form_element)
		}
	}
}
