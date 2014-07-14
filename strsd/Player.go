//Player type and FoV logic
package strsd

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/akavel/polyclip-go"
)

//All player instances
var AllPlayers Players = make(Players,0)

type Players map[string]*Player
type Player struct {
	Name string `json:"name"`
	X int `json:"x"`
	Y int `json:"y"`
	Way int `json:"way"`
	Life int `json:"life"`
	Granades int `json:"granades"`
	Bullets int `json:"bullets"`
	FoV [][]bool `json:"-"`

	Score int `json:"-"`
	TimeSinceLastAction int `json:"-"`
	Action int `json:"-"`
	Game *Game `json:"-"`
}

//Player instance creation
func MakePlayer(name string, game *Game) *Player {
	fmt.Println("Player '"+name+"' joined the game!")
	p := Player{Name:name,Action:ACTION_NONE,Game:game,Score:0}
	p.Spawn()

	AllPlayers[name] = &p
	return &p
}

//Generate simple box-shaped polygon from position 
func GetBox(x, y int) polyclip.Polygon {
	return polyclip.Polygon{
		{
			{float64(x), float64(y)},
			{float64(x), float64(y+1)},
			{float64(x+1), float64(y+1)},
			{float64(x+1), float64(y)},
		},
	}
}

//Compute angle between 2 positions
func GetAngle(p1 polyclip.Point, p2 polyclip.Point) float64 {
	return math.Atan2(p1.Y - p2.Y, p1.X - p2.X)
}

//Damaging player instance
func (o *Player) Hit(attacker *Player, force int) {
	o.Life -= force
	if(o.Life <= 0) {
		//Update score (no suicide)
		if(o != attacker) {
			attacker.Score++
		}
		//Respawn
		o.Spawn();
	}
}

//Compute Field of View
func (p *Player) ComputeFoV() {
	//Mark all tiles as visible
	p.FoV = make([][]bool,p.Game.Width);
	for x:=0; x<p.Game.Width; x++ {
		p.FoV[x] = make([]bool,p.Game.Height)
		for y:=0; y<p.Game.Height; y++ {
			p.FoV[x][y] = true
		}
	}

	//Disable FieldOfView?
	if(!p.Game.OptionFieldOfView){
		return
	}

	//Visible distance
	size := float64(p.Game.Width)
	if(size < float64(p.Game.Height)) { size = float64(p.Game.Height); }

	//Player's FoV
	pp := polyclip.Point{float64(p.X), float64(p.Y)}
	t1, t2 := pp, pp

	//Depending of Player facing we need to compute view points
	//FYI: Tiny float modifications are for computation corrections
	switch(p.Way) {
	case WAY_LEFT:
		pp.X += 1
		pp.Y += 0.1
		t1, t2 = pp, pp
		t1.X -= size
		t1.Y -= size
		t2.X -= size
		t2.Y += size
	case WAY_RIGHT:
		pp.Y += 0.1
		t1, t2 = pp, pp
		t1.X += size
		t1.Y -= size
		t2.X += size
		t2.Y += size
	case WAY_UP:
		pp.X += 0.1
		pp.Y += 1
		t1, t2 = pp, pp
		t1.X -= size
		t1.Y -= size
		t2.X += size
		t2.Y -= size
	case WAY_DOWN:
		pp.X += 0.1
		t1, t2 = pp, pp
		t1.X -= size
		t1.Y += size
		t2.X += size
		t2.Y += size
	}

	//Hide everything outside of Field of View
	fov := polyclip.Polygon{{pp, t1, t2}}
	for x := 0; x < p.Game.Width; x++ {
		for y := 0; y < p.Game.Height; y++ {
			if(p.FoV[x][y] && len(GetBox(x,y).Construct(polyclip.INTERSECTION, fov)) == 0) {
				p.FoV[x][y] = false;
			}
		}
	}

	//Shadowing obstacles
	pp = polyclip.Point{float64(p.X)+.5, float64(p.Y)+.5}
	for _, wall := range p.Game.Walls {
		wp1, wp2 := polyclip.Point{float64(wall.X), float64(wall.Y)}, polyclip.Point{float64(wall.X), float64(wall.Y)}

		//Select correct points from box for computing projection points
		//Hint: Look at numkeys - 5 = player, number = box
		if(wall.X == p.X && wall.Y > p.Y) { //#2
			wp1.Y++
			wp2.X++
			wp2.Y++
		} else if(wall.X < p.X && wall.Y == p.Y) { //#4
			wp2.Y++
		} else if(wall.X > p.X && wall.Y == p.Y) { //#6
			wp1.X++
			wp2.Y++
			wp2.X++
		} else if(wall.X == p.X && wall.Y < p.Y) { //#8
			wp2.X++
		} else if( (wall.X < p.X && wall.Y < p.Y) || (wall.X > p.X && wall.Y > p.Y) ) { //#7 + #3
			wp1.X++
			wp2.Y++
		} else if( (wall.X > p.X && wall.Y < p.Y) || (wall.X < p.X && wall.Y > p.Y) ) { //#9 + #1
			wp2.X++
			wp2.Y++
		}

		//Compute angle for walls
		ag1 := GetAngle(wp1,pp)
		ag2 := GetAngle(wp2,pp)

		//Compute projection points
		ep1, ep2 := wp1, wp2
		ep1.X += math.Cos(ag1) * size
		ep1.Y += math.Sin(ag1) * size
		ep2.X += math.Cos(ag2) * size
		ep2.Y += math.Sin(ag2) * size

		//Create polygon of invisible space
		wallFov := polyclip.Polygon{{ep1, ep2, wp1, wp2}}

		//Mark invisible tiles
		for x := 0; x < p.Game.Width; x++ {
			for y := 0; y < p.Game.Height; y++ {
				if(p.FoV[x][y] && len(GetBox(x,y).Construct(polyclip.INTERSECTION, wallFov)) > 0) {
					p.FoV[x][y] = false;
				}
			}
		}

	}

}

