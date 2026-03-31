const SIZE = 6
const app = new PIXI.Application()
app.stage.scale = SIZE

async function init() {
    await app.init({
        background: "#1f1f1f",
        width: window.innerWidth,
        height: window.innerHeight,
        useContextAlpha: false,
        antialias: true,
        autoDensity: true,
        resolution: 1,
    })
    document.body.appendChild(app.canvas)

    app.stage.scale = 64
    app.canvas.style.imageRendering = "pixelated";
    app.canvas.style.imageRendering = "crisp-edges";
    
    await load_textures()
    await init_tiles()
    connect()
}
