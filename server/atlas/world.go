package atlas

import (
	"math/rand"
	"sync"
)

type World struct {
	rnd          *rand.Rand
	points       []Point
	triangles    []Triangle
	cellByOrigin map[Point]*Cell

	Cells []*Cell `json:"cells"`
	Size  int     `json:"size"`
}

func newWorld(size int, density int, seed int64) *World {
	world := &World{
		rnd:          rand.New(rand.NewSource(seed)),
		points:       make([]Point, density),
		triangles:    []Triangle{},
		cellByOrigin: make(map[Point]*Cell),
		Cells:        []*Cell{},
		Size:         size,
	}
	world.triangulate()
	world.assignVertices()
	world.fillCells()

	return world
}

func (world *World) GetNearestCell(pt Point) *Cell {
	nearest := world.Cells[0].Origin
	nearestDist := distance(nearest, pt)

	for idx := range world.Cells {
		cell := world.Cells[idx]
		cellDist := distance(cell.Origin, pt)
		if nearestDist > cellDist {
			nearest = cell.Origin
			nearestDist = cellDist
		}
	}

	return world.cellByOrigin[nearest]
}

func (world *World) triangulate() {
	for idx := range world.points {
		world.points[idx] = Point{
			X: float64(world.Size) * world.rnd.Float64(),
			Y: float64(world.Size) * world.rnd.Float64(),
		}

		world.Cells = append(world.Cells, NewCell(world.points[idx]))
		cell := world.Cells[idx]
		world.cellByOrigin[cell.Origin] = cell
	}
	world.triangles = createTriangles(world.points)
}

func (world *World) assignVertices() {
	for _, triangle := range world.triangles {
		world.newVertices(triangle)
	}
}

func (world *World) fillCells() {
	wg := &sync.WaitGroup{}
	wg.Add(world.Size)

	assignCol := func(x int, world *World) {
		defer wg.Done()
		for y := 0; y < world.Size; y++ {
			tile := Tile{
				X:     x,
				Y:     y,
				Value: 0,
			}
			cell := world.GetNearestCell(tile.point())

			cell.mu.Lock()
			cell.addTile(tile)
			cell.mu.Unlock()
		}
	}
	for x := 0; x < world.Size; x++ {
		go assignCol(x, world)
	}
	wg.Wait()

	for _, cell := range world.Cells {
		cell.griddify()
	}
}
