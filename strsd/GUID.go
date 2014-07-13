//Unique ID generator
package strsd

//TODO: Make it channel based
type GUID int
var CurrentGUID GUID = GUID(0)
func MakeGUID() GUID {
	CurrentGUID++
	return CurrentGUID
}
