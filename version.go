package struct2webform

const (
	internal_BUILD_NUMBER   = 5
	internal_VERSION_STRING = "0.0.2"
)

func BuildNumber() int64 {
	return internal_BUILD_NUMBER
}
func Version() string {
	return internal_VERSION_STRING
}
