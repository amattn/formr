package struct2webform

import (
	"html/template"
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

type AllSupportedTypes struct {
	AString  string  `s2w_id:"string_id" s2w_name:"string_name" s2w_label:"AString:"`
	AInt64   int64   `s2w_id:"int64_id" s2w_name:"int64_name" s2w_label:"AInt64:"`
	AFloat64 float64 `s2w_id:"float64_id" s2w_name:"float64_name" s2w_label:"AFloat64:"`
	ABool    bool    `s2w_id:"bool_id" s2w_name:"bool_name" s2w_label:"ABool:"`

	ATime     time.Time `s2w_id:"time_id" s2w_name:"time_name" s2w_label:"ATime:"`
	ADate     time.Time `s2w_id:"date_id" s2w_name:"date_name" s2w_label:"ADate:"`
	ADateTime time.Time `s2w_id:"datetime_id" s2w_name:"datetime_name" s2w_label:"ADateTime:"`

	NoTagField string

	IgnoreThis bool `s2w_name:"-"`
}

func TestAllSupportedTypes(t *testing.T) {
	is_a_struct := AllSupportedTypes{}
	form_elements, err := ProcessStruct(is_a_struct)
	if err != nil {
		t.Errorf("2044655070 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", is_a_struct, is_a_struct)
	}

	some_type := reflect.TypeOf(is_a_struct)
	expected_num := some_type.NumField() - 1

	if len(form_elements) != expected_num {
		t.Errorf("2044655071 Expected %d form elements of output, but got: %d", expected_num, len(form_elements))
	}

	label_expecteds := make([]template.HTML, 0, expected_num)
	element_expecteds := make([]template.HTML, 0, expected_num)

	label_expecteds = append(label_expecteds, `<label for="string_id">AString:</label>`)
	label_expecteds = append(label_expecteds, `<label for="int64_id">AInt64:</label>`)
	label_expecteds = append(label_expecteds, `<label for="float64_id">AFloat64:</label>`)
	label_expecteds = append(label_expecteds, `<label for="bool_id">ABool:</label>`)
	label_expecteds = append(label_expecteds, `<label for="time_id">ATime:</label>`)
	label_expecteds = append(label_expecteds, `<label for="date_id">ADate:</label>`)
	label_expecteds = append(label_expecteds, `<label for="datetime_id">ADateTime:</label>`)
	label_expecteds = append(label_expecteds, `<label></label>`)

	element_expecteds = append(element_expecteds, `<input type="text" id="string_id" class="s2w_string" name="string_name" value="">`)
	element_expecteds = append(element_expecteds, `<input type="text" id="int64_id" class="s2w_int64" name="int64_name" value="0">`)
	element_expecteds = append(element_expecteds, `<input type="text" id="float64_id" class="s2w_float64" name="float64_name" value="0">`)
	element_expecteds = append(element_expecteds, `<input type="checkbox" id="bool_id" class="s2w_bool" name="bool_name">`)
	element_expecteds = append(element_expecteds, `<input type="text" id="time_id" class="s2w_time.Time" name="time_name" value="0001-01-01T00:00:00Z">`)
	element_expecteds = append(element_expecteds, `<input type="text" id="date_id" class="s2w_time.Time" name="date_name" value="0001-01-01T00:00:00Z">`)
	element_expecteds = append(element_expecteds, `<input type="text" id="datetime_id" class="s2w_time.Time" name="datetime_name" value="0001-01-01T00:00:00Z">`)
	element_expecteds = append(element_expecteds, `<input type="text" class="s2w_string" name="NoTagField" value="">`)

	for i, form_element := range form_elements {
		if i >= len(label_expecteds) {
			t.Fatalf("2388324744 i >= len(label_expecteds): %d >= %d", i, len(label_expecteds))
		}
		if i >= len(element_expecteds) {
			t.Fatalf("2388324745 i >= len(element_expecteds): %d >= %d", i, len(element_expecteds))
		}

		if form_element.Label != label_expecteds[i] {
			t.Errorf("2044655072 %d %s\n%s\n%s\n%s", i, "Expected:", label_expecteds[i], "     Got:", form_element.Label)
		}
		if form_element.Element != element_expecteds[i] {
			t.Errorf("2044655073 %d %s\n%s\n%s\n%s", i, "Expected:", element_expecteds[i], "     Got:", form_element.Element)
		}
	}

}

type ComplexStruct struct {
	ShortString string  `s2w_id:"uname" s2w_name:"user_name" s2w_label:"Name:"`
	LongString  string  `s2w_id:"ustory" s2w_name:"story"`
	BigInt      int64   `s2w_id:"age_id" s2w_name:"age_name" s2w_label:"Age:"`
	Float       float64 `s2w_id:"ratio_id" s2w_name:"ratio_name" s2w_label:"Ratio:"`
	TrueOrFalse bool    `s2w_id:"is_checked_id" s2w_name:"is_checked_name" s2w_label:"Checked:"`

	IgnoreThisOne bool `s2w_name:"-"`
	IgnoreThisToo bool `s2w_name:"-" s2w_id:"should_be_ignored" s2w_label:"What?:"`

	DateTime    time.Time  `s2w_name:"time_name" s2w_label:"Date:"`
	DateTimePtr *time.Time `s2w_id:"time_ptr_id" s2w_name:"time_ptr_name" s2w_label:"DatePtr:"`
}

// s2w only counts struct fields that have the s2w_name tag
func TestNumFormFields(t *testing.T) {
	is_a_struct := ComplexStruct{}
	form_strings, _ := ProcessStruct(is_a_struct)

	some_type := reflect.TypeOf(is_a_struct)
	expected_num := some_type.NumField() - 2

	if len(form_strings) != expected_num {
		t.Errorf("1655060571 Expected %d form elements of output, but got: %d", expected_num, len(form_strings))
	}
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

	// log.Println(1142352648, form_strings)

	some_type := reflect.TypeOf(is_a_struct)

	expected_num := some_type.NumField() - 2

	if len(form_strings) != expected_num {
		t.Errorf("Expected %d form elements of output, but got: %d", expected_num, len(form_strings))
		t.Log("Output Lines:")
		for i, form_element := range form_strings {
			t.Logf("%d %s", i, form_element)
		}
	}
}

type CheckboxStruct struct {
	StartTrue  bool `s2w_name:"StartTrue"`
	StartFalse bool `s2w_name:"StartFalse"`
	FilledOut  bool `s2w_id:"cool_id" s2w_name:"FilledOut" s2w_label:"Are we cool?" `
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
	expected_num_fields := some_type.NumField()
	if len(form_elements) != expected_num_fields {
		t.Errorf("Expected %d form elements of output, but got: %d", expected_num_fields, len(form_elements))
		t.Log("Output Lines:")
		for i, form_element := range form_elements {
			t.Logf("%d %s", i, form_element)
		}
	}

	label_expecteds := make([]template.HTML, expected_num_fields)
	element_expecteds := make([]template.HTML, expected_num_fields)

	label_expecteds[0] = `<label></label>`
	label_expecteds[1] = `<label></label>`
	label_expecteds[2] = `<label for="cool_id">Are we cool?</label>`

	element_expecteds[0] = `<input type="checkbox" class="s2w_bool" name="StartTrue" checked>`
	element_expecteds[1] = `<input type="checkbox" class="s2w_bool" name="StartFalse">`
	element_expecteds[2] = `<input type="checkbox" id="cool_id" class="s2w_bool" name="FilledOut" checked>`

	for i, form_element := range form_elements {
		if i >= len(label_expecteds) {
			t.Fatalf("i > len(label_expecteds): %d >= %d", i, len(label_expecteds))
		}
		if i >= len(element_expecteds) {
			t.Fatalf("i > len(element_expecteds): %d >= %d", i, len(element_expecteds))
		}

		if form_element.Label != label_expecteds[i] {
			t.Errorf("328727756 %d %s\n%s\n%s\n%s", i, "Expected:", label_expecteds[i], "     Got:", form_element.Label)
		}
		if form_element.Element != element_expecteds[i] {
			t.Errorf("328727757 %d %s\n%s\n%s\n%s", i, "Expected:", element_expecteds[i], "     Got:", form_element.Element)
		}
	}

}

type GorillaCompatibility struct {
	ShortString   string `schema:"short_string"`
	OtherString   string
	IgnoreThis    string `schema:"-"`
	IgnoreThisToo string `schema:"-"`
}

func TestGorillaCompatibility(t *testing.T) {

	is_a_struct := GorillaCompatibility{}
	is_a_struct.ShortString = "ShortString"
	is_a_struct.OtherString = "OtherString"
	is_a_struct.IgnoreThis = "IgnoreThis"

	form_elements, err := ProcessStruct(is_a_struct)
	if err != nil {
		t.Errorf("3366803382 %s type:%T desc:%v", "Expected err != nil but it seems to be nil", is_a_struct, is_a_struct)
	}

	some_type := reflect.TypeOf(is_a_struct)
	expected_num_fields := some_type.NumField() - 2

	if len(form_elements) != expected_num_fields {
		t.Errorf("Expected %d form elements of output, but got: %d", expected_num_fields, len(form_elements))
		t.Log("Output Lines:")
		for i, form_element := range form_elements {
			t.Logf("%d %s", i, form_element)
		}
	}

	label_expecteds := make([]template.HTML, 0, expected_num_fields)
	element_expecteds := make([]template.HTML, 0, expected_num_fields)

	label_expecteds = append(label_expecteds, `<label></label>`)
	label_expecteds = append(label_expecteds, `<label></label>`)

	element_expecteds = append(element_expecteds, `<input type="text" class="s2w_string" name="short_string" value="ShortString">`)
	element_expecteds = append(element_expecteds, `<input type="text" class="s2w_string" name="OtherString" value="OtherString">`)

	for i, form_element := range form_elements {
		if i >= len(label_expecteds) {
			t.Fatalf("i > len(label_expecteds): %d >= %d", i, len(label_expecteds))
		}
		if i >= len(element_expecteds) {
			t.Fatalf("i > len(element_expecteds): %d >= %d", i, len(element_expecteds))
		}

		if form_element.Label != label_expecteds[i] {
			t.Errorf("328727756 %d %s\n%s\n%s\n%s", i, "Expected:", label_expecteds[i], "     Got:", form_element.Label)
		}
		if form_element.Element != element_expecteds[i] {
			t.Errorf("328727757 %d %s\n%s\n%s\n%s", i, "Expected:", element_expecteds[i], "     Got:", form_element.Element)
		}
	}
}
