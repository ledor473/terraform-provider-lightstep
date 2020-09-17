package lightstep

// TODO: Check if there's a reason that the API sometimes ask for *string instead of string
func stringPointer(s string) *string {
	return &s
}
