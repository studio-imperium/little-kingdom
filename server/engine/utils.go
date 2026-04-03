package engine

import (
	"math"
	"math/rand/v2"
)

type Entity interface {
	GetX() float32
	GetY() float32
	Damage(uint16)
}

type Object interface {
	GetX() float32
	GetY() float32
}

func NearbyPoint(obj Object, within float32) *Point {
	offsetX := (rand.Float32() * within * 2) - within
	offsetY := (rand.Float32() * within * 2) - within

	return &Point{
		x: obj.GetX() + offsetX,
		y: obj.GetY() + offsetY,
	}
}

func Distance(obj1 Object, obj2 Object) float64 {
	dx := obj1.GetX() - obj2.GetX()
	dy := obj1.GetY() - obj2.GetY()
	return math.Sqrt(math.Pow(float64(dx), 2) + math.Pow(float64(dy), 2))
}

type Point struct {
	x float32
	y float32
}

func (p Point) GetX() float32 { return p.x }
func (p Point) GetY() float32 { return p.y }
