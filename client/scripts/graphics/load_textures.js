let textures
const spritesheet_data = {
    frames : {},
    meta : {
        image: "assets.png",
        format: "RGBA8888",
        size: { w: 512, h: 512},
        scale: 1
    }
}

async function load_textures() {
    await load_tiles()
    await load_items()
    await load_npcs()

    const assets_texture = await PIXI.Assets.load("assets/assets.png")
    const sheet = new PIXI.Spritesheet(assets_texture, spritesheet_data)
    sheet.parse()
    sheet.textureSource.source.scaleMode = "nearest"
    textures = sheet.textures

    console.log("Loaded " + Object.keys(textures).length + " textures")
}

async function load_tiles() {
    const tiles_json = await (await fetch("/assets/tiles.json")).json()

    for (let {id, size, x, y} of tiles_json) {
        tile_data[id] = {size,x,y}
        spritesheet_data.frames[id] = {
            frame: { x: x+8, y: y+8, w: size, h: size },
            sourceSize: {w: size, h: size},
            spriteSourceSize: {x: 0, y: 0, w: size, h: size},
        }
    }
}

function load_recursive(obj) {
    if (obj.type == "container") {
        for (let child of obj.children) {
            load_recursive(child)
        }
    }
    else {
        const {w, h} = obj.sprite
        spritesheet_data.frames[obj.label] = {
            frame: obj.sprite,
            sourceSize: {w, h},
            spriteSourceSize: {x: 0, y: 0, w, h},
        }
    }
}

async function load_items() {
    const items_json = await (await fetch("/assets/items.json")).json()

    for (let id of Object.keys(items_json)) {
        let {sprite, hand, equipped} = items_json[id]
        
        spritesheet_data.frames[id] = {
            frame: sprite,
            sourceSize: {w: sprite.w, h: sprite.h},
            spriteSourceSize: {x: 0, y: 0, w: sprite.w, h: sprite.h},
        }

        if (equipped && Object.keys(equipped).length) {
            load_recursive(equipped)
        }
        if (hand && Object.keys(hand).length) {
            load_recursive(hand)
        }
    }
    item_data = items_json
}

async function load_npcs() {
    const npc_json = await (await fetch("/assets/npcs.json")).json()

    for (let id of Object.keys(npc_json)) {
        let {body} = npc_json[id]
        
        for (let bodypart of body) {
            load_recursive(bodypart)
        }
    }
    npc_data = npc_json
}
