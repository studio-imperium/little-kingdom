package server

import "github.com/studio-imperium/atlas"

// Tiles
var DEEPWATER int8 = 0
var WATER int8 = 1
var GRASS int8 = 2
var STONE int8 = 3
var SAND int8 = 4
var SANDSTONE int8 = 5
var DRYGRASS int8 = 6
var RUBBLE int8 = 7
var DARKSTONE int8 = 8
var SNOW int8 = 9
var ICE int8 = 10

var OceanMap []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(WATER),
	),
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
	),
}

var IslandMap []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(GRASS),
	),
	atlas.NewBiome(
		atlas.NewFill(GRASS),
	),
	atlas.NewBiome(
		atlas.NewFill(GRASS),
	),
	atlas.NewBiome(
		atlas.NewFill(SAND),
	),
	atlas.NewBiome(
		atlas.NewFill(WATER),
	),
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
	),
}

var DesertMountainsMap []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(20, WATER, GRASS),
		atlas.NewSelectiveBorder(STONE, WATER),
	),
	atlas.NewBiome(
		atlas.NewPattern(3.2, GRASS, STONE),
		atlas.NewBorder(STONE),
	),
	atlas.NewBiome(
		atlas.NewPattern(5, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, SAND, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
	),
	atlas.NewBiome(
		atlas.NewCropCircle(40, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, RUBBLE, DARKSTONE, DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewBorder(DARKSTONE),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, RUBBLE, DARKSTONE, DARKSTONE, DARKSTONE, DARKSTONE, DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, DARKSTONE, DARKSTONE, SNOW, ICE),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewPattern(27, RUBBLE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewSelectiveBorder(DARKSTONE, RUBBLE),
		atlas.NewBorder(ICE),
		atlas.NewSelectiveExternalBorder(SNOW, ICE),
	),
	atlas.NewBiome(
		atlas.NewPattern(27, ICE, WATER),
		atlas.NewSelectiveBorder(WATER, ICE),
	),
	atlas.NewBiome(
		atlas.NewFill(WATER),
	),
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
	),
}

var Sandy []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(3.2, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, SAND, SANDSTONE, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, DRYGRASS),
	),
	atlas.NewBiome(
		atlas.NewCropCircle(40, DRYGRASS, SAND),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SAND, DRYGRASS),
		atlas.NewSelectiveBorder(SANDSTONE, SAND),
	),
}

var Snowy []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, RUBBLE, DARKSTONE, DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewBorder(DARKSTONE),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, RUBBLE, DARKSTONE, DARKSTONE, DARKSTONE, DARKSTONE, DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(40, RUBBLE, DARKSTONE, DARKSTONE, SNOW, ICE),
		atlas.NewSelectiveBorder(DARKSTONE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewBorder(SNOW),
	),
	atlas.NewBiome(
		atlas.NewPattern(27, RUBBLE, SNOW),
		atlas.NewSelectiveBorder(SNOW, ICE),
		atlas.NewSelectiveBorder(DARKSTONE, RUBBLE),
		atlas.NewBorder(ICE),
		atlas.NewSelectiveExternalBorder(SNOW, ICE),
	),
	atlas.NewBiome(
		atlas.NewPattern(27, ICE, WATER),
		atlas.NewSelectiveBorder(WATER, ICE),
	),
	atlas.NewBiome(
		atlas.NewFill(WATER),
	),
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
	),
}

var PatternDemo []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewPattern(5, WATER, DEEPWATER),
	),
}

var CropCircleDemo []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewCropCircle(5, WATER, DEEPWATER),
	),
}

var VoronoiDemo []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewVoronoi(100, WATER, DEEPWATER),
	),
}

// Islands
var Islands1 []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewVoronoi(100, GRASS, DEEPWATER),
		atlas.NewSelectiveExternalBorder(SAND, GRASS),
		atlas.NewSelectiveBorder(WATER, SAND),
		atlas.NewSelectiveBorder(SAND, GRASS),
	),
}

// Big islands
var BorderDemo2 []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
		atlas.NewBorder(WATER),
	),
}

// Big islands
var BorderDemo []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewBorder(DEEPWATER),
	),
}

// Funky islands
var FunkyDemo []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewPattern(7, DEEPWATER, GRASS),
		atlas.NewSelectiveExternalBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
	),
}

// Final demo
var Final []atlas.Biome = []atlas.Biome{
	atlas.NewBiome(
		atlas.NewFill(GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewBorder(DEEPWATER),
	),
	atlas.NewBiome(
		atlas.NewFill(GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewBorder(DEEPWATER),
	),
	atlas.NewBiome(
		atlas.NewFill(GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
		atlas.NewBorder(DEEPWATER),
	),
	atlas.NewBiome(
		atlas.NewVoronoi(10, GRASS, DEEPWATER),
		atlas.NewSelectiveExternalBorder(SAND, GRASS),
		atlas.NewSelectiveBorder(WATER, SAND),
		atlas.NewSelectiveBorder(SAND, GRASS),
	),
	atlas.NewBiome(
		atlas.NewPattern(7, DEEPWATER, GRASS),
		atlas.NewSelectiveExternalBorder(WATER, GRASS),
		atlas.NewSelectiveBorder(SAND, GRASS),
	),
	atlas.NewBiome(
		atlas.NewFill(DEEPWATER),
	),
}

func create_map() *atlas.World {
	world := atlas.NewWorld(200, 200, 21)
	world.Infect(Sandy, 0.2)
	return world
}
