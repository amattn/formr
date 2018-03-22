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
	ShortString string  `s2w_id:"uname" s2w_name:"user_name" s2w_label:"Name:"`
	LongString  string  `s2w_id:"ustory" s2w_name:"story" s2w_label:"Story:" `
	BigInt      int64   `s2w_id:"age_id" s2w_name:"age_name" s2w_label:"Age:"`
	Float       float64 `s2w_id:"ratio_id" s2w_name:"ratio_name" s2w_label:"Ratio:"`
	TrueOrFalse bool

	DateTime    time.Time
	DateTimePtr *time.Time
}

func TestComplexStruct(t *testing.T) {
	is_a_struct := ComplexStruct{}
	is_a_struct.ShortString = "Bobber McTester"
	is_a_struct.LongString = "This is not my story"
	is_a_struct.BigInt = 42
	is_a_struct.Float = 1.55
	is_a_struct.TrueOrFalse = true
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

type CheckboxStruct struct {
	StartTrue  bool
	StartFalse bool
	FilledOut  bool `s2w_id:"cool_id" s2w_name:"cool_name" s2w_label:"Are we cool?" `
}

func TestCheckbox(t *testing.T) {

	is_a_struct := CheckboxStruct{}
	is_a_struct.StartTrue = true
	is_a_struct.FilledOut = true

	form_elements, err := ProcessStruct(is_a_struct)
	if err != nil {
		t.Errorf("3366803382 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", is_a_struct, is_a_struct)
	}

	some_type := reflect.TypeOf(is_a_struct)
	if len(form_elements) != some_type.NumField() {
		t.Errorf("Expected %d form elements of output, but got: %d", some_type.NumField(), len(form_elements))
		t.Log("Output Lines:")
		for i, form_element := range form_elements {
			t.Logf("%d %s", i, form_element)
		}
	}

	label_expecteds := make([]string, 3)
	element_expecteds := make([]string, 3)

	label_expecteds[0] = `<label for="PLACEHOLDER">PLACEHOLDER</label>`
	label_expecteds[1] = `<label for="PLACEHOLDER">PLACEHOLDER</label>`
	label_expecteds[2] = `<label for="cool_id">Are we cool?</label>`

	element_expecteds[0] = `<input type="checkbox" id="PLACEHOLDER" class="s2w_bool" name="PLACEHOLDER" checked>`
	element_expecteds[1] = `<input type="checkbox" id="PLACEHOLDER" class="s2w_bool" name="PLACEHOLDER" >`
	element_expecteds[2] = `<input type="checkbox" id="cool_id" class="s2w_bool" name="cool_name" checked>`

	for i, form_element := range form_elements {
		if form_element.Label != label_expecteds[i] {
			t.Errorf("328727756 %d %s\n%s\n%s\n%s", i, "Expected:", label_expecteds[i], "     Got:", form_element.Label)
		}
		if form_element.Element != element_expecteds[i] {
			t.Errorf("328727757 %d %s\n%s\n%s\n%s", i, "Expected:", element_expecteds[i], "     Got:", form_element.Element)
		}
	}

}
