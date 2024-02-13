package internal

import (
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"

	//"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func setupUI() *widget.Container {

	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(20)
	headerFace, _ := loadFont(30)

	//This creates the root container for this UI.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			// Padding defines how much space to put around the outside of the grid.
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			// Spacing defines how much space to put between each column and row
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x7f})))

	// construct a new container that serves as the root of the UI hierarchy
	menuPage := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x7f})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	menuPage.AddChild(newSeparator(widget.RowLayoutData{
		Stretch: true,
	}))
	menuPage.AddChild(newText("Mode", face))

	rootContainer.AddChild(headerContainer(headerFace, color.White))
	menuPage.AddChild(newSeparator(widget.RowLayoutData{
		Stretch: true,
	}))
	rootContainer.AddChild(menuPage)

	buttons := []*widget.Button{}
	buttons = append(buttons, createButton("SR", "langton", face, buttonImage))
	buttons = append(buttons, createButton("SDSR", "sdsr", face, buttonImage))
	buttons = append(buttons, createButton("Evo", "evo", face, buttonImage))
	buttons = append(buttons, createButton("Sexy", "sexy", face, buttonImage))
	buttons = append(buttons, createButton("Reset", "reset", face, buttonImage))

	// add the button as a child of the container
	for _, b := range buttons {
		menuPage.AddChild(b)
	}

	elements1 := []widget.RadioGroupElement{}
	for _, cb := range buttons {
		elements1 = append(elements1, cb)
	}
	widget.NewRadioGroup(widget.RadioGroupOpts.Elements(elements1...))

	menuPage.AddChild(newSeparator(widget.RowLayoutData{
		Stretch: true,
	}))

	menuPage.AddChild(newText("Input", face))
	menuPage.AddChild(newSeparator(widget.RowLayoutData{
		Stretch: true,
	}))
	menuPage.AddChild(textInput(face))

	toggles := []*widget.Button{}
	toggles = append(toggles, createButton("60", "60", face, buttonImage))
	toggles = append(toggles, createButton("45", "45", face, buttonImage))
	toggles = append(toggles, createButton("30", "30", face, buttonImage))
	toggles = append(toggles, createButton("15", "15", face, buttonImage))

	for _, t := range toggles {
		menuPage.AddChild(t)
	}

	elements := []widget.RadioGroupElement{}
	for _, cb := range toggles {
		elements = append(elements, cb)
	}
	widget.NewRadioGroup(widget.RadioGroupOpts.Elements(elements...))

	menuPage.AddChild(newSeparator(widget.RowLayoutData{
		Stretch: true,
	}))

	footerContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:  25,
			Right: 25,
		}),
	)))
	rootContainer.AddChild(footerContainer)
	closeButton := createButton("Exit", "exit", face, buttonImage)
	menuPage.AddChild(closeButton)

	return rootContainer
}

func newText(s string, face font.Face) *widget.Text {
	return widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
			Stretch:  false,
		})),
		widget.TextOpts.Text(s, face, color.White))
}

func createButton(text, modeName string, face font.Face, buttonImage *widget.ButtonImage) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			})),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(text, face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			handleButtonClick(modeName)
		}),
	)
}

func handleButtonClick(modeName string) {
	switch modeName {
	case "15", "30", "45", "60":
		fps, _ := strconv.Atoi(modeName)
		ebiten.SetMaxTPS(fps)
	case "reset":
		generations = 0
		clear(grid)
	case "exit":
		exit = true
	default:
		mode = modeName
		automata = getAutomata(mode)
	}
}

func textInput(face font.Face) *widget.TextInput {
	// construct a standard textinput widget
	return widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),

		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.White,
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.White,
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(8)),
		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),
		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Stop at generation: "),

		//This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			num, err := strconv.Atoi(args.InputText)
			if err == nil {
				stopGeneration = num
			}
		}),
	)
}

func headerContainer(face font.Face, col color.Color) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(15))),
	)

	c.AddChild(header("SETTINGS", face,
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	))

	c2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
		)),
	)
	c.AddChild(c2)

	return c
}

func header(label string, face font.Face, opts ...widget.ContainerOpt) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(append(opts, []widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x7f})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Left:   25,
			Right:  25,
			Top:    4,
			Bottom: 4,
		}))),
	}...)...)

	c.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text(label, face, color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	return c
}

func newSeparator(ld interface{}) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(color.NRGBA{170, 170, 180, 255})),
	))

	return c
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