//Player can see another player?
func (o *Player) CanSee(p *Player) bool {
	//Must be same game
	if(p.Game != o.Game) {
		return false
	}

	return o.FoV[p.X][p.Y]
}

//Respawn player, refill life, recompute FoV
func (o *Player) Spawn() {
	//Refill life
	o.Life = MAX_PLAYER_LIFE;

	//Place player
	o.X = rand.Intn(o.Game.Width);
	o.Y = rand.Intn(o.Game.Height);
	for(!o.Game.SpaceIsEmpty(o.X,o.Y,o)) {
		o.X = rand.Intn(o.Game.Width);
		o.Y = rand.Intn(o.Game.Height);
	}

	//Compute Field of View
	o.ComputeFoV()
}

//Process player's turn
func (o *Player) Step() {
	switch(o.Action) {
	//Action: Shoot bullet
	case ACTION_SHOOT_BULLET:
		bullet := Bullet{Way:o.Way,X:o.X,Y:o.Y,Player:o}
		o.Game.AddBullet(&bullet)
		o.Game.AddEvent(o.X,o.Y,EVENT_TYPE_SHOT)

	//Action: Throw granade
	case ACTION_THROW_GRANADE:
		granade := Granade{Way:o.Way,X:o.X,Y:o.Y,Player:o,Timer:GRANADE_FLY_TIME}
		o.Game.AddGranade(&granade)

	//Action: Move forward
	case ACTION_MOVE:
		if(o.Way == WAY_UP && o.Game.SpaceIsEmpty(o.X,o.Y-1,o)) {
			o.Y--;
		} else if(o.Way == WAY_RIGHT && o.Game.SpaceIsEmpty(o.X+1,o.Y,o)) {
			o.X++;
		} else if(o.Way == WAY_DOWN && o.Game.SpaceIsEmpty(o.X,o.Y+1,o)) {
			o.Y++;
		} else if(o.Way == WAY_LEFT && o.Game.SpaceIsEmpty(o.X-1,o.Y,o)) {
			o.X--;
		}

	//Action: Turn right
	case ACTION_TURN_RIGHT:
		o.Way++
		if(o.Way > WAY_LEFT) {
			o.Way = WAY_UP
		}

	//Action: Turn left
	case ACTION_TURN_LEFT:
		o.Way--
		if(o.Way < 0) {
			o.Way = WAY_LEFT
		}
	}

	//Reset action to NONE
	o.Action = ACTION_NONE
}
