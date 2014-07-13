package strsd

type Granade struct {
	Id GUID `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Way int `json:"way"`
	Timer int `json:"timer"`
	Player *Player `json:"-"`
}

func (o *Granade) Step() bool {
	o.Timer--
	if(o.Timer > 0) {
		//Move
		switch(o.Way) {
		case WAY_UP:
			if(o.Player.Game.SpaceIsInBounds(o.X,o.Y-1)) {
				o.Y--
			} else {
				o.Way = WAY_DOWN
				o.Y++
			}
		case WAY_RIGHT:
			if(o.Player.Game.SpaceIsInBounds(o.X+1,o.Y)) {
				o.X++
			} else {
				o.Way = WAY_LEFT
				o.X--
			}
		case WAY_DOWN:
			if(o.Player.Game.SpaceIsInBounds(o.X,o.Y+1)) {
				o.Y++
			} else {
				o.Way = WAY_UP
				o.Y--
			}
		case WAY_LEFT:
			if(o.Player.Game.SpaceIsInBounds(o.X-1,o.Y)) {
				o.X--
			} else {
				o.Way = WAY_RIGHT
				o.X++
			}
		}
	} else if(o.Timer == 0){
		o.Player.Game.AddEvent(o.X,o.Y,EVENT_TYPE_EXPLOSION)
		if(o.Player.Game.SpaceIsInBounds(o.X+1,o.Y)) { o.Player.Game.AddEvent(o.X+1,o.Y,EVENT_TYPE_EXPLOSION) }
		if(o.Player.Game.SpaceIsInBounds(o.X-1,o.Y)) { o.Player.Game.AddEvent(o.X-1,o.Y,EVENT_TYPE_EXPLOSION) }
		if(o.Player.Game.SpaceIsInBounds(o.X,o.Y+1)) { o.Player.Game.AddEvent(o.X,o.Y+1,EVENT_TYPE_EXPLOSION) }
		if(o.Player.Game.SpaceIsInBounds(o.X,o.Y-1)) { o.Player.Game.AddEvent(o.X,o.Y-1,EVENT_TYPE_EXPLOSION) }
		//Boom
		for _, p := range o.Player.Game.Players {
			if( (p.X == o.X && p.Y == o.Y) || (p.X == o.X-1 && p.Y == o.Y) || (p.X == o.X+1 && p.Y == o.Y) || (p.X == o.X && p.Y == o.Y-1) || (p.X == o.X && p.Y == o.Y+1) ) {
				p.Hit(o.Player,3);
			}
		}
	} else {
		return false
	}
	return true
}
