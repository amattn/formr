# formr

Convert your go (golang) structs into scaffolding-style webform.


If your struct looks like this:

	type SomeStruct struct {
		ShortString string  `formr:"user_name" formr_id:"uname" formr_label:"Name:" formr_value:"Bobber McTester"`
		BigInt      int64   `formr:"age_name" formr_id:"age_id" formr_label:"Age:" formr_value:"42"`
		Untagged    string
	}

The output of `ProcessStruct(SomeStruct{})` will be a slice of strings like this:

	[
		<label for="uname">Name:</label><input type="text" id="uname" class="formr_string" name="user_name" value="Bobber McTester">
		<label for="age_id">Age:</label><input type="text" id="age_id" class="formr_" name="age_name" value="42">
		<label></label><input type="text" class="formr_string" name="Untagged" value="">
	]

## Install

	go get github.com/amattn/formr

or if you want to update to latest:

	go get -u github.com/amattn/formr


`formr` is currently considered beta software.  We recommend vendoring for now.

## Usage

Currently supported types:

- string
- int64
- float64
- bool
- time.Time
- pointers to any of the above types

Basically, make a struct with the supported types.  You can add struct tags to customize the id, label.  Default values are taken from the struct itself.
Pass that struct to `ProcessStruct()` and that will generate a `[]FormElement` slice that you can put into a `html/template` to populate a form.

Using a package such as https://github.com/gorilla/schema you can parse the generated forms.  See end2end.go in the examples folder for a working implementation.

## License

MIT

see LICENSE file for details