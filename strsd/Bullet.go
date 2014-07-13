//A single bullet shot by player
package strsd

type Bullet struct {
	Id GUID `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Way int `json:"way"`
	Player *Player `json:"-"`
}

//Behaviour of bullet for each turn
func (o *Bullet) Step() bool {
	switch(o.Way) {
	case WAY_UP: o.Y--;
	case WAY_RIGHT: o.X++;
	case WAY_DOWN: o.Y++;
	case WAY_LEFT: o.X--;
	}
	//Check if bullet is out of the game
	if(!o.Player.Game.SpaceIsInBounds(o.X,o.Y)) {
		return false
	}
	//Check if bullet can hit any player
	for _, player := range o.Player.Game.Players {
		if(player.X == o.X && player.Y == o.Y) {
			player.Hit(o.Player,1); //Damage player for 1 point
			return false
		}
	}
	return true
}
