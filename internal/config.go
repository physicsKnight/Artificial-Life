package internal

import "image/color"

var (
	screenW        = 1280
	screenH        = 900
	offsetX        = gridsize / 2
	offsetY        = gridsize / 2
	scaleX         = 1.0
	scaleY         = 1.0
	selectedX      = 0.0
	selectedY      = 0.0
	gridsize       = 130.0
	mouseOffsetX   = 638.0
	mouseOffsetY   = 482.0
	screenX1       = 0
	screenY1       = 0
	screenX2       = 0
	screenY2       = 0
	generations    = 0
	stopGeneration = -1
	steps          = 1
	running        = false
	wasRunning     = false
	exit           = false

	startLangton = `
          22222222
         2170140142
         2022222202
         272    212
         212    212
         202    212
         272    212
         21222222122222
         20710710711111
          2222222222222
    `

	startByl = `
          22
         2312
         2342
          25
    `

	startEvo = `
        02222222222222220
        20170170170170172
        27222222222222202
        21200000000000212
        20200000000000272
        27200000000000202
        21200000000000212
        20200000000000272
        27200000000000202
        21200000000000212
        20200000000000272
        27200000000000202
        21200000000000212
        21200000000000272
        21222222222222202
        21111111110410412
        02222222222222250
    `

	colors = []color.Color{
		color.RGBA{32, 5, 54, 255},     // black
		color.RGBA{0, 0, 255, 255},     // blue
		color.RGBA{255, 0, 0, 255},     // red
		color.RGBA{0, 255, 0, 255},     // green
		color.RGBA{255, 255, 0, 255},   // yellow
		color.RGBA{255, 0, 255, 255},   // magenta
		color.RGBA{255, 255, 255, 255}, // white
		color.RGBA{0, 255, 255, 255},   // cyan
		color.RGBA{160, 160, 160, 255}, // white
		color.RGBA{255, 200, 0, 255},   // orange
		color.RGBA{255, 175, 175, 255}, // pink
		color.RGBA{64, 64, 64, 255},    // dark gray
	}

	LangtonRules, BylRules, SDSRRules, EvoRules, SexyRules []string
	inputState                                             = "0"
	mode                                                   string

	automata      Automata
	LangtonObject Langton
	BylObject     Byl
	SDSRObject    SDSR
	EvoObject     Evo
	SexyObject    Sexy

	langtonMap = map[string]string{}
	bylMap     = map[string]string{}
	SDSRMap    = map[string]string{}
	evoMap     = map[string]string{}
	sexyMap    = map[string]string{}
)
