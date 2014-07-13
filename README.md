Super Turn-based Realtime Shooter Deluxe
=====

Turn-base/Realtime 3rd perspective multiplayer shooter game trough REST api.

*Currently frontend uses [sprites](http://spritedatabase.net/game/1868) from [Chaos Engine](http://en.wikipedia.org/wiki/The_Chaos_Engine) and all rights are reserved by their respective owners. If anyone have issues with this, just ping the project and we will remove the graphics imidiately.* __This will change when if project will be able to find pixel-art artist willing to create the sprites.__

Quick start
=====

```bash
	$ git clone ...
	$ go get github.com/akavel/polyclip-go 
	$ go run run.go
```

Open [http://localhost:8080](http://localhost:8080) in your browser (Hopefully you are up to date and use latest browser)

REST Api
=====

User is *currently* registered automatically with 1st game request

__{api path}__ is always game instance's apiPath

### GET {api path}/{username}/shoot
Shoots a bullet in next turn

### GET {api path}/{username}/throw
Throws a granade in next turn

### GET {api path}/{username}/go
Go forward in next turn

### GET {api path}/{username}/left
Turn left in next turn

### GET {api path}/{username}/right
Turn right in next turn

### Sample JSON response
```json
{
	"turn":16,
	"width":5,
	"height":5,
	"players":[
		{
			"name":"Bobesa",
			"x":1,
			"y":2,
			"way":0,
			"life":2,
			"granades":0,
			"bullets":0
		}
	],
	"walls":[
		{
			"id":33,
			"x":2,
			"y":1,
			"life":5
		}
	],
	"bullets":[
		{
			"id":22,
			"x":2,
			"y":2,
			"way":1
		}
	],
	"granades":[
		{
			"id":22,
			"x":2,
			"y":2,
			"way":1,
			"timer":3
		}
	],
	"fov":[
		[true,true,true,false,false],
		[true,true,true,false,false],
		[true,true,true,false,false],
		[true,true,false,false,false],
		[true,false,false,false,false]
	],
	"fragLimit":5,
	"events":[],
	"score":{
		"Bobesa":0
	}
}
```
