package waf

const (
	errorStringsExpand    string = "%v was not successfully converted to []string, type: %T"
	errorIntExpand        string = "%v was not an int, actual: %T"
	errorBoolExpand       string = "%v was not a bool, actual: %T"
	errorStringExpand     string = "%v was not a string, actual: %T"
	errorInterfacesExpand string = "%v was not a []interface{}, actual: %T"
)
