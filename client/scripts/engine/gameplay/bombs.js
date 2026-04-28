const bombs = {}

class Bomb {
  constructor(id, x, y, target_x, target_y, mine = false) {
    this.object = build_bomb(id)
    this.which = id
    this.object.x = x
    this.object.y = y
    this.origin_x = x
    this.origin_y = y
    this.target_x = target_x
    this.target_y = target_y

    this.elapsed_time = 0
    this.mine = mine
    this.dead = false

    this.particle_emitter = new ParticleEmitter(
      this.object,
      bomb_data[id].trail.label,
      0.1,
    )

    add_object(this.object)
  }

  tick(deltaMS) {
    const data = bomb_data[this.which]
    const airtime = data.airtime

    this.particle_emitter.tick(deltaMS, !this.dead)
    this.elapsed_time += deltaMS

    if (!this.dead) {
      const t =
        airtime > 0 ? Math.min(this.elapsed_time / (airtime * 1000), 1) : 1

      this.object.angle += Math.sin(Math.PI * t) * deltaMS
      this.object.scale.set(1 / 128 + Math.sin(Math.PI * t) / 64)

      const dx = this.target_x - this.origin_x
      const dy = this.target_y - this.origin_y

      this.object.x = this.origin_x + dx * t
      this.object.y =
        this.origin_y - 2 * airtime * Math.sin(t * Math.PI) + dy * t
    }
  }

  kill(id) {
    this.dead = true
    this.object.visible = false
    this.particle_emitter.radiate(bomb_data[this.which].radius, 7)
    misc_layer.attach(this.particle_emitter.container)

    setTimeout(() => {
      delete bombs[id]
      this.object.destroy()
      this.particle_emitter.kill()
    }, 2000)
  }
}

function build_bomb(bomb_id) {
  let bomb = new PIXI.Container()

  for (let part of bomb_data[bomb_id].object) {
    let obj = build_object(part)
    bomb_layer.attach(obj)
    bomb.addChild(obj)
  }
  return bomb
}

function bomb_tick(deltaMS) {
  for (let id of Object.keys(bombs)) {
    const bomb = bombs[id]
    bomb.tick(deltaMS)

    if (
      !bomb.dead &&
      bomb.elapsed_time > bomb_data[bomb.which].airtime * 1000
    ) {
      bomb.kill(id)
    }
  }
}
