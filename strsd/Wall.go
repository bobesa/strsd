//Wall tile that blocks the bullets and bounces the granade
package strsd

//TODO: Destroyable wall
type Wall struct {
	Id GUID `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Life int `json:"life"`
}
