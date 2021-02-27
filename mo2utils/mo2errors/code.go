package mo2errors

const (
	Mo2NoExist        = 400
	Mo2NotFound       = 404
	Mo2Conflict       = 409
	Mo2LengthRequired = 411
)

var codeText = map[int]string{
	Mo2NoExist:        "object not exist",
	Mo2NotFound:       "object not found",
	Mo2Conflict:       "have conflict with exist object",
	Mo2LengthRequired: "length not match required",
}

// CodeText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func CodeText(code int) string {
	return codeText[code]
}
