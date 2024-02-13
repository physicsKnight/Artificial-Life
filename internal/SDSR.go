package internal

type SDSR struct {
	rules   []string
	ruleMap map[string]string
	start   string
}

func (s SDSR) findRule(rule string) (string, bool) {
	val, ok := s.ruleMap[rule]
	return val, ok
}

func (s SDSR) update(cell Cell) string {
	return sdsr(cell)
}

func (s SDSR) getConfig() string {
	return s.start
}

func sdsr(cell Cell) string {
	endState := ""
	center := cell.val
	inTube := cell.checkInTube()

	switch center {
	case "0":
		if inTube && cell.nextTo("1") {
			endState = "1"
		} else {
			endState = "0"
		}
	case "1":
		if inTube {
			if cell.nextTo("7") {
				endState = "7"
			}
			if cell.nextTo("6") {
				endState = "6"
			}
			if cell.nextTo("4") {
				endState = "4"
			}
		}
	case "2":
		if cell.nextTo("3") {
			endState = "1"
		} else if cell.nextTo("2") {
			endState = "2"
		}
	case "4", "6", "7":
		if inTube && cell.nextTo("0") {
			endState = "0"
		}
	}

	if center == "8" {
		endState = "0"
	}
	if cell.nextTo("8") {
		switch center {
		case "0", "1":
			if cell.nextToRange(2, 7) {
				endState = "8"
			} else {
				endState = center
			}
		case "2", "3", "5":
			endState = "0"
		case "4", "6", "7":
			endState = "1"
		}
	}

	if endState == "" {
		switch center {
		case "0":
			endState = "0"
		case "1", "2", "3", "4", "5", "6", "7":
			endState = "8"
		}
	}

	return endState
}
