package atlas

import (
	"sync"
)

type Cell struct {
	vertices []Vertex
	biome    uint8

	Origin Point
	Tiles  []Tile `json:"tiles"`

	grid map[Point]*Tile
	Idx  int
	mu   sync.Mutex
}

func (cell *Cell) GetBiome() uint8 {
	return cell.biome
}

func NewCell(origin Point, idx int) *Cell {
	return &Cell{
		biome:    0,
		vertices: []Vertex{},
		Origin:   origin,
		Tiles:    []Tile{},
		grid:     make(map[Point]*Tile),
		Idx:      idx,
	}
}

func (cell *Cell) griddify() {
	for idx, tile := range cell.Tiles {
		cell.grid[tile.point()] = &cell.Tiles[idx]
	}
}

func (cell *Cell) addTile(tile Tile) {
	cell.Tiles = append(cell.Tiles, tile)
}

func (cell *Cell) GetAdjacentCells() []*Cell {
	seen := make(map[Point]bool)
	seen[cell.Origin] = true

	var cells []*Cell
	for _, vertex := range cell.vertices {
		for _, adjacentCell := range vertex.cells {
			if !seen[adjacentCell.Origin] {
				cells = append(cells, adjacentCell)
				seen[adjacentCell.Origin] = true
			}
		}
	}
	return cells
}
