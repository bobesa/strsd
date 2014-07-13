//Basic constants
var TILE_SIZE = 32;
var WAY_UP = 0,	WAY_RIGHT = 1, WAY_DOWN = 2, WAY_LEFT = 3;

//Load our game code after css/html is loaded
window.addEventListener("load",function(){
	var gUserName = "", gPlayers = {}, gBullets = {}, gGranades = {}, gWalls = {}, gEvents = [], gLastTurn = -1, gUserAvatar = 0;
	var vGame = document.getElementById("game"), vFoV = document.getElementById("fov"), vScoreboard = document.getElementById("scoreboard");

	//Simple trick to be able to go forEach with NodeList
	NodeList.prototype.forEach = Array.prototype.forEach;

	//Sets Player gfx tile
	var getTile = function(dom,x,y,sprite){
		if(sprite == 0) {
			sprite = gUserAvatar;
		} else if(sprite == gUserAvatar) {
			sprite = 0;
		}
		x *= -TILE_SIZE;
		y *= -TILE_SIZE;
		if(typeof sprite == "number") {
			switch(sprite) {
			case 3: y -= 180; break; //3
			case 5: y -= 360; break; //5
			case 4: x -= 340; break; //4
			case 2: x -= 340; y -= 180; break; //2
			case 1: x -= 340; y -= 360; break; //1
			}
		}
		dom.style.backgroundPosition = x + " " + y;
	};

	//Avatars
	var faceDoms = document.querySelectorAll(".face");
	faceDoms.forEach(function(dom,index){
		dom.addEventListener("click",function(){
			faceDoms.forEach(function(sDom){
				sDom.className = "face";
			});
			dom.className = "face active";
			gUserAvatar = index;
		});
	});

	//Controls
	var buttonShoot = document.createElement("span"), 
	buttonThrow = buttonShoot.cloneNode(true), 
	buttonLeft = buttonShoot.cloneNode(true), 
	buttonRight = buttonShoot.cloneNode(true), 
	buttonGo = buttonShoot.cloneNode(true);

	buttonShoot.innerHTML = "SHOOT (X)";
	buttonThrow.innerHTML = "GRANADE (C)";
	buttonLeft.innerHTML = "TURN LEFT (&larr;)";
	buttonRight.innerHTML = "TURN RIGHT (&rarr;)";
	buttonGo.innerHTML = "GO (&uarr;)";

	var clearControls = function(){
		buttonShoot.className = "";
		buttonThrow.className = "";
		buttonGo.className = "";
		buttonLeft.className = "";
		buttonRight.className = "";
	};

	buttonShoot.addEventListener("click",function(){ clearControls(); buttonShoot.className="active"; doGameRequest("shoot"); });
	buttonThrow.addEventListener("click",function(){ clearControls(); buttonThrow.className="active"; doGameRequest("throw"); });
	buttonLeft.addEventListener("click",function(){ clearControls(); buttonLeft.className="active"; doGameRequest("left"); });
	buttonRight.addEventListener("click",function(){ clearControls(); buttonRight.className="active"; doGameRequest("right"); });
	buttonGo.addEventListener("click",function(){ clearControls(); buttonGo.className="active"; doGameRequest("go"); });

	var controls = document.getElementById("controls");
	controls.appendChild(buttonShoot);
	controls.appendChild(buttonThrow);
	controls.appendChild(buttonLeft);
	controls.appendChild(buttonRight);
	controls.appendChild(buttonGo);

	document.addEventListener("keydown",function(e){
		if(e.which == 39) { //RIGHT
			buttonRight.click();
		} else if(e.which == 37) { //LEFT
			buttonLeft.click();
		} else if(e.which == 38) { //FORWARD
			buttonGo.click();
		} else if(e.which == 88) { //X SHOOT
			buttonShoot.click();
		} else if(e.which == 67) { //C THROW
			buttonThrow.click();
		}
	});

	//XHR
	var xhr = function(url,callback){
		var oReq = new XMLHttpRequest();
		oReq.onload = function(){
			callback(JSON.parse(this.responseText));
		};
		oReq.open("get", url, true);
		oReq.send();
	};

	//Login function
	var doLogin = function(username){
		gUserName = username;
		doGameRequest();
		document.getElementById("tabLogin").style.display = "none";
		document.getElementById("tabGame").style.display = "";
	};

	//Send a command you want to do
	var doGameRequest = function(action){
		if(action) {
			action = "/" + action;
		} else {
			action = "";
		}
		xhr("/api/"+gUserName+action,doParse);	
	};

	//Parse response from the XHR 
	var doParse = function(gamedata){
		//New Map
		if(gLastTurn > gamedata.turn) {
			gPlayers = {}; gBullets = {}; gGranades = {}; gWalls = {}; gEvents = []; gLastTurn = -1;
			vGame.innerHTML = "";
			vGame.appendChild(vFoV);
		}

		//New Turn
		if(gLastTurn != gamedata.turn) {
			gLastTurn = gamedata.turn;
			vGame.style.width = TILE_SIZE * gamedata.width + "px";
			vGame.style.height = TILE_SIZE * gamedata.height + "px";
			vGame.style.marginLeft = (-TILE_SIZE * 0.5 * gamedata.width) + "px";
			vGame.style.marginTop = (-TILE_SIZE * 0.5 * gamedata.height) + "px";
			clearControls();

			//Events
			gEvents.forEach(function(dom){
				dom.parentNode.removeChild(dom);				
			});
			gEvents = [];
			gamedata.events.forEach(function(e){
				var dom = document.createElement("div");
				dom.style.left = e.x * TILE_SIZE + "px"
				dom.style.top = e.y * TILE_SIZE + "px"
				switch(e.type) {
					case 0: //shot
						dom.className = "shot";
						break
					case 1: //explosion
						dom.className = "miniExplosion";
						break
				}
				vGame.appendChild(dom);
				gEvents.push(dom);
			});

			//Fov
			if(gamedata.fov) {
				vFoV.innerHTML = "";
				gamedata.fov.forEach(function(yValues,x){
					yValues.forEach(function(visible,y){
						if(!visible) {
							var fog = document.createElement("div");
							fog.className = "fog";
							fog.style.left = x * TILE_SIZE + "px";
							fog.style.top = y * TILE_SIZE + "px";
							vFoV.appendChild(fog);
						}
					});
				});
			}

			//Score
			var score = [];
			Object.keys(gamedata.score).forEach(function(playerName){
				score.push({name:playerName,score:gamedata.score[playerName]});
			});
			score.sort(function(a,b){
				if(a.score < b.score) {
					return +1;
				} else if(a.score > b.score) {
					return -1;
				} else {
					return 0;
				}
			});
			vScoreboard.innerHTML = "Scores: ";
			score.forEach(function(player,index){
				if(index != 0) {
					vScoreboard.innerHTML += " - ";
				}
				var name = player.name;
				if(name == gUserName) {
					name = "<u>"+name+"</u>";
				}
				vScoreboard.innerHTML += "<span><b>"+name+"</b> "+player.score+"/"+gamedata.fragLimit+"</span>";
			});
		}

		//Players
		foundPlayers = [];
		gamedata.players.forEach(function(player){
			var o;
			foundPlayers.push(player.name);
			if(gPlayers.hasOwnProperty(player.name)) {
				o = gPlayers[player.name];
			} else {
				o = {
					dom: document.createElement("div"),
					domName: document.createElement("div"),
					sprite: player.name == gUserName ? 0 : 1 + Math.floor(Math.random()*5)
				};
				o.dom.className = "user";
				o.domName.className = "userTitle";
				o.domName.innerHTML = player.name == gUserName ? "<b>"+player.name+"</b>" : player.name;
				gPlayers[player.name] = o;
				vGame.appendChild(o.dom);
				vGame.appendChild(o.domName);
			}

			o.dom.style.left = TILE_SIZE * player.x + "px";
			o.dom.style.top = TILE_SIZE * player.y + "px";
			o.domName.style.left = TILE_SIZE * player.x + "px";
			o.domName.style.top = TILE_SIZE * player.y + "px";
			switch(player.way) {
				case WAY_UP: getTile(o.dom,2,0,o.sprite); break;
				case WAY_RIGHT: getTile(o.dom,2,1,o.sprite);break;
				case WAY_DOWN: getTile(o.dom,4,2,o.sprite);break;
				case WAY_LEFT: getTile(o.dom,6,3,o.sprite);break;
			}
		});
		for(var b in gPlayers) {
			gPlayers[b].dom.className = foundPlayers.indexOf(b) == -1 ? "user faded" : "user";
			gPlayers[b].domName.className = foundPlayers.indexOf(b) == -1 ? "userTitle faded" : "userTitle";
		}

		//Bullets
		foundBullets = [];
		gamedata.bullets.forEach(function(bullet){
			var bName = "bullet_"+bullet.id;
			foundBullets.push(bName);
			var o;
			if(gBullets.hasOwnProperty(bName)) {
				o = gBullets[bName];
			} else {
				o = document.createElement("div");
				o.className = "bullet";
				gBullets[bName] = o;
				vGame.appendChild(o);
			}

			o.style.left = TILE_SIZE * bullet.x + "px";
			o.style.top = TILE_SIZE * bullet.y + "px";
			switch(bullet.way) {
				case WAY_UP: o.style.backgroundPosition = "-92px -152px"; break;
				case WAY_RIGHT: o.style.backgroundPosition = "-111px -152px";break;
				case WAY_DOWN: o.style.backgroundPosition = "-128px -152px";break;
				case WAY_LEFT: o.style.backgroundPosition = "-146px -152px";break;
			}
		});
		for(var b in gBullets) {
			if(foundBullets.indexOf(b) == -1) {
				gBullets[b].parentNode.removeChild(gBullets[b]);
				delete gBullets[b];
			}
		}

		//Granades
		foundGranades = [];
		gamedata.granades.forEach(function(granade){
			var gName = "granade_"+granade.id;
			foundGranades.push(gName);
			var o;
			if(gGranades.hasOwnProperty(gName)) {
				o = gGranades[gName];
			} else {
				o = document.createElement("div");
				o.className = "granade";
				gGranades[gName] = o;
				vGame.appendChild(o);
			}

			if(granade.timer == 0) {
				o.className = "granade explosion";				
			}

			o.style.left = TILE_SIZE * granade.x + "px";
			o.style.top = TILE_SIZE * granade.y + "px";
		});
		for(var b in gGranades) {
			if(foundGranades.indexOf(b) == -1) {
				gGranades[b].parentNode.removeChild(gGranades[b]);
				delete gGranades[b];
			}
		}

		//Walls
		foundWalls = [];
		gamedata.walls.forEach(function(wall){
			var gName = "wall_"+wall.id;
			foundWalls.push(gName);
			var o;
			if(gWalls.hasOwnProperty(gName)) {
				o = gWalls[gName];
			} else {
				o = document.createElement("div");
				o.className = "wall";
				gWalls[gName] = o;
				vGame.appendChild(o);
			}
			o.style.left = TILE_SIZE * wall.x + "px";
			o.style.top = TILE_SIZE * wall.y + "px";
		});
		for(var b in gWalls) {
			if(foundWalls.indexOf(b) == -1) {
				gWalls[b].parentNode.removeChild(gWalls[b]);
				delete gWalls[b];
			}
		}
	};

	//Do game request every 0.5 seconds
	setInterval(function(){
		if(gUserName != "") doGameRequest();
	},500);

	//Bind doLogin function to login button
	document.getElementById("buttonLogin").addEventListener("click",function(){
		var domInput = document.getElementById("textUsername");
		var name = domInput.value;
		if(name == "") {
			alert("You must type your name");
			domInput.focus();
		} else if(name.indexOf(" ") > -1) {
			alert("Name must have no spaces");	
			domInput.focus();		
		} else {
			doLogin(name);			
		}
	});
});
