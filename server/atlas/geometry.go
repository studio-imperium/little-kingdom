package atlas

import "math"

func addPoint(triangles []Triangle, pt Point, pts *[]Point) []Triangle {
	var newTriangles []Triangle

	for _, triangle := range triangles {
		if triangle.withinCircumcircle(pt) {
			newTriangles = append(newTriangles, triangle.reform(pt, pts)...)
		} else {
			newTriangles = append(newTriangles, triangle)
		}
	}

	return newTriangles
}

func createTriangles(pts []Point) []Triangle {
	maxY := pts[0].Y
	maxX := pts[0].X
	for _, pt := range pts {
		maxY = math.Max(maxY, pt.Y)
		maxX = math.Max(maxX, pt.X)
	}

	addedPts := []Point{
		{0, 0},
		{maxX + 1, 0},
		{0, maxY + 1},
		{maxX + 1, maxY + 1},
	}

	triangles := []Triangle{
		newTriangle([3]Point{
			addedPts[0],
			addedPts[1],
			addedPts[2],
		}),
		newTriangle([3]Point{
			addedPts[1],
			addedPts[2],
			addedPts[3],
		}),
	}

	for _, pt := range pts {
		triangles = addPoint(triangles, pt, &addedPts)
		addedPts = append(addedPts, pt)
	}

	return triangles
}
