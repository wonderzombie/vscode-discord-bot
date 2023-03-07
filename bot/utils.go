package bot

func Empty(lines []string) bool {
	ret := false
	if len(lines) == 0 {
		ret = true
	} else if l := len(lines); l == 1 && lines[0] == "" {
		ret = true
	}
	return ret
}
