package internal

// Automata is an interface that defines the methods required for a cellular automaton.
type Automata interface {
	update(Cell) string
	findRule(string) (string, bool)
	getConfig() string
}

// initRules initializes the rule maps for each automaton.
func initRules() {
	// rotateString rotates the given string by moving the first character to the end.
	rotateString := func(state string) string {
		return state[1:] + string(state[0])
	}

	// Define arrays for the different automata rules and their maps.
	var rules [][]string = [][]string{LangtonRules, BylRules, SDSRRules, EvoRules, SexyRules}
	var maps []map[string]string = []map[string]string{langtonMap, bylMap, SDSRMap, evoMap, sexyMap}

	// Iterate through the rules and maps, populating the maps with the rules.
	for i, ruleArray := range rules {
		for _, rule := range ruleArray {
			first := string(rule[0])
			middle := rule[1:5]
			end := string(rule[5])
			for j := 0; j < 4; j++ {
				maps[i][first+middle] = end
				middle = rotateString(middle)
			}
		}
	}

	// Initialize the automaton objects with their respective rules, rule maps, and starting configurations.
	LangtonObject = Langton{
		rules:   LangtonRules,
		ruleMap: langtonMap,
		start:   startLangton,
	}
	BylObject = Byl{
		rules:   BylRules,
		ruleMap: bylMap,
		start:   startByl,
	}
	SDSRObject = SDSR{
		rules:   SDSRRules,
		ruleMap: SDSRMap,
		start:   startLangton,
	}
	EvoObject = Evo{
		rules:   EvoRules,
		ruleMap: evoMap,
		start:   startEvo,
	}
	SexyObject = Sexy{
		rules:   SexyRules,
		ruleMap: sexyMap,
		start:   startEvo,
	}
}

// updateMode checks if the mode has changed and updates the automata accordingly.
func updateMode(automata Automata) {
	if mode != currentMode {
		currentMode = mode
		automata = getAutomata(mode)
	}
}

// getAutomata returns the corresponding automaton object based on the provided mode.
func getAutomata(mode string) Automata {
	switch mode {
	case "langton":
		return LangtonObject
	case "byl":
		return BylObject
	case "sdsr":
		return SDSRObject
	case "evo":
		return EvoObject
	case "sexy":
		return SexyObject
	default:
		panic("Invalid mode provided")
	}
}
