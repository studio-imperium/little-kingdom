package atlas

func NewWorld(size int, density int, seed int64) *World {
	world := newWorld(size, density, seed)
	return world
}

func NewTemplateWorld(size int) *World {
	world := newWorld(size, size, 0)
	return world
}
