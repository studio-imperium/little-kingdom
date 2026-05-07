const bomb_layer = new PIXI.RenderLayer()
const head_layer = new PIXI.RenderLayer()
const body_layer = new PIXI.RenderLayer()
const hand_layer = new PIXI.RenderLayer()
const loot_layer = new PIXI.RenderLayer()
const misc_layer = new PIXI.RenderLayer()
app.stage.addChild(
  misc_layer,
  loot_layer,
  hand_layer,
  body_layer,
  head_layer,
  bomb_layer,
)

const outline = new PIXI.filters.OutlineFilter({
  thickness: 3,
  color: "black",
  quality: 0.1,
  alpha: 1,
})
const shadow = new PIXI.filters.OutlineFilter({
  thickness: 3,
  color: "black",
  quality: 0.1,
  alpha: 0.125,
})
const cache = {}

function create_texture(texture, do_outline = true, renderer = app.renderer) {
  const container = new PIXI.Container()
  const sprite = new PIXI.Sprite(texture)
  const scaleFactor = OBJECT_SIZE

  sprite.scale.set(scaleFactor)
  container.addChild(sprite)

  if (do_outline) {
    container.filters = [outline, shadow]
  }

  const generated_texture = renderer.textureGenerator.generateTexture({
    target: container,
    resolution: 1,
    antialias: false,
    textureSourceOptions: {
      scaleMode: "nearest",
    },
  })

  generated_texture.source.scaleMode = "nearest"
  return generated_texture
}

function create_sprite(texture) {
  const sprite = new PIXI.Sprite(texture)
  return sprite
}

function build_object(obj, renderer = app.renderer, texture_cache = cache) {
  if (obj.type == "container") {
    const object = new PIXI.Container()
    const { x, y, angle, scale, label } = obj

    object.x = x ? x * OBJECT_SIZE : 0
    object.y = y ? y * OBJECT_SIZE : 0
    object.angle = angle ? angle : 0
    object.scale = scale ? scale : 1
    object.label = label

    for (let child of obj.children) {
      var thing = build_object(child, renderer, texture_cache)
      object.addChild(thing)
    }

    return object
  } else {
    if (!texture_cache[obj.label]) {
      texture_cache[obj.label] = create_texture(
        textures[obj.label],
        obj.outline === undefined ? true : obj.outline,
        renderer,
      )
    }
    const texture = texture_cache[obj.label]
    const { x, y, angle, scale } = obj
    const sprite = create_sprite(texture)

    sprite.x = x ? x * OBJECT_SIZE : 0
    sprite.y = y ? y * OBJECT_SIZE : 0
    sprite.angle = angle ? angle : 0
    sprite.scale = scale ? scale : 1
    return sprite
  }
}

function add_object(object) {
  object.scale.set(object.scale.x / 64)
  app.stage.addChild(object)
}
