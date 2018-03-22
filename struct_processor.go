package struct2webform

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"reflect"

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
	form_output := []FormElement{}

	thing_type := reflect.TypeOf(some_struct)
	thing_kind := thing_type.Kind()
	thing_value := reflect.ValueOf(some_struct)

	if thing_kind != reflect.Struct {
		derr := deeperror.New(417084840, "ProcessStruct() expected struct got "+fmt.Sprintf("%T", some_struct), nil)
		return form_output, derr
	}

	all_errors := []error{}

	for i := 0; i < thing_type.NumField(); i++ {

		struct_field := thing_type.Field(i)

		log.Println(struct_field)

		switch field_kind := struct_field.Type.Kind(); {
		case field_kind == reflect.String:
			current_value := thing_value.Field(i)
			unwrapped_value := current_value.String()
			form_field, err := form_output_string(struct_field, unwrapped_value)
			if err != nil {
				all_errors = append(all_errors, err)
			} else {
				form_output = append(form_output, form_field)
			}

		case field_kind == reflect.Int64:
			current_value := thing_value.Field(i)
			unwrapped_value := current_value.Int()
			form_field, err := form_output_int64(struct_field, unwrapped_value)
			if err != nil {
				all_errors = append(all_errors, err)
			} else {
				form_output = append(form_output, form_field)
			}

		case field_kind == reflect.Bool:
			current_value := thing_value.Field(i)
			unwrapped_value := current_value.Bool()
			form_field, err := form_output_bool(struct_field, unwrapped_value)
			if err != nil {
				all_errors = append(all_errors, err)
			} else {
				form_output = append(form_output, form_field)
			}
		default:

		}
	}

	return form_output, nil
}

const (
	STRUCT_TAG_KEY_STRUCT_FIELD_TYPE = "s2w_type"
	STRUCT_TAG_KEY_FIELD_ID          = "s2w_id"
	STRUCT_TAG_KEY_FIELD_NAME        = "s2w_name"
	STRUCT_TAG_KEY_LABEL_CONTENTS    = "s2w_label"
	// STRUCT_TAG_KEY_DEFAULT_VALUE     = "s2w_value" // TODO
	STRUCT_TAG_KEY_CURRENT_VALUE = "s2w_value"
)

// need field_id, field_name, label_contents
func form_output_string(struct_field reflect.StructField, current_value string) (FormElement, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	if form_field_id == "" {
		form_field_id = "PLACEHOLDER"
	}
	if form_field_name == "" {
		form_field_name = form_field_id
	}
	if label_contents == "" {
		label_contents = form_field_name
	}

	// TODO consider switch to a template??

	label_string := fmt.Sprintf("<label for=\"%s\">%s</label>",
		form_field_id,
		label_contents)

	element_string := fmt.Sprintf("<input type=\"text\" id=\"%s\" class=\"s2w_%s\" name=\"%s\" value=\"%s\">",
		form_field_id,
		struct_field.Type.String(),
		form_field_name,
		current_value)

	return FormElement{Label: label_string, Element: element_string}, nil
}

func form_output_int64(struct_field reflect.StructField, current_value int64) (FormElement, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	if form_field_id == "" {
		form_field_id = "PLACEHOLDER"
	}
	if form_field_name == "" {
		form_field_name = form_field_id
	}
	if label_contents == "" {
		label_contents = form_field_name
	}

	data := map[string]interface{}{
		STRUCT_TAG_KEY_STRUCT_FIELD_TYPE: struct_field.Type.String(),
		STRUCT_TAG_KEY_FIELD_ID:          form_field_id,
		STRUCT_TAG_KEY_FIELD_NAME:        form_field_name,
		STRUCT_TAG_KEY_LABEL_CONTENTS:    label_contents,
		STRUCT_TAG_KEY_CURRENT_VALUE:     current_value,
	}

	label_template := `<label for="{{ .s2w_id }}">{{ .s2w_label }}</label>`
	element_template := `<input type="text" id="{{ .s2w_id }}" class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}" value="{{ .s2w_value }}">`

	return execute_templates(1522519197, label_template, element_template, data)
}

func form_output_bool(struct_field reflect.StructField, current_value bool) (FormElement, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)

	if form_field_id == "" {
		form_field_id = "PLACEHOLDER"
	}
	if form_field_name == "" {
		form_field_name = form_field_id
	}
	if label_contents == "" {
		label_contents = form_field_name
	}

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

	label_template := `<label for="{{ .s2w_id }}">{{ .s2w_label }}</label>`
	element_template := `<input type="checkbox" id="{{ .s2w_id }}" class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}" {{if .s2w_value }}{{ .s2w_value }}{{end}}>`

	return execute_templates(1963919951, label_template, element_template, data)
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
