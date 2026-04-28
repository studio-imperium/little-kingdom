package engine

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
)

type Tile struct {
	X   uint16
	Y   uint16
	Val uint8
}

type Origin struct {
	X uint16
	Y uint16
}

type Cell struct {
	Idx    uint16
	Tiles  []Tile
	Origin Origin
	biome  uint8
	adj    []*Cell
}

type Map struct {
	size  uint16
	cells []*Cell
	tiles []uint8
}

func (c *Cell) GetAdjacentCells() []*Cell {
	return c.adj
}

func (m *Map) GetNearestCell(obj Object) *Cell {
	nearest := m.cells[0]
	nearestDist := float64(m.size)

	for _, cell := range m.cells {
		if dist := Dist(obj, cell); nearestDist > dist {
			nearest = cell
			nearestDist = dist
		}
	}

	return nearest
}

func LoadMap(name string) *Map {
	path := filepath.Join("..", "engine", "assets", "maps", name+".map")
	f, _ := os.Open(path)
	defer f.Close()

	var cellCount uint16

	binary.Read(f, binary.LittleEndian, &cellCount)

	m := &Map{
		cells: make([]*Cell, cellCount),
	}

	for i := range m.cells {
		m.cells[i] = &Cell{
			uint16(i),
			nil,
			Origin{0, 0},
			0,
			nil,
		}
	}

	for i := range m.cells {
		binary.Read(f, binary.LittleEndian, &m.cells[i].Origin.X)
		binary.Read(f, binary.LittleEndian, &m.cells[i].Origin.Y)
		binary.Read(f, binary.LittleEndian, &m.cells[i].biome)
		m.cells[i].Idx = uint16(i)

		var adj_len uint8
		binary.Read(f, binary.LittleEndian, &adj_len)

		var adj []*Cell = make([]*Cell, adj_len)
		for j := uint8(0); j < adj_len; j++ {
			var idx uint16
			binary.Read(f, binary.LittleEndian, &idx)
			adj[j] = m.cells[idx]
		}
		m.cells[i].adj = adj
	}

	binary.Read(f, binary.LittleEndian, &m.size)

	m.tiles = make([]uint8, int(m.size)*int(m.size))

	for y := uint16(0); y < m.size; y++ {
		for x := uint16(0); x < m.size; x++ {
			var tile uint8
			var idx uint16
			binary.Read(f, binary.LittleEndian, &tile)
			binary.Read(f, binary.LittleEndian, &idx)

			m.tiles[x+y*m.size] = tile
			m.cells[idx].Tiles = append(m.cells[idx].Tiles, Tile{x, y, tile})
		}
	}

	return m
}

func Dist(obj Object, cell *Cell) float64 {
	x := obj.GetX()
	y := obj.GetY()

	dx := float64(x) - float64(cell.Origin.X)
	dy := float64(y) - float64(cell.Origin.Y)
	return math.Hypot(dx, dy)
}
