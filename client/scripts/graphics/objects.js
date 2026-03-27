const object_layer = new PIXI.RenderLayer()
app.stage.addChild(object_layer)

function create_sprite(texture) {
    const sprite = new PIXI.Sprite(texture)
    
    const color = "black"
    const quality = 0.1
    const alpha = 1
    
    const outline = new PIXI.filters.OutlineFilter(
        {
            thickness: 2,
            color,
            quality,
            alpha
        }
    )
    const shadow = new PIXI.filters.OutlineFilter(
        {
            thickness: 2,
            color,
            quality,
            alpha: alpha/8
        }
    )
    
    sprite.filters = [
        outline,
        shadow,
    ]
    return sprite
}

function build_object(obj) {
    if (obj.type == "container") {
        const object = new PIXI.Container()
        const {x, y, angle, scale, label} = obj
        
        object.x = x ? x : 0
        object.y = y ? y : 0
        object.angle = angle ? angle : 0
        object.scale = scale ? scale : 1
        object.label = label

        for (let child of obj.children) {
            object.addChild(build_object(child))
        }
        
        return object
    }
    else {
        const texture = textures[obj.label]
        const {x, y, angle, scale} = obj
        const sprite = create_sprite(texture)

        sprite.x = x ? x : 0
        sprite.y = y ? y : 0
        sprite.angle = angle ? angle : 0
        sprite.scale = scale ? scale : 1

        return sprite
    }
}

function build_character(hand, gear) {
    let character = new PIXI.Container()

    if (hand && item_data[hand].hand) {
        character.addChild(build_object(item_data[hand].hand))
    }

    if (gear.body) {
        character.addChild(build_object(item_data[gear.body].equipped))
    }
    else {
        character.addChild(build_object(item_data["default_body"].equipped))
    }

    if (gear.head) {
        character.addChild(build_object(item_data[gear.head].equipped))
    }
    else {
        character.addChild(build_object(item_data["default_head"].equipped))
    }

    return character
}

function build_npc(npc_id) {
    let npc = new PIXI.Container()

    for (let bodypart of npc_data[npc_id].body) {
        npc.addChild(build_object(bodypart))
    }

    return npc
}

function add_object(object) {
    object.zIndex = object.scale.x
    app.stage.addChild(object)
    object_layer.attach(object)
    object_layer.sortRenderLayerChildren()
}
