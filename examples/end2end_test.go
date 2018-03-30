package main

import (
	"bytes"
	"html/template"
	"testing"
	"time"

	"../../formr"
)

func TestTemplate(t *testing.T) {
	some_struct := ExampleStruct{}
	some_struct.AInt64 = 99
	some_struct.ATime = time.Unix(1600000000, 0)
	elements, err := formr.ProcessStruct(some_struct)
	template_data := ExampleTemplateData{}
	template_data.FormElements = elements

	tmpl, err := template.New("index_template").Parse(index_template)
	if err != nil {
		t.Fatal(2987939295, err)
	}

	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, template_data)
	if err != nil {
		t.Fatal(2987939296, err)
	}
}
