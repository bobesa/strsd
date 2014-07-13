package strsd

type Wall struct {
	Id GUID `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Life int `json:"life"`
}
