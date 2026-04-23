package atlas

type Tile struct {
	X     int  `json:"x"`
	Y     int  `json:"y"`
	Value int8 `json:"value"`
}

func (tile Tile) point() Point {
	return Point{
		float64(tile.X),
		float64(tile.Y),
	}
}
