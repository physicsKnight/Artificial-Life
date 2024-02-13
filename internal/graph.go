package internal

import (
	//"fmt"
	"strings"
	"sync"
	"unicode"
)

var (
	grid               [][]Cell
	gridCopy           [][]Cell
	currentMode, start string
	update             func(Cell, []string) string
	firstCell          []string
	chunks, goroutines int
	wg                 sync.WaitGroup
)

// inGrid checks if the given row and column are within the bounds of the grid.
func inGrid(r, c int) bool {
	return r >= 0 && r < int(gridsize) && c >= 0 && c < int(gridsize)
}

func onEdge(r, c int) bool {
	return r < 0 || r >= int(gridsize) || c < 0 || c >= int(gridsize)
}

// clear sets all the cells in the grid to "0" and adds the automata configuration to the grid.
func clear(cells [][]Cell) {
	for i, row := range cells {
		for j := range row {
			cells[i][j].val = "0"
		}
	}
	addConfiguration(int(gridsize)/2-5, int(gridsize)/2-10, automata)
}

// initCells initializes the cells of the grid, setting their values to "0", row, column
// and assigning their neighbors.
func initCells(cells [][]Cell) {
	for i := range cells {
		for j := range cells[i] {
			cells[i][j].val = "0"
			cells[i][j].row = i
			cells[i][j].col = j
			if i > 0 {
				cells[i][j].neighbours[0][0] = &cells[i-1][j] // top neighbour
			}
			if j < len(cells[i])-1 {
				cells[i][j].neighbours[1][0] = &cells[i][j+1] // right neighbour
			}
			if i < len(cells)-1 {
				cells[i][j].neighbours[2][0] = &cells[i+1][j] // bottom neighbour
			}
			if j > 0 {
				cells[i][j].neighbours[3][0] = &cells[i][j-1] // left neighbour
			}
		}
	}
}

// initGraph initializes the grid, gridCopy and automata based on the given mode, and
// optionally shifts and expands the grid when required.
func initGraph(shift, expand bool) {
	automata = getAutomata(mode)

	// If shift is true, create a new grid and shift cells based on expand.
	if shift {
		newSlice := make([][]Cell, int(gridsize))
		for i := range newSlice {
			newSlice[i] = make([]Cell, int(gridsize))
		}
		initCells(newSlice)
		shiftCells(newSlice, expand)
	} else {
		// Create the grid and add the automata configuration.
		grid = make([][]Cell, int(gridsize))
		for i := range grid {
			grid[i] = make([]Cell, int(gridsize))
		}
		initCells(grid)
		addConfiguration(int(gridsize)/2-5, int(gridsize)/2-10, automata)
	}

	// gridCopy is always created.
	gridCopy = make([][]Cell, int(gridsize))
	for i := range grid {
		gridCopy[i] = make([]Cell, int(gridsize))
	}

	updateNumRoutines(30)
}

// addConfiguration adds the automata configuration to the grid starting at the specified row and col.
func addConfiguration(row, col int, automata Automata) {
	startConfig := automata.getConfig()
	firstCell = strings.Split(startConfig, "\n")
	r := row
	for _, line := range firstCell {
		c := col
		for _, char := range line {
			if !unicode.IsSpace(char) {
				if inGrid(row, c) {
					grid[r][c].val = string(char)
				}
			}
			c++
		}
		r++
	}
}

// shiftCells shifts the cells in the new grid based on the expand parameter.
func shiftCells(newSlice [][]Cell, expand bool) {
	if expand {
		for i := 0; i < int(gridsize)-2; i++ {
			for j := 0; j < int(gridsize)-2; j++ {
				newSlice[i+1][j+1] = grid[i][j]
				newSlice[i+1][j+1].row = i + 1
				newSlice[i+1][j+1].col = j + 1
			}
		}
	} else {
		for i := 0; i < int(gridsize)-2; i++ {
			for j := 0; j < int(gridsize)-2; j++ {
				newSlice[i][j] = grid[i+1][j+1]
				newSlice[i][j].row = i
				newSlice[i][j].col = j
			}
		}
	}
	grid = newSlice
}

// updateNumRoutines sets the number of goroutines based on the grid size and the number of chunks.
func updateNumRoutines(chunks int) {
	goroutines = (int(gridsize) + chunks - 1) / chunks
}

// updateGridChunk updates a chunk of the grid between startRow and endRow, inclusive, using the provided automata.
func updateGridChunk(startRow, endRow int, wg *sync.WaitGroup, automata Automata) {
	defer wg.Done()

	for i := startRow; i < endRow; i++ {
		for j, cell := range grid[i] {
			gridCopy[i][j].val = updateCell(cell, automata)
		}
	}
}

// processNextGen updates the grid to the next generation using the specified mode and automata.
func processNextGen(mode string, automata Automata) {
	// If the current generation is the stopping generation, stop processing.
	if generations == stopGeneration {
		running = false
		return
	}

	updateMode(automata)

	// Process the grid in parallel using goroutines.
	for i := 0; i < goroutines; i++ {
		startRow := i * (int(gridsize) / goroutines)
		endRow := startRow + (int(gridsize) / goroutines)

		if i == goroutines-1 {
			endRow = int(gridsize)
		}

		// Update a chunk of the grid using a goroutine.
		wg.Add(1)
		go updateGridChunk(startRow, endRow, &wg, automata)
	}
	wg.Wait()

	// Swap the references of grid and gridCopy
	temp := grid
	grid = gridCopy
	gridCopy = temp
	generations++
}

// rotate rotates the given neighbours array by moving the first element to the end.
func rotate(neighbours [][]int) [][]int {
	first := neighbours[0]
	rotated := append(neighbours[1:], first)
	return rotated
}

// getStates returns an array of 4 strings representing the states of the given cell and its neighbors.
func getStates(cell Cell) []string {
	states := make([]string, 4)
	center := cell.val

	// Iterate over the neighboring cells starting with TLBR
	neighbours := [][]int{{0, -1}, {-1, 0}, {0, 1}, {1, 0}}
	for i := 0; i < 4; i++ {
		currState := center
		for _, offset := range neighbours {
			r, c := cell.row+offset[0], cell.col+offset[1]
			if inGrid(r, c) {
				currState += grid[r][c].val
			} else {
				currState = ""
				break
			}
		}

		states[i] = currState
		neighbours = rotate(neighbours)
	}

	return states
}

func mod(a, b int) int {
	// C version of '%' is a remainder function not modulo
	// so we must implement our own mod function
	r := a % b
	if r < 0 {
		return r + b
	} else {
		return r
	}
}

// updateCell returns the new value for the given cell based on its current state and its neighbors.
func updateCell(cell Cell, automata Automata) string {
	states := getStates(cell)
	checked := false
	for i, state := range states {
		// if we're not on the edge or we have a zero
		// state, leave the state alone
		if state != "" && state != "00000" {
			checked = true
			if val, ok := automata.findRule(state); ok { // rule was found
				return val
			}
		}

		if i == 3 && checked {
			return automata.update(cell)
		}
	}

	return cell.val
}
