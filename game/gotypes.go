package game

var (
	PlayerBlack = Player{isBlack: true}
	PlayerWhite = Player{isBlack: false}
)

type Player struct {
	isBlack bool
}

func (p Player) other() Player {
	if p.isBlack {
		return PlayerWhite
	}

	return PlayerBlack
}

func (p Player) equals(other Player) bool {
	return p.isBlack == other.isBlack
}

type Point struct {
	Row int
	Col int
}

func (p Point) Neighbors() []Point {
	return []Point{
		{p.Row - 1, p.Col},
		{p.Row + 1, p.Col},
		{p.Row, p.Col - 1},
		{p.Row, p.Col + 1},
	}
}
