const npcs = {}

class Npc {
  constructor(id, x, y, health) {
    this.object = build_npc(id)
    this.object.angle = 0
    this.object.x = x
    this.object.y = y
    this.animator = new Animator(this.object)
    this.interpolator = new Interpolator(this.object)
    this.colorAnimator = new ColorAnimator(this.object)
    this.which = id
    this.dying = false
    this.death_timeout = null
    this.confirmed_dead = false

    add_object(this.object)
  }

  revive() {
    // A fresh frame arrived after we started dying — cancel the pending
    // destroy and restore a sane visible state so the entity doesn't flicker.
    this.dying = false
    clearTimeout(this.death_timeout)
    this.death_timeout = null
    this.colorAnimator.active = false
    this.object.tint = 0xffffff
    // Base entity scale is 1/64 (see add_object and the animator's
    // object_scale / 64). The death animation shrinks toward 0, so restore the
    // base here — NOT 1, which would render the enemy 64x too big (gigantic).
    this.object.scale.set(1 / 64)
    this.animator.animation = null
    this.animator.timestamp = 0
  }

  update(x, y, health) {
    if (this.confirmed_dead) return
    if (this.dying) {
      this.revive()
    }
    this.object.health = health
    this.interpolator.add_npc_frame(x, y)
  }

  damage() {
    this.colorAnimator.animate(0xffb3b3, 300)
  }

  kill(id) {
    if (this.dying) return
    this.dying = true
    this.colorAnimator.animate(0xffb3b3, 30000)
    this.animator.animate(0, 0.2)
    this.interpolator.frames = []
    this.interpolator.last_frame = Date.now()

    this.death_timeout = setTimeout(() => {
      this.object.destroy()
      delete npcs[id]
    }, 300)
  }
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
