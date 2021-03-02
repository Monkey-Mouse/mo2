package mo2errors

const (
	Mo2NoError        = 200
	Mo2NoExist        = 400
	Mo2NotFound       = 404
	Mo2Conflict       = 409
	Mo2LengthRequired = 411
	Mo2Error          = 418
)

var codeText = map[int]string{
	Mo2NoExist:        "object not exist",
	Mo2NotFound:       "object not found",
	Mo2Conflict:       "have conflict with exist object",
	Mo2LengthRequired: "length not match required",
	Mo2Error:          "some internal error",
	Mo2NoError:        "no error",
}

// CodeText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func CodeText(code int) string {
	return codeText[code]
}
