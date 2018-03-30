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
	Label   string
	Element string
}

// Not a stable interface yet.  I expect this function signature to change eventually.
func ProcessStruct(some_struct interface{}) ([]FormElement, error) {
	form_elements := []FormElement{}
	all_errors := []error{}

	thing_type := reflect.TypeOf(some_struct)
	thing_kind := thing_type.Kind()

	if thing_kind != reflect.Struct {
		derr := deeperror.New(417084840, "ProcessStruct() expected struct got "+fmt.Sprintf("%T", some_struct), nil)
		return form_elements, derr
	}

	for i := 0; i < thing_type.NumField(); i++ {

		struct_field := thing_type.Field(i)
		form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
		if form_field_name != "" {
			form_element, err := process_field(some_struct, i)
			if err != nil {
				all_errors = append(all_errors, err)
			} else if form_element != nil {
				form_elements = append(form_elements, *form_element)
			}

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

func process_field(some_struct interface{}, field_num int) (*FormElement, error) {
	var form_element FormElement
	var err error

	thing_type := reflect.TypeOf(some_struct)
	thing_value := reflect.ValueOf(some_struct)

	i := field_num
	struct_field := thing_type.Field(i)
	// log.Println(3073721749, struct_field)

	switch field_kind := struct_field.Type.Kind(); {
	case field_kind == reflect.String:
		current_value := thing_value.Field(i)
		unwrapped_value := current_value.String()
		form_element, err = form_output_input_type_text(struct_field, unwrapped_value)

	case field_kind == reflect.Int64:
		current_value := thing_value.Field(i)
		unwrapped_value := current_value.Int()
		string_value := strconv.FormatInt(unwrapped_value, 10)
		form_element, err = form_output_input_type_text(struct_field, string_value)

	case field_kind == reflect.Float64:
		current_value := thing_value.Field(i)
		unwrapped_value := current_value.Float()
		string_value := fmt.Sprintf("%g", unwrapped_value)
		form_element, err = form_output_input_type_text(struct_field, string_value)

	case field_kind == reflect.Bool:
		current_value := thing_value.Field(i)
		unwrapped_value := current_value.Bool()
		form_element, err = form_output_bool(struct_field, unwrapped_value)

	case field_kind == reflect.Struct:

		current_value := thing_value.Field(i)
		unknown_interface := current_value.Interface()
		unwrapped_value, is_expected_type := unknown_interface.(time.Time)
		if is_expected_type {
			marshalled_bytes, _ := unwrapped_value.MarshalText()
			string_value := string(marshalled_bytes)
			form_element, err = form_output_input_type_text(struct_field, string_value)
		}

	default:

	}

	if err != nil {
		return nil, err
	} else {
		return &form_element, err
	}

}

const (
	STRUCT_TAG_KEY_STRUCT_FIELD_TYPE = "s2w_type"
	STRUCT_TAG_KEY_FIELD_ID          = "s2w_id"
	STRUCT_TAG_KEY_FIELD_NAME        = "s2w_name"
	STRUCT_TAG_KEY_LABEL_CONTENTS    = "s2w_label"
	STRUCT_TAG_KEY_CURRENT_VALUE     = "s2w_value"

	EMPTY_VALUE_PLACEHOLDER = "PLACEHOLDER"

	DEFAULT_LABEL_TEMPLATE = `<label{{ if .s2w_id }} for="{{ .s2w_id }}"{{ end }}>{{ .s2w_label }}</label>`
)

func form_output_input_type_text(struct_field reflect.StructField, current_value string) (FormElement, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	if form_field_name == "" {
		form_field_name = EMPTY_VALUE_PLACEHOLDER
	}

	data := map[string]interface{}{
		STRUCT_TAG_KEY_STRUCT_FIELD_TYPE: struct_field.Type.String(),
		STRUCT_TAG_KEY_FIELD_ID:          form_field_id,
		STRUCT_TAG_KEY_FIELD_NAME:        form_field_name,
		STRUCT_TAG_KEY_LABEL_CONTENTS:    label_contents,
		STRUCT_TAG_KEY_CURRENT_VALUE:     current_value,
	}

	element_template := `<input type="text"{{ if .s2w_id }} id="{{ .s2w_id }}"{{ end }} class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}" value="{{ .s2w_value }}">`

	return execute_templates(1522519197, DEFAULT_LABEL_TEMPLATE, element_template, data)
}

func form_output_bool(struct_field reflect.StructField, current_value bool) (FormElement, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	is_checked := ""
	if current_value {
		is_checked = "checked"
	}

	data := map[string]interface{}{
		STRUCT_TAG_KEY_STRUCT_FIELD_TYPE: struct_field.Type.String(),
		STRUCT_TAG_KEY_FIELD_ID:          form_field_id,
		STRUCT_TAG_KEY_FIELD_NAME:        form_field_name,
		STRUCT_TAG_KEY_LABEL_CONTENTS:    label_contents,
		STRUCT_TAG_KEY_CURRENT_VALUE:     is_checked,
	}

	element_template := `<input type="checkbox"{{ if .s2w_id }} id="{{ .s2w_id }}"{{ end }} class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}"{{if .s2w_value }} {{ .s2w_value }}{{end}}>`

	return execute_templates(1963919951, DEFAULT_LABEL_TEMPLATE, element_template, data)
}

func execute_templates(debug_num int64, label_template, element_template string, data interface{}) (FormElement, error) {
	label_string, err := execute_single_template(debug_num, label_template, data)
	if err != nil {
		return FormElement{}, err
	}

	element_string, err := execute_single_template(debug_num, element_template, data)
	if err != nil {
		return FormElement{}, err
	}

	return FormElement{Label: label_string, Element: element_string}, nil
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
