package strsd

type GUID int
var CurrentGUID GUID = GUID(0)
func MakeGUID() GUID {
	CurrentGUID++
	return CurrentGUID
}
