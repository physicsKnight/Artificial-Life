package internal

type Evo struct {
	rules   []string
	ruleMap map[string]string
	start   string
}

func (e Evo) findRule(rule string) (string, bool) {
	val, ok := e.ruleMap[rule]
	return val, ok
}

func (e Evo) getConfig() string {
	return e.start
}

func (e Evo) update(cell Cell) string {
	return evo(cell)
}

func evo(cell Cell) string {
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

func CodesToNucleotides(code string) []string {
	mapping := map[string]string{
		"071": "G",
		"041": "T",
		"1":   "C",
	}

	var results []string
	for len(code) > 0 {
		for i := 3; i >= 1; i -= 2 {
			if i <= len(code) {
				if val, ok := mapping[code[:i]]; ok {
					results = append(results, val)
					code = code[i:]
					break
				}
			}
		}
	}
	return results
}

func NucleotidesToCodes(nucleotides string) []string {
	codes := make([]string, len(nucleotides))
	for i, nucleotide := range nucleotides {
		switch string(nucleotide) {
		case "G":
			codes[i] = "071"
		case "T":
			codes[i] = "041"
		case "C":
			codes[i] = "1"
		default:
			return nil
		}
	}
	return codes
}
