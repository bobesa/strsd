//Events in turn (Can be a shot animation etc.)
package strsd

type Event struct {
	Type int `json:"type"`
	X int `json:"x"`
	Y int `json:"y"`
}
