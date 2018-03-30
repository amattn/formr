package struct2webform

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"

	"github.com/amattn/deeperror"
)

func init() {

}

type FormElement struct {
	Label   template.HTML
	Element template.HTML
}

// Not a stable interface yet.  I expect this function signature to change eventually.
func ProcessStruct(some_struct interface{}) ([]FormElement, error) {
	form_elements := []FormElement{}
	all_errors := []error{}

	thing_type := reflect.TypeOf(some_struct)
	thing_value := reflect.ValueOf(some_struct)
	thing_kind := thing_type.Kind()

	if thing_kind != reflect.Struct {
		derr := deeperror.New(417084840, "ProcessStruct() expected struct got "+fmt.Sprintf("%T", some_struct), nil)
		return form_elements, derr
	}

	for i := 0; i < thing_type.NumField(); i++ {

		struct_field := thing_type.Field(i)
		form_element, err := process_field(struct_field, thing_value.Field(i))
		if err != nil {
			all_errors = append(all_errors, err)
		} else if form_element != nil {
			form_elements = append(form_elements, *form_element)
		}
	}

	if len(all_errors) > 0 {
		err_string := fmt.Sprintf("%d errors encountered", len(all_errors))
		for _, inner_err := range all_errors {
			err_string += inner_err.Error()
		}

		err := fmt.Errorf("3102914638 %s", err_string)
		return form_elements, err
	}

	return form_elements, nil
}

func process_field(struct_field reflect.StructField, struct_field_value reflect.Value) (*FormElement, error) {
	var form_element *FormElement
	var err error

	// log.Println(2531510246, struct_field)
	// log.Println(2531510247, struct_field_value)

	switch field_kind := struct_field_value.Kind(); {
	case field_kind == reflect.String:
		unwrapped_value := struct_field_value.String()
		form_element, err = form_output_input_type_text(struct_field, unwrapped_value)

	case field_kind == reflect.Int64:
		unwrapped_value := struct_field_value.Int()
		string_value := strconv.FormatInt(unwrapped_value, 10)
		form_element, err = form_output_input_type_text(struct_field, string_value)

	case field_kind == reflect.Float64:
		unwrapped_value := struct_field_value.Float()
		string_value := fmt.Sprintf("%g", unwrapped_value)
		form_element, err = form_output_input_type_text(struct_field, string_value)

	case field_kind == reflect.Bool:
		unwrapped_value := struct_field_value.Bool()
		form_element, err = form_output_bool(struct_field, unwrapped_value)

	case field_kind == reflect.Struct:
		unknown_interface := struct_field_value.Interface()
		unwrapped_value, is_expected_type := unknown_interface.(time.Time)
		if is_expected_type {
			marshalled_bytes, _ := unwrapped_value.MarshalText()
			string_value := string(marshalled_bytes)
			form_element, err = form_output_input_type_text(struct_field, string_value)
		}

	case field_kind == reflect.Ptr:
		if struct_field_value.IsNil() {

			dereferenced_type := struct_field_value.Type().Elem()
			return process_field(struct_field, reflect.Zero(dereferenced_type))
		} else {
			return process_field(struct_field, struct_field_value.Elem())
		}
	default:

	}

	if err != nil {
		return nil, err
	} else {
		return form_element, err
	}

}

const (
	STRUCT_TAG_KEY_STRUCT_FIELD_TYPE = "s2w_type"
	STRUCT_TAG_KEY_FIELD_ID          = "s2w_id"
	STRUCT_TAG_KEY_FIELD_NAME        = "s2w_name"
	STRUCT_TAG_KEY_LABEL_CONTENTS    = "s2w_label"
	STRUCT_TAG_KEY_CURRENT_VALUE     = "s2w_value"

	STRUCT_TAG_KEY_GORILLA_SCHEMA_KEY_NAME = "schema"

	EMPTY_VALUE_MARK = "-"

	DEFAULT_LABEL_TEMPLATE = `<label{{ if .s2w_id }} for="{{ .s2w_id }}"{{ end }}>{{ .s2w_label }}</label>`
)

func process_tag(struct_field reflect.StructField) (map[string]interface{}, bool) {
	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	// first check for s2w_name
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	if form_field_name == "" {

		// then check for "schema"
		gorilla_schema_tag_name := struct_field.Tag.Get(STRUCT_TAG_KEY_GORILLA_SCHEMA_KEY_NAME)

		if gorilla_schema_tag_name != "" {
			form_field_name = gorilla_schema_tag_name
		} else {
			// if no tag, use field name
			form_field_name = struct_field.Name
		}
	}

	if form_field_name == "-" {
		// Bail out
		return nil, false
	}

	// log.Println(3013947950, form_field_name)

	data := map[string]interface{}{
		STRUCT_TAG_KEY_STRUCT_FIELD_TYPE: struct_field.Type.String(),
		STRUCT_TAG_KEY_FIELD_ID:          form_field_id,
		STRUCT_TAG_KEY_FIELD_NAME:        form_field_name,
		STRUCT_TAG_KEY_LABEL_CONTENTS:    label_contents,
	}
	return data, true
}

func form_output_input_type_text(struct_field reflect.StructField, current_value string) (*FormElement, error) {

	data, is_ok := process_tag(struct_field)
	if is_ok == false {
		return nil, nil
	}

	data[STRUCT_TAG_KEY_CURRENT_VALUE] = current_value

	element_template := `<input type="text"{{ if .s2w_id }} id="{{ .s2w_id }}"{{ end }} class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}" value="{{ .s2w_value }}">`

	return execute_templates(1522519197, DEFAULT_LABEL_TEMPLATE, element_template, data)
}

func form_output_bool(struct_field reflect.StructField, current_value bool) (*FormElement, error) {

	data, is_ok := process_tag(struct_field)
	if is_ok == false {
		return nil, nil
	}

	is_checked := ""
	if current_value {
		is_checked = "checked"
	}

	data[STRUCT_TAG_KEY_CURRENT_VALUE] = is_checked

	element_template := `<input type="checkbox"{{ if .s2w_id }} id="{{ .s2w_id }}"{{ end }} class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}"{{if .s2w_value }} {{ .s2w_value }}{{end}}>`

	element, err := execute_templates(1963919951, DEFAULT_LABEL_TEMPLATE, element_template, data)
	return element, err
}

func execute_templates(debug_num int64, label_template, element_template string, data interface{}) (*FormElement, error) {
	label_string, err := execute_single_template(debug_num, label_template, data)
	if err != nil {
		return nil, err
	}

	element_string, err := execute_single_template(debug_num, element_template, data)
	if err != nil {
		return nil, err
	}

	element := FormElement{Label: template.HTML(label_string), Element: template.HTML(element_string)}

	return &element, nil
}

func execute_single_template(debug_num int64, raw_template string, data interface{}) (string, error) {
	tmpl, err := template.New(fmt.Sprintf("%d", debug_num)).Parse(raw_template)
	if err != nil {
		derr := deeperror.New(debug_num, "Parse failure", err)
		derr.AddDebugField("raw_template", raw_template)
		return "", derr
	}

	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, data)
	if err != nil {
		derr := deeperror.New(debug_num, "Execute failure", err)
		derr.AddDebugField("raw_template", raw_template)
		derr.AddDebugField("data", data)
		return "", derr
	}

	output := buff.String()

	return output, nil
}
