const loots = {}

class Loot {
  constructor(id, x, y) {
    this.object = build_loot(id)
    this.dead = false
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
    this.dead = true

    setTimeout(() => {
      this.object.destroy()
      delete loots[id]
    }, 300)
  }
}

function build_loot(loot_id) {
  let loot = new PIXI.Container()
  let blueprint = item_data[loot_id].hand

  loot_obj = build_object(blueprint)
  loot_obj.angle = -blueprint.angle
  loot_obj.getChildAt(0).x = -5
  loot_obj.getChildAt(0).y = -5
  loot_layer.attach(loot_obj)
  loot.addChild(loot_obj)

  return loot
}
