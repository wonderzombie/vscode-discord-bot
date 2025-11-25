package bot

func Empty(lines []string) bool {
	ret := false
	ret = ret || len(lines) == 0
	ret = ret || (len(lines) == 1 && lines[0] == "")
	return ret
}
