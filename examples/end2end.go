package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	formr "../../struct2webform"
	gschema "github.com/gorilla/schema"
)

type ExampleStruct struct {
	AString string `s2w_id:"string_id" s2w_label:"AString:" schema:"string_name"`
	AInt64  int64  `s2w_id:"int64_id" s2w_label:"AInt64:" schema:"int64_name"`
	ABool   bool   `s2w_id:"bool_id" s2w_label:"ABool:" schema:"bool_name"`

	ATime time.Time `s2w_id:"time_id" s2w_label:"ATime:" schema:"time_name"`
}

const index_template = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>My test page</title>
  </head>
  <body>
    <form action="/post_handler" method="post">
        {{ range .FormElements}}
            {{ .Label }}
            {{ .Element }}
            <br/>
        {{ end }}

        <input type="submit" value="Send Request">

    </form>
  </body>
</html>
`

type ExampleTemplateData struct {
	FormElements []formr.FormElement
}

func main() {

	some_struct := ExampleStruct{}
	some_struct.AInt64 = 99
	some_struct.ATime = time.Unix(1500000000, 0)

	elements, err := formr.ProcessStruct(some_struct)
	template_data := ExampleTemplateData{}
	template_data.FormElements = elements

	tmpl, err := template.New("index_template").Parse(index_template)
	if err != nil {
		log.Fatalln(3156803913, err)
	}

	http.HandleFunc("/post_handler", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		// we use gorilla/schema to parse the form.
		err := r.ParseForm()
		if err != nil {
			log.Println(3318985455, err)
			http.Error(w, "3318985455 ParseForm Failure", http.StatusInternalServerError)
			return
		}

		form_response := new(ExampleStruct)
		decoder := gschema.NewDecoder()
		err = decoder.Decode(form_response, r.PostForm)
		if err != nil {
			log.Println(3318985457, err)
			http.Error(w, "3318985457 Decode Failure", http.StatusInternalServerError)
			return
		}

		// It's not good practice to use the post-handler to present content.  this will break if a user tries to bookmark this page for example.
		// normally you redirect to a different page like so:
		// http.Redirect(w, r, "/your_answers", http.StatusSeeOther)
		fmt.Fprintf(w, "<html>You submitted:\n<br/>AString:%v\n<br/>AInt64:%v\n<br/>ABool:%v\n<br/>ATime:%v\n</html>",
			form_response.AString,
			form_response.AInt64,
			form_response.ABool,
			form_response.ATime,
		)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		buff := new(bytes.Buffer)
		err := tmpl.Execute(buff, template_data)
		if err != nil {
			log.Println(83632357, "HandleFunc ERROR:", err)
			http.Error(w, "503 error: 83632357", http.StatusInternalServerError)
			return
		}

		_, err = buff.WriteTo(w)
		if err != nil {
			log.Println(1468074579, "Write error 2232810591", err)
		}

	})

	host_and_port := ":8273"
	log.Println("Starting HTTP server at", "\nhttp://"+host_and_port)
	log.Fatal(http.ListenAndServe(host_and_port, nil))

}
