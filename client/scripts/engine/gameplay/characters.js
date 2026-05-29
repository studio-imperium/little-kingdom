const characters = {}

class Character {
  constructor(x, y, angle, health, hand, head, body) {
    this.object = build_character(hand, head, body)
    this.kit = (hand ? hand : 0) + head + body
    this.object.x = x
    this.object.y = y
    this.object.angle = angle
    this.object.health = health
    this.hand = hand
    this.head = head
    this.body = body
    this.animator = new Animator(this.object)
    this.interpolator = new Interpolator(this.object)
    this.colorAnimator = new ColorAnimator(this.object)
    this.dying = false
    this.death_timeout = null
    this.confirmed_dead = false

    add_object(this.object)
  }

  revive() {
    // A fresh frame arrived after we started dying — cancel the pending
    // destroy and restore a sane visible state so the player doesn't flicker.
    this.dying = false
    clearTimeout(this.death_timeout)
    this.death_timeout = null
    this.colorAnimator.active = false
    this.object.tint = 0xffffff
    // Base entity scale is 1/64 (see add_object and the animator's
    // object_scale / 64). Resetting to 1 would render the player 64x too big.
    this.object.scale.set(1 / 64)
    this.animator.animation = null
    this.animator.timestamp = 0
  }

  update(x, y, angle, health, hand, head, body) {
    if (this.confirmed_dead) return
    if (this.dying) {
      this.revive()
    }
    if (hand + head + body != this.kit) {
      let tmp_x = this.object.x
      let tmp_y = this.object.y

      this.object.destroy()
      this.object = build_character(hand, head, body)

      this.animator.set_object(this.object)
      if (this.interpolator) {
        this.interpolator.set_object(this.object)
      }
      this.object.x = tmp_x
      this.object.y = tmp_y
      this.object.angle = angle
      this.colorAnimator.object = this.object
      this.kit = hand + head + body
      add_object(this.object)
    }
    this.object.health = health
    this.hand = hand
    this.head = head
    this.body = body

    if (this.interpolator) {
      this.interpolator.add_char_frame(x, y, angle)
    } else if (
      Math.abs(this.object.x - x) > 1 &&
      Math.abs(this.object.y - y) > 1
    ) {
      this.object.x = x
      this.object.y = y
      this.object.angle = angle
    }
  }

  tick(deltaMS) {}

  damage() {
    this.colorAnimator.animate(0xffb3b3, 300)
  }

  kill(id) {
    if (this.dying) return
    this.dying = true
    this.colorAnimator.animate(0xff0000, 300)
    this.animator.animate(0, 0.2)
    this.interpolator.frames = []
    this.interpolator.last_frame = Date.now()

    this.death_timeout = setTimeout(() => {
      this.object.destroy()
      delete characters[id]
    }, 300)
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

function character_tick(deltaMS) {
  for (let id of Object.keys(characters)) {
    const char = characters[id]
    char.tick(deltaMS)
  }
  character.tick(deltaMS)
}
