const tile_layer = new PIXI.RenderLayer()
app.stage.addChild(tile_layer)

let default_vertices = get_vertices()
let tile_data = {}
let tile_map = {}

async function init_tiles() {
    let elapsed = 0 
    app.ticker.add((ticker) => {
        elapsed += ticker.deltaMS / 1000;
        const offsetX = Math.sin(elapsed * 0.8) * 0.5;
        const offsetY = Math.cos(elapsed * 0.6) * 0.5;

        for (const { tile, mesh, uvs } of Object.values(tile_map)) {
            if (tile == "lava" || tile == "water") {
                const uv_buffer = mesh.geometry.getBuffer("aUV")

                for (let i = 0; i < uv_buffer.data.length; i += 2) {
                    uv_buffer.data[i] = uvs.modified[i] + offsetX
                    uv_buffer.data[i + 1] = uvs.modified[i + 1] + offsetY
                }
                uv_buffer.update()
            }
        }
    })
}

function get_vertices() {
    return new Float32Array([
        0, 0,
        SIZE, 0,
        SIZE, SIZE,
        0, SIZE,
    ])
}
function randi() {
    return Math.floor(Math.random() * 1000)
}
function get_uvs(tile) {
    const { size } = tile_data[tile.tile]
    const subtiles = size / SIZE
    const a = 1/subtiles
    const x = (randi() % subtiles)/subtiles
    const y = (randi() % subtiles)/subtiles

    return new Float32Array([
        x,y,
        x+a,y,
        x+a,y+a,
        x,y+a,
    ])
}

function border_uvs(tile) {
    const { size } = tile_data[tile.tile]
    const subtiles = size / SIZE
    const a = 1/subtiles

    return new Float32Array([
        0,0,
        0+a,0,
        0+a,0+a,
        0,0+a,
    ])
}

function add_tile(x, y, tile_id, tile_key = `${x},${y}`) {
    const texture = textures[tile_id]
    const vertices = get_vertices()

    const tile_mesh = new PIXI.MeshSimple({
        texture,
        vertices,
        uvs: new Float32Array([0, 0, 1, 0, 1, 1, 0, 1]),
        indices: new Uint32Array([0, 1, 2, 0, 2, 3])
    })
    tile_mesh.x = x*SIZE
    tile_mesh.y = y*SIZE

    if (tile_map[tile_key]) {
        tile_map[tile_key].mesh.destroy()
    }
    tile_map[tile_key] = {
        "x" : x,
        "y" : y,
        "tile" : tile_id,
        "uvs" : {
            "modified" : new Float32Array([0, 0, 1, 0, 1, 1, 0, 1]),  // for water animations
            "base" : new Float32Array([0, 0, 1, 0, 1, 1, 0, 1]),      // for tile consistency
        },
        "mesh" : tile_mesh,
    }
    const tile = tile_map[tile_key]
    tile.uvs.base = get_uvs(tile)
    tile.uvs.modified.set(tile.uvs.base)
    tile.mesh.geometry.getBuffer("aUV").data.set(tile.uvs.base)

    app.stage.addChild(tile_mesh)
    tile_layer.attach(tile_mesh)
}

function reset_tiles() {
    for (let tile of Object.values(tile_map)) {
        const uvBuffer = tile.mesh.geometry.getBuffer("aUV")
        uvBuffer.data.set(tile.uvs.base)
        uvBuffer.update()
        tile.uvs.modified.set(tile.uvs.base)
        tile.mesh.vertices.set(default_vertices)
    }
}

function render_tiles() {
    reset_tiles()

    for (const bot_left of Object.values(tile_map)) {
        const [x,y] = [bot_left.x, bot_left.y]
        const [
            top_left,
            top_right,
            bot_right,
        ] = [
            tile_map[`${x},${y-1}`],
            tile_map[`${x+1},${y-1}`],
            tile_map[`${x+1},${y}`],
        ]

        if (top_left && top_right && bot_right) {
            smooth_vertices([top_left, top_right, bot_left, bot_right])
        }
    }
}

function smooth_vertices(tiles) {
    const count = {}
    const types = []

    for (let i = 0; i < tiles.length; i++) {
        let tile = tiles[i].tile
        if (count[tile]) {
            count[tile] += 1
        }
        else {
            count[tile] = 1
            types.push(tile)
            odd_one_out_idx = i
        }
    }

    if (types.length == 2 && (count[types[0]] == 1 || count[types[0]] == 3)) {
        let odd_tile = count[types[0]] == 1 ? types[0] : types[1]
        let odd_tile_idx = 0

        for (let i = 0; i < tiles.length; i++) {
            let tile = tiles[i].tile
            if (tile == odd_tile) {
                odd_tile_idx = i
                break
            }
        }

        let offset = 2
        let x_offset
        let y_offset

        if (odd_tile_idx == 0) {
            x_offset = -offset
            y_offset = -offset
        }
        else if (odd_tile_idx == 1) {
            x_offset = offset
            y_offset = -offset
        }
        else if (odd_tile_idx == 2) {
            x_offset = -offset
            y_offset = offset
        }
        else {
            x_offset = offset
            y_offset = offset
        }

        function shift_vertex(tile, idx) {
            const uvs = tile.uvs.modified
            const size = tile_data[tile.tile].size

            uvs[idx] += x_offset/size
            uvs[idx+1] += y_offset/size
            
            tile.mesh.geometry.getBuffer("aUV").data.set(uvs)
            tile.mesh.vertices[idx] += x_offset
            tile.mesh.vertices[idx+1] += y_offset
        }
        
        shift_vertex(tiles[0], 4)
        shift_vertex(tiles[1], 6)
        shift_vertex(tiles[2], 2)
        shift_vertex(tiles[3], 0)
    }
}
