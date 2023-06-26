package main

func Find(arr []string, x string) int {
	for i, n := range arr {
		if n == x {
			return i
		}
	}
	return -1
}

func  isIn(s string ,arr []string) bool {
	return Find(arr, s) != -1
}
	