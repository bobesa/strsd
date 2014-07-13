//Game type, workers and turn logic
package strsd

import (
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"math/rand"
	"fmt"
	"encoding/json"
)

type Game struct {
	Turn int `json:"turn"`
	Width int `json:"width"`
	Height int `json:"height"`
	Players []*Player `json:"players"`
	Walls []*Wall `json:"walls"`
	Bullets []*Bullet `json:"bullets"`
	Granades []*Granade `json:"granades"`
	FoV [][]bool `json:"fov,omitempty"`
	FragLimit int `json:"fragLimit"`
	Events []Event `json:"events"`
	Score map[string]int `json:"score,omitempty"`
	Queue chan *Request `json:"-"`

	OptionFogOfWar bool `json:"-"`
	OptionPersonalSpace bool `json:"-"`
	OptionFieldOfView bool `json:"-"`

	lastMapName string
	apiPath string
}

//Create new game instance with workers etc.
func MakeGame(apiPath string) *Game {
	var g = Game{
		Turn: 0,
		Players: make([]*Player,0),
		Walls: make([]*Wall,0),
		Bullets: make([]*Bullet,0),
		Granades: make([]*Granade,0),
		Width: 5,
		Height: 5,
		FragLimit: FRAG_LIMIT,
		Score: make(map[string]int),
		Queue: make(chan *Request),
		Events: make([]Event,0),

		OptionFogOfWar: DEFAULT_OPTION_FOG_OF_WAR,
		OptionPersonalSpace: DEFAULT_OPTION_PERSONAL_SPACE,
		OptionFieldOfView: DEFAULT_OPTION_FIELD_OF_VIEW,

		lastMapName: "",
		apiPath: apiPath,
	}

	//Load random map
	g.NewRandomMap();

	//Run request worker
	go func(){
		for {
			req := <-g.Queue
			if(req.Action != ACTION_NONE) {
				req.Player.Action = req.Action
			}
			req.Callback <- req.Player.Game.GetResponse(req.Player)
		}
	}()

	//Run turn timer
	go func(){
		for {
			time.Sleep(time.Second)
			g.Step()
		}
	}()

	return &g
}

//MakePlayer shortcut
func (g *Game) MakePlayer(name string) *Player {
	return MakePlayer(name,g)
}

//Function used for http.HandleFunc
func (g *Game) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.String()[len(g.apiPath):]

	if(query != "") {
		actions := strings.Split(query,"/");

		//Get/Create player
		player, exist := AllPlayers[actions[0]]
		if(!exist) {
			player = MakePlayer(actions[0],g);
			g.AddPlayer(player)
			AllPlayers[actions[0]] = player
		}

		//Process action
		action := ACTION_NONE
		if(len(actions) == 2) {
			switch(actions[1]){
			//Shoot bullet
			case "shoot":
				action = ACTION_SHOOT_BULLET
			//Throw granade
			case "throw":
				action = ACTION_THROW_GRANADE
			//Movement
			case "go":
				action = ACTION_MOVE
			case "right":
				action = ACTION_TURN_RIGHT
			case "left":
				action = ACTION_TURN_LEFT
			}
		}

		//Create callback and wait for processing
		callback := make(chan []byte)
		g.Queue <- &Request{Player:player,Action:action,Callback:callback}
		w.Write(<-callback)
	} else {
		w.Write([]byte("{}"))
	}
}

//Load random next map (that is different from last map)
func (g *Game) NewRandomMap() {
	files, err := ioutil.ReadDir("maps")
	if(err == nil) {
		file := files[rand.Intn(len(files))]
		for(len(files) > 1 && file.Name() == g.lastMapName) {
			file = files[rand.Intn(len(files))]
		}
		g.lastMapName = file.Name()
		g.NewMap(file.Name())
	}

}

