# struct2webform

Convert your go (golang) structs into scaffolding-style webform.


If your struct looks like this:

	type SomeStruct struct {
		ShortString string  `s2w_id:"uname" s2w_name:"user_name" s2w_label:"Name:" s2w_value:"Bobber McTester"`
		BigInt      int64   `s2w_id:"age_id" s2w_name:"age_name" s2w_label:"Age:" s2w_value:"42"`
	}

The output of `ProcessStruct(SomeStruct{})` will be a slice of strings like this:

	[
		<label for="uname">Name:</label><input type="text" id="uname" class="s2w_string" name="user_name" value="Bobber McTester">
		<label for="age_id">Age:</label><input type="text" id="age_id" class="s2w_" name="age_name" value="42">
	]

## Install

TBD

## Usage

Currently supported types:

- string
- int64
- float64
- bool

TBD

## License

MIT

see LICENSE file for details