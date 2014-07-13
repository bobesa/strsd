//A simple http request type
package strsd

type Request struct {
	Player *Player
	Action int
	Callback chan []byte
}
