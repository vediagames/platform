package graphql

var CategoryPageLayouts = map[int][]string{
	1: { // "shooting"
		"first-person-shooter",
		"third-person-shooter",
	},
	2: { // "io"
		"hypercasual",
	},
	4: { // "board-and-puzzle"
		"card",
		"jigsaw",
		"mahjong",
		"word",
		"match-3",
		"math",
	},
	5: { // "multiplayer"
		"play-with-friends",
		"online",
	},
	6: { // "sports"
		"soccer",
		"basketball",
		"boxing",
		"football",
		"pool",
		"ball",
	},
	10: { // "driving"
		"bike",
		"motor-racing",
		"car",
		"truck",
		"parking",
	},
	11: { // "adventure"
		"build",
	},
	12: { // "beauty"
		"dress-up",
		"makeup",
	},
	13: { // "arcade"
		"match-3",
		"pool",
		"runner",
	},
}
