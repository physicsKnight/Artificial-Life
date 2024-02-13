package internal

type Byl struct {
	rules   []string
	ruleMap map[string]string
	start   string
}

func (b Byl) findRule(rule string) (string, bool) {
	val, ok := b.ruleMap[rule]
	return val, ok
}

func (b Byl) getConfig() string {
	return b.start
}

func (b Byl) update(cell Cell) string {
	return cell.val
}

// for byl you need to check the states one by one
func byl(cell Cell, rules []string) string {
	endState := cell.val
	n := len(rules)
	states := getStates(cell)
	for _, state := range states {
		if len(state) == 5 && state != "00000" { // if we're not on the edge
			index := findRuleByl(state, n, rules)
			if index < n { // rule was found
				rule := rules[index][:5]
				if rule == state {
					endState = string(rules[index][5]) // new state
					break
				}
			}
		}
	}

	return endState
}

func findRuleByl(state string, n int, rules []string) int {
	index := n

	for i, rule := range rules {
		mismatch := false
		for j := 0; j < 5; j++ {
			// 9 is ignore state
			//fmt.Printf("%s: %s\n", state, rule)
			if rule[j] != state[j] && rule[j] != '9' {
				mismatch = true
				break
			}
		}

		if !mismatch {
			index = i
			break
		}
	}

	return index
}
