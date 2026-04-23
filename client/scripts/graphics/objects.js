const head_layer = new PIXI.RenderLayer()
const body_layer = new PIXI.RenderLayer()
const hand_layer = new PIXI.RenderLayer()
const misc_layer = new PIXI.RenderLayer()
app.stage.addChild(hand_layer, body_layer, head_layer, misc_layer)

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

function create_texture(texture) {
  const container = new PIXI.Container()
  const sprite = new PIXI.Sprite(texture)
  const scaleFactor = OBJECT_SIZE

  sprite.scale.set(scaleFactor)
  container.addChild(sprite)
  container.filters = [outline, shadow]

  const generated_texture = app.renderer.textureGenerator.generateTexture({
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

function build_object(obj) {
  if (obj.type == "container") {
    const object = new PIXI.Container()
    const { x, y, angle, scale, label } = obj

    object.x = x ? x * OBJECT_SIZE : 0
    object.y = y ? y * OBJECT_SIZE : 0
    object.angle = angle ? angle : 0
    object.scale = scale ? scale : 1
    object.label = label

    for (let child of obj.children) {
      object.addChild(build_object(child))
    }

    return object
  } else {
    if (!cache[obj.label]) {
      cache[obj.label] = create_texture(textures[obj.label])
    }
    const texture = cache[obj.label]
    const { x, y, angle, scale } = obj
    const sprite = create_sprite(texture)

    sprite.x = x ? x * OBJECT_SIZE : 0
    sprite.y = y ? y * OBJECT_SIZE : 0
    sprite.angle = angle ? angle : 0
    sprite.scale = scale ? scale : 1
    return sprite
  }
}

function build_character(hand, head, body) {
  let character = new PIXI.Container()
  character.sortableChildren = true

  if (hand && item_data[hand].hand) {
    const hand_obj = build_object(item_data[hand].hand)
    character.addChild(hand_obj)
    hand_layer.attach(hand_obj)
  }

  const body_obj = body
    ? build_object(item_data[body].equipped)
    : build_object(item_data[1].equipped)
  character.addChild(body_obj)
  body_layer.attach(body_obj)

  const head_obj = head
    ? build_object(item_data[head].equipped)
    : build_object(item_data[0].equipped)
  character.addChild(head_obj)
  head_layer.attach(head_obj)

  return character
}

function build_npc(npc_id) {
  let npc = new PIXI.Container()

  for (let bodypart of npc_data[npc_id].body) {
    let obj = build_object(bodypart)

    if (bodypart.label == "hand") {
      hand_layer.attach(obj)
    } else if (bodypart.label == "body") {
      body_layer.attach(obj)
    } else if (bodypart.label == "head") {
      head_layer.attach(obj)
    } else {
      misc_layer.attach(obj)
    }
    npc.addChild(obj)
  }

  return npc
}

function build_projectile(projectile_id) {
  let projectile = new PIXI.Container()

  for (let part of projectile_data[projectile_id].object) {
    let obj = build_object(part)

    body_layer.attach(obj)

    projectile.addChild(obj)
  }
  projectile.alpha = 0
  return projectile
}

function add_object(object) {
  object.scale.set(object.scale.x / 64)
  app.stage.addChild(object)
}
