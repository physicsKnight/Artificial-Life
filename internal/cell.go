package internal

import "fmt"

// Cell represents a single cell within the grid.
type Cell struct {
	row, col   int
	val        string
	neighbours [4][1]*Cell
}

// nextTo checks if the cell has a neighbour with the specified value.
func (cell Cell) nextTo(val string) bool {
	rOffset := []int{-1, 1, 0, 0}
	cOffset := []int{0, 0, -1, 1}
	for i := 0; i < 4; i++ {
		r := cell.row + rOffset[i]
		c := cell.col + cOffset[i]
		if inGrid(r, c) && grid[r][c].val == val {
			return true
		}
	}

	return false
}

// nextToRange checks if the cell has a neighbour with a value in the specified range.
func (cell Cell) nextToRange(left, right int) bool {
	rOffset := []int{-1, 1, 0, 0}
	cOffset := []int{0, 0, -1, 1}
	for i := 0; i < 4; i++ {
		r := cell.row + rOffset[i]
		c := cell.col + cOffset[i]
		n := int(grid[r][c].val[0] - '0')
		if inGrid(r, c) && n >= left && n <= right {
			return true
		}
	}
	return false
}

// checkInTube checks if the cell is in a "tube" structure by counting the number
// of neighbours with specific states and returning true if the count is greater than or equal to 2.
func (cell Cell) checkInTube() bool {
	count := 0
	states := []string{"1", "2", "4", "6", "7"}
	for _, n := range states {
		if cell.nextTo(n) {
			count++
		}
	}
	return count >= 2
}

func (cell Cell) getIntState(state string) (int, error) {
	switch state {
	case "0":
		return 0, nil
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	case "3":
		return 3, nil
	case "4":
		return 4, nil
	case "5":
		return 5, nil
	case "6":
		return 6, nil
	case "7":
		return 7, nil
	case "8":
		return 8, nil
	case "9":
		return 9, nil
	case "A", "a":
		return 10, nil
	case "B", "b":
		return 11, nil
	default:
		return 0, fmt.Errorf("invalid hex digit: %s", state)
	}
}
