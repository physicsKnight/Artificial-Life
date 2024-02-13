package internal

type Langton struct {
	rules   []string
	ruleMap map[string]string
	start   string
}

func (l Langton) findRule(rule string) (string, bool) {
	val, ok := l.ruleMap[rule]
	return val, ok
}

func (l Langton) getConfig() string {
	return l.start
}

func (l Langton) update(cell Cell) string {
	return cell.val
}
