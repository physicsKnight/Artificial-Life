package internal

type Sexy struct {
	rules   []string
	ruleMap map[string]string
	start   string
}

func (s Sexy) findRule(rule string) (string, bool) {
	val, ok := s.ruleMap[rule]
	return val, ok
}

func (s Sexy) getConfig() string {
	return s.start
}

func (s Sexy) update(cell Cell) string {
	return sexy(cell)
}

func sexy(cell Cell) string {
	center := cell.val
	endState := ""

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
