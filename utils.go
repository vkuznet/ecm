package main

import "strings"

// custom split function based on '---\n' delimiter
func pwmSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "---\n"); i >= 0 {
		return i + len("---\n"), data[0:i], nil
		//         return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}
