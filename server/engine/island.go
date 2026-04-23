package engine

import "server/atlas"

// Tiles
var WATER int8 = 1
var GRASS int8 = 2
var WOOD int8 = 3
var STONE int8 = 4
var DRYGRASS int8 = 5
var SAND int8 = 6
var SANDSTONE int8 = 7
var COLDGRASS int8 = 8
var SNOW int8 = 9
var ICE int8 = 10
var GRAVEL int8 = 11
var RUIN int8 = 12
var LAVA int8 = 13

var Beach []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(GRASS),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(6, GRASS, WATER, SAND),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewBorder(SAND),
		atlas.NewSelectiveExternalBorder(SAND, SAND),
	),
	atlas.NewBiome(
		atlas.NewPattern(6.1, GRASS, WATER),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewSelectiveBorder(SAND, WATER),
	),
	atlas.NewBiome(
		atlas.NewFill(WATER),
	),
}
var Sandy []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(7, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
	),
	atlas.NewBiome(
		atlas.NewCropCircle(3.5, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(20, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
	),
}
var Sandy2 []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(3.5, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(20, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
	),
}
var Sandy3 []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(3.5, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(20, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
	),
	atlas.NewBiome(
		atlas.NewCropCircle(7, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, SANDSTONE),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
	),
}
var Snowy []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(3, ICE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(10, ICE, COLDGRASS),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewSelectiveBorder(SNOW, ICE),
	),
}
var Snowy2 []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewPattern(10, ICE, COLDGRASS),
		atlas.NewSelectiveBorder(SNOW, ICE),
	),
	atlas.NewBiome(
		atlas.NewPattern(4.1, ICE, COLDGRASS),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewSelectiveBorder(SNOW, COLDGRASS),
	),
}
var Glaciers []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewVoronoi(40, STONE, STONE, STONE, STONE, STONE, STONE, STONE, SNOW),
		atlas.NewSelectiveBorder(STONE, SNOW),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, STONE, STONE, STONE, SNOW, ICE),
		atlas.NewSelectiveBorder(STONE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewPattern(27, STONE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewSelectiveBorder(STONE, STONE),
		atlas.NewBorder(ICE),
		atlas.NewSelectiveExternalBorder(SNOW, ICE),
	),
}
var Hot []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewPattern(3.3, LAVA, GRAVEL),
	),
	atlas.NewBiome(
		atlas.NewPattern(15, LAVA, RUIN),
		atlas.NewSelectiveBorder(GRAVEL, LAVA),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(20, RUIN, LAVA),
		atlas.NewSelectiveBorder(GRAVEL, LAVA),
	),
}

func AppendBiomes(biomes ...[]atlas.Biome) []atlas.Biome {
	if len(biomes) == 1 {
		return biomes[0]
	} else {
		return append(biomes[0], AppendBiomes(biomes[1:]...)...)
	}
}

var Island = AppendBiomes(Sandy, Sandy2, Sandy3, Beach)

func CreateIsland(size int) *atlas.World {
	world := atlas.NewWorld(size, 2000, 11)
	world.Infect(Island, 1)

	return world
}
