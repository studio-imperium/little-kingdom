const objects = {}

class Character {
    constructor(hand, gear, angle, x, y) {
        this.object = build_character(hand, gear)
        this.object.kit = hand + gear.head + gear.body        
        this.object.angle = angle
        this.object.x = x
        this.object.y = y
        
        add_object(this.object)
    }

    update(hand, gear, angle, x, y) {
        if (hand + gear.head + gear.body != this.kit) {
            this.object.destroy()
            this.object = build_character(hand, gear)
            add_object(this.object)
        }
        this.object.angle = angle
        this.object.x = x
        this.object.y = y
    }
    
    attack_animation() {
    }
}

class Npc {
    constructor(id, angle, x, y) {
        this.object = build_npc(id) 
        this.object.angle = angle
        this.object.x = x
        this.object.y = y
        
        add_object(this.object)
    }

    update(angle, x, y) {
        this.object.angle = angle
        this.object.x = x
        this.object.y = y
    }
    
    attack_animation() {
    }
}