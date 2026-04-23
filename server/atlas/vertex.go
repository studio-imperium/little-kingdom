package atlas

type Vertex struct {
	location Point
	cells    [3]*Cell
}

func (world *World) newVertices(triangle Triangle) {
	vertex := Vertex{
		location: triangle.center,
		cells: [3]*Cell{
			world.cellByOrigin[triangle.points[0]],
			world.cellByOrigin[triangle.points[1]],
			world.cellByOrigin[triangle.points[2]],
		},
	}

	for _, cell := range vertex.cells {
		if cell == nil {
			return
		}
	}
	for _, cell := range vertex.cells {
		cell.vertices = append(cell.vertices, vertex)
	}
}
