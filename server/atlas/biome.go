package atlas

type Biome struct {
	modifiers []Modifier
}

func NewBiome(mods ...Modifier) Biome {
	return Biome{
		modifiers: mods,
	}
}

func (biome *Biome) SetModifier(idx int, mod Modifier) {
	biome.modifiers[idx] = mod
}
