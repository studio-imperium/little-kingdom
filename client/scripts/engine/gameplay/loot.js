const loots = {}

class Loot {
  constructor(id, x, y) {
    this.object = build_loot(id)
    this.object.angle = Math.random() * 360
    this.object.x += x
    this.object.y += y
    this.last_update = Date.now()

    add_object(this.object)
  }

  update() {
    this.last_update = Date.now()
  }

  kill(id) {
    this.object.destroy()
    delete loots[id]
  }
}

function build_loot(loot_id) {
  let loot = new PIXI.Container()
  let sprite = item_data[loot_id].sprite
  let loot_blueprint = {
    x: -sprite.w / 2,
    y: -sprite.h / 2,
    angle: 0,
    scale: 1,
    label: loot_id,
  }

  loot_obj = build_object(loot_blueprint)
  loot_layer.attach(loot_obj)
  loot.addChild(loot_obj)

  return loot
}
