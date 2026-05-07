package engine

import (
	"encoding/binary"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sync"
)

type Biome uint8

const (
	Beach Biome = iota
	Sandy3
	Sandy2
	Sandy
	Snowy
	Snowy2
	Glaciers
	Hot
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
	Idx        uint16
	Tiles      []Tile
	Origin     Origin
	biome      uint8
	Characters map[uint32]*Character
	npcs       map[uint32]*Npc
	adj        []*Cell
	m          *Map
}

type Map struct {
	size   uint16
	cells  []*Cell
	tiles  []uint8
	engine *Engine
	Mu     sync.Mutex
}

func (c *Cell) Spawn(npcs []SummonData) {
	for _, npcData := range npcs {
		id, npc := c.m.engine.SpawnNpc(npcData.ID, npcData.X+float32(c.Origin.X), npcData.Y+float32(c.Origin.Y))
		c.npcs[id] = npc
	}
}

func (c *Cell) Load() {
	for id, npc := range c.npcs {
		if npc.Dead {
			delete(c.npcs, id)
		}
	}
	for id, char := range c.Characters {
		if Dist(char, c) > 8 {
			delete(c.Characters, id)
		}
	}

	spawns, ok := biomeSpawns[c.biome]
	if len(c.Characters) == 0 && len(c.npcs) == 0 && ok {
		odds := float32(0.0)
		seed := rand.Float32()

		for i := 0; i < len(spawns); i++ {
			odds += spawns[i].Chance
			if seed <= odds {
				c.Spawn(spawns[i].Npcs)
				break
			}
		}
	}
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

var beachpoints = []Origin{}

func (engine *Engine) GetBeachpoint() (float32, float32) {
	r := rand.Int() % len(beachpoints)
	x := float32(beachpoints[r].X)
	y := float32(beachpoints[r].Y)
	return x, y
}

func (engine *Engine) LoadMap(name string) *Map {
	path := filepath.Join("..", "engine", "assets", "maps", name+".map")
	f, _ := os.Open(path)
	defer f.Close()

	var cellCount uint16

	binary.Read(f, binary.LittleEndian, &cellCount)

	m := &Map{
		cells:  make([]*Cell, cellCount),
		engine: engine,
	}

	for i := range m.cells {
		m.cells[i] = &Cell{
			uint16(i),
			nil,
			Origin{0, 0},
			0,
			make(map[uint32]*Character, 0),
			make(map[uint32]*Npc, 0),
			nil,
			m,
		}
	}

	for i := range m.cells {
		binary.Read(f, binary.LittleEndian, &m.cells[i].Origin.X)
		binary.Read(f, binary.LittleEndian, &m.cells[i].Origin.Y)
		binary.Read(f, binary.LittleEndian, &m.cells[i].biome)
		m.cells[i].Idx = uint16(i)

		if m.cells[i].biome == 23 {
			beachpoints = append(beachpoints, m.cells[i].Origin)
		}

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

	engine.Map = m
	return m
}

func Dist(obj Object, cell *Cell) float64 {
	x := obj.GetX()
	y := obj.GetY()

	dx := float64(x) - float64(cell.Origin.X)
	dy := float64(y) - float64(cell.Origin.Y)
	return math.Hypot(dx, dy)
}
