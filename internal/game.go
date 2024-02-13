package internal

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func worldToScreen(worldX, worldY float64, screenX, screenY *int) {
	*screenX = int((worldX - offsetX) * scaleX)
	*screenY = int((worldY - offsetY) * scaleY)
}

func screenToWorld(screenX, screenY int, worldX, worldY *float64) {
	*worldX = float64(screenX)/scaleX + offsetX
	*worldY = float64(screenY)/scaleY + offsetY
}

var numberKeyCodes = []ebiten.Key{ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9, ebiten.KeyA, ebiten.KeyB, ebiten.KeyN}

// Load a font with size 15
var f, _ = loadFont(15)

// Game struct holds the game state and user input information
type Game struct {
	mouseX, mouseY                  int
	leftMousePressed, leftMouseHeld bool
	mouseReleased, rightPressed     bool
	downPressed, upPressed          bool
	spacePressed, enterPressed      bool
	escapePressed, leftPressed      bool
	keyExpand, keyMinimise          bool
	menu                            bool
	inputMode                       bool
	speed                           int
	ui                              *ebitenui.UI
}

func Run() {
	mode = "langton"
	offsetX = float64(-screenW) / 2
	offsetY = float64(-screenH) / 2
	worldToScreen(0, 0, &screenX1, &screenY1)
	worldToScreen(gridsize, gridsize, &screenX2, &screenY2)

	// Construct the UI
	ui := ebitenui.UI{
		Container: setupUI(),
	}
	game := Game{
		ui:    &ui,
		speed: 60,
	}

	// Set the maximum update rate, window size, and window title
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Langton Loops")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Initialize rules and the grid
	initRules()
	initGraph(false, false)

	// Start the game loop
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// handleKeyPress is a utility function that handles key press and release events
// It takes a key, a pointer to a boolean to track the key press state,
// and onPressed/onReleased functions that will be called when the key is pressed or released
func handleKeyPress(key ebiten.Key, pressed *bool, onPressed, onReleased func()) {
	// Check if the key is currently pressed
	if ebiten.IsKeyPressed(key) {
		// If the key was not already pressed, update the state and call onPressed if provided
		if !*pressed {
			*pressed = true
			if onPressed != nil {
				onPressed()
			}
		}
	} else {
		// If the key is not pressed but was previously pressed, update the state and call onReleased if provided
		if *pressed {
			*pressed = false
			if onReleased != nil {
				onReleased()
			}
		}
	}
}

func (g *Game) Update() error {
	// Exit the application if the exit flag is set
	if exit {
		return errors.New("closing the application")
	}
	// If the game is running, process the next generation of cells in the automaton
	if running {
		for i := 0; i < steps; i++ {
			processNextGen(mode, automata)
		}
	}

	// If not in the menu, handle user input and update the game state
	if !g.menu {
		// Get the current cursor position
		g.mouseX, g.mouseY = ebiten.CursorPosition()
		x, y := float64(g.mouseX), float64(g.mouseY)
		g.leftMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		// mouse button is being held
		if g.leftMousePressed {
			// if we're holding it for the first time
			// then record the initial mouse position
			if !g.leftMouseHeld {
				mouseOffsetX = x
				mouseOffsetY = y

				g.leftMouseHeld = true
				g.mouseReleased = false
			}
			// no need for else since we only need
			// the initial pos
		} else {
			// mouse button was just released
			if g.leftMouseHeld {
				g.mouseReleased = true
			}

			g.leftMouseHeld = false
		}

		// x, y is the initial pos
		// mouseoffset is the current pos
		if g.leftMouseHeld {
			offsetX -= (x - mouseOffsetX) / scaleX
			offsetY -= (y - mouseOffsetY) / scaleY
			mouseOffsetX = x
			mouseOffsetY = y
		}

		// Calculate the mouse coordinates before zooming
		var mouseX_beforeZoom, mouseY_beforeZoom float64
		screenToWorld(g.mouseX, g.mouseY, &mouseX_beforeZoom, &mouseY_beforeZoom)

		// Get the scroll amount from the mouse wheel
		_, scroll := ebiten.Wheel()
		// Zoom in or out based on the scroll amount
		if scroll < 0 {
			scaleX *= 1.01
			scaleY *= 1.01
		}
		if scroll > 0 {
			scaleX *= 0.99
			scaleY *= 0.99
		}

		// Handle key press events for various game functions
		handleKeyPress(ebiten.KeyEqual, &g.keyExpand, nil, func() {
			gridsize += 2
			initGraph(true, true)
		})
		handleKeyPress(ebiten.KeyMinus, &g.keyMinimise, nil, func() {
			gridsize -= 2
			initGraph(true, false)
		})
		handleKeyPress(ebiten.KeySpace, &g.spacePressed, nil, func() {
			running = !running
		})
		handleKeyPress(ebiten.KeyLeft, &g.leftPressed, nil, func() {
			if !running {
				// Swap the references of grid and gridCopy
				temp := grid
				grid = gridCopy
				gridCopy = temp
				generations--
			}
		})
		handleKeyPress(ebiten.KeyRight, &g.rightPressed, nil, func() {
			if !running {
				processNextGen(mode, automata)
			}
		})
		handleKeyPress(ebiten.KeyUp, &g.upPressed, nil, func() {
			steps++
		})
		handleKeyPress(ebiten.KeyDown, &g.downPressed, nil, func() {
			if steps > 1 {
				steps--
			}
		})

		// Check if any of the number keys are pressed
		for i, keyCode := range numberKeyCodes {
			if ebiten.IsKeyPressed(keyCode) {
				inputState = strconv.Itoa(i)
			}
		}

		var mouseX_afterZoom, mouseY_afterZoom float64
		screenToWorld(g.mouseX, g.mouseY, &mouseX_afterZoom, &mouseY_afterZoom)

		offsetX += (mouseX_beforeZoom - mouseX_afterZoom)
		offsetY += (mouseY_beforeZoom - mouseY_afterZoom)

		// Calculate the mouse coordinates after zooming
		if g.mouseReleased {
			// selected is in world space
			selectedX = mouseX_afterZoom
			selectedY = mouseY_afterZoom
			g.mouseReleased = false
		}

		handleKeyPress(ebiten.KeyEnter, &g.enterPressed, nil, func() {
			// If a cell is selected and the Enter key is pressed, update its value
			x, y := int(math.Floor(selectedX)), int(math.Floor(selectedY))
			if selectedX >= 0 && selectedX <= gridsize &&
				selectedY >= 0 && selectedY <= gridsize {
				if inputState == "12" {
					addConfiguration(y, x, automata)
				} else {
					grid[y][x].val = inputState
				}
			}
		})
	} else {
		// Run the menu screen updates
		g.ui.Update()
	}

	// we always want to switch to and from the menu
	handleKeyPress(ebiten.KeyEscape, &g.escapePressed, nil, func() {
		g.menu = !g.menu

		if g.menu {
			wasRunning = running
			running = false
		} else {
			running = wasRunning
		}
	})

	return nil
}

// draw graphics to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(color.RGBA{0, 0, 0, 0xff})

	// Calculate the world coordinates of the top-left and bottom-right corners of the screen
	var worldLeft, worldTop, worldRight, worldBottom float64
	screenToWorld(0, 0, &worldLeft, &worldTop)
	screenToWorld(screenW, screenH, &worldRight, &worldBottom)

	// Set the cell size on screen based on the current zoom level
	w, h := 1*scaleX, 1*scaleY

	// Loop through each cell in the grid and draw them on the screen
	for row := 0.0; row < gridsize; row++ {
		for col := 0.0; col < gridsize; col++ {
			// Get the cell value and corresponding color
			cstr := grid[int(col)][int(row)].val
			c, err := strconv.Atoi(cstr)
			if err != nil {
			}
			clr := colors[c]

			// Convert the cell's world coordinates to screen coordinates
			var screen_x1, screen_y1 int
			worldToScreen(float64(row), float64(col), &screen_x1, &screen_y1)

			// Draw the cell rectangle on the screen
			ebitenutil.DrawRect(screen, float64(screen_x1), float64(screen_y1), w, h, clr)
		}
	}

	// Display the current generation and mode on the screen
	text.Draw(screen, fmt.Sprintf("Generation: %d", generations), f, 5, 15, color.White)
	text.Draw(screen, fmt.Sprintf("Mode: %s", mode), f, 5, 30, color.White)
	text.Draw(screen, fmt.Sprintf("Steps: %d", steps), f, 5, 45, color.White)

	// If the menu is active, draw the UI
	if g.menu {
		g.ui.Draw(screen)
	}
}

// Set the layout for the game screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}
