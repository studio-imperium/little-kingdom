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

    add_object(this.object)
  }

  update(x, y, angle, health, hand, head, body) {
    if (hand + head + body != this.kit) {
      this.object.destroy()
      this.object = build_character(hand, head, body)
      this.animator.set_object(this.object)
      this.interpolator.set_object(this.object)
      this.kit = hand + head + body
      this.object.x = x
      this.object.y = y
      this.object.angle = angle
      this.interpolator.object = this.object
      this.colorAnimator.object = this.object
      add_object(this.object)
    }
    this.object.health = health
    this.hand = hand
    this.head = head
    this.body = body
    this.interpolator.add_char_frame(x, y, angle)
  }

  damage() {
    this.colorAnimator.animate(0xffb3b3, 300)
  }

  kill(id) {
    this.colorAnimator.animate(0xff0000, 300)
    this.animator.animate(0, 0.2)
    this.interpolator.frames = []
    this.interpolator.last_frame = Date.now()

    setTimeout(() => {
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