//Load map by filename
func (g *Game) NewMap(filename string) {
	//Load Map
	data, err := ioutil.ReadFile("maps/"+filename)
	if(err == nil) {

		//Reset Game
		g.Turn = 0
		g.Walls = make([]*Wall,0)
		g.Bullets = make([]*Bullet,0)
		g.Granades = make([]*Granade,0)
		g.Events = make([]Event,0)

		//Load map
		lines := strings.Split(string(data),"\n")
		g.Height = len(lines) - 1
		g.Width = 2
		for y, line := range lines {
			if(y == 0) {
				fmt.Println("Number of frags: "+line) //TODO: parse number of frags
			} else {
				if(g.Width < len(line)-1) {
					g.Width = len(line)-1
				}
				for x := 0; x < len(line); x++ {
					switch(line[x]) {
					case 'W':
						g.AddWall(x,y-1)
					}
				}
			}
		}
		fmt.Println("Loaded map file '"+filename+"'.")

		//Reset players
		for _, player := range g.Players {
			player.Score = 0
			player.Spawn()
		}

	} else {
		fmt.Println("Unable to find map file '"+filename+"'! Loading another random map...")
		g.NewRandomMap()
	}

}

//Response generator (by modifying non-pointer game instance)
func (g Game) GetResponse(p *Player) []byte {
	if(g.OptionFieldOfView) {
		//Send only visible players if FoV is on
		visiblePlayers := make([]*Player,0)
		for _, e := range g.Players {
			if(e == p || p.CanSee(e)) {
				visiblePlayers = append(visiblePlayers,e)
			}
		}
		g.Players = visiblePlayers
		//Send fog of war?
		if(g.OptionFogOfWar) {
			g.FoV = p.FoV
		}
	}
	d, _ := json.Marshal(g)
	return d
}

//Check if no player is in place & place is in map bounds
func (g *Game) SpaceIsEmpty(x int, y int, p *Player) bool {
	//Bounds & Wall check
	if(!g.SpaceIsInBounds(x,y)) {
		return false
	}
	//Player check + Personal space check
	for _, o := range g.Players {
		if(o != p && ( (o.X == x && o.Y == y) || (
			g.OptionPersonalSpace && (
			(o.X-1 == x && o.Y == y)  ||
			(o.X+1 == x && o.Y == y)  ||
			(o.X == x && o.Y-1 == y)  ||
			(o.X == x && o.Y+1 == y) ) ) ) ){
			return false
		}
	}
	return true
}

//Check if place is in map bounds and no wall is hit
func (g *Game) SpaceIsInBounds(x int, y int) bool {
	//Map bounds check
	if(x < 0 || y < 0 || x >= g.Width || y >= g.Height) {
		return false
	}
	//Wall check
	for _, w := range g.Walls {
		if(x == w.X && y == w.Y) {
			return false
		}
	}
	return true
}

//Adds event to game
func (g *Game) AddEvent(x,y,t int) {
	g.Events = append(g.Events,Event{X:x,Y:y,Type:t})
}

//Adds bullet to game
func (g *Game) AddBullet(o *Bullet) {
	o.Id = MakeGUID()
	o.Step()
	g.Bullets = append(g.Bullets, o)
}

//Adds granade to game
func (g *Game) AddGranade(o *Granade) {
	o.Id = MakeGUID()
	o.Step()
	g.Granades = append(g.Granades, o)
}

//Adds wall to game map
func (g *Game) AddWall(x, y int) {
	g.Walls = append(g.Walls, &Wall{X:x,Y:y,Life:MAX_WALL_LIFE,Id:MakeGUID()})
}

//Adds player to game
func (g *Game) AddPlayer(p *Player) {
	g.Players = append(g.Players, p)
}

//Process of single turn
func (g *Game) Step() {
	g.Turn++

	g.Events = make([]Event,0)

	//Bullets
	visibleBullets := make([]*Bullet,0)
	for _, bullet := range g.Bullets {
		if(bullet.Step()) {
			visibleBullets = append(visibleBullets,bullet)
		}
	}
	g.Bullets = visibleBullets

	//Granades
	visibleGranades := make([]*Granade,0)
	for _, granade := range g.Granades {
		if(granade.Step()) {
			visibleGranades = append(visibleGranades,granade)
		}
	}
	g.Granades = visibleGranades

	//Players
	topScore := 0
	g.Score = make(map[string]int)
	for _, player := range g.Players {
		player.Step()
		player.ComputeFoV()
		g.Score[player.Name] = player.Score;
		if(topScore < player.Score) { topScore = player.Score; }
	}

	//Reset map
	if(g.FragLimit <= topScore) {
		g.NewRandomMap()
	}
}
