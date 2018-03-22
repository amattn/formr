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

// Not a stable interface yet.  I expect this funciton signature to change eventually.
func ProcessStruct(some_struct interface{}) ([]string, error) {
	form_output := []string{}

	thing_type := reflect.TypeOf(some_struct)
	thing_kind := thing_type.Kind()

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
			form_output = append(form_output, form_output_string(struct_field))

		case field_kind == reflect.Int64:
			form_field, err := form_output_int64(struct_field)
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
	STRUCT_TAG_KEY_DEFAULT_VALUE     = "s2w_value" // TODO
)

// need field_id, field_name, label_contents
func form_output_string(struct_field reflect.StructField) string {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)
	default_value := struct_field.Tag.Get(STRUCT_TAG_KEY_DEFAULT_VALUE)

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
	return fmt.Sprintf("<label for=\"%s\">%s</label><input type=\"text\" id=\"%s\" class=\"s2w_%s\" name=\"%s\" value=\"%s\">",
		form_field_id,
		label_contents,
		form_field_id,
		struct_field.Type.String(),
		form_field_name,
		default_value)
}

func form_output_int64(struct_field reflect.StructField) (string, error) {

	form_field_id := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_ID)
	form_field_name := struct_field.Tag.Get(STRUCT_TAG_KEY_FIELD_NAME)
	label_contents := struct_field.Tag.Get(STRUCT_TAG_KEY_LABEL_CONTENTS)
	default_value := struct_field.Tag.Get(STRUCT_TAG_KEY_DEFAULT_VALUE)

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
		STRUCT_TAG_KEY_FIELD_ID:       form_field_id,
		STRUCT_TAG_KEY_FIELD_NAME:     form_field_name,
		STRUCT_TAG_KEY_LABEL_CONTENTS: label_contents,
		STRUCT_TAG_KEY_DEFAULT_VALUE:  default_value,
	}
	log.Println("3327233781 data", data)
	raw_template := `<label for="{{ .s2w_id }}">{{ .s2w_label }}</label><input type="text" id="{{ .s2w_id }}" class="s2w_{{ .s2w_type }}" name="{{ .s2w_name }}" value="{{ .s2w_value }}">`

	tmpl, err := template.New("form_output_int64").Parse(raw_template)
	if err != nil {
		derr := deeperror.New(727509799, "", err)
		derr.AddDebugField("raw_template", raw_template)
		return "", derr
	}

	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, data)
	if err != nil {
		derr := deeperror.New(3396615818, "", err)
		derr.AddDebugField("raw_template", raw_template)
		derr.AddDebugField("data", data)
		return "", derr
	}

	output := buff.String()

	return output, nil
}
