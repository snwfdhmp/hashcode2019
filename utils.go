package main

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func containsStr(str string, arr []string) bool {
	for i := range arr {
		if arr[i] == str {
			return true
		}
	}
	return false
}
