package strsd

type Bullet struct {
	Id GUID `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Way int `json:"way"`
	Player *Player `json:"-"`
}

func (o *Bullet) Step() bool {
	switch(o.Way) {
	case WAY_UP: o.Y--;
	case WAY_RIGHT: o.X++;
	case WAY_DOWN: o.Y++;
	case WAY_LEFT: o.X--;
	}
	if(!o.Player.Game.SpaceIsInBounds(o.X,o.Y)) {
		return false
	}
	for _, player := range o.Player.Game.Players {
		if(player.X == o.X && player.Y == o.Y) {
			player.Hit(o.Player,1);
			return false;
		}
	}
	return true
}
