package main

import "fmt"

var SALT1 string = "+-a^+6"
var SALT2 string = "+-3^+b+-f"

type Gtts struct {
	first_seed  int
	second_seed int
}

func rshift(val int, n int) int {
	if val >= 0 {
		return val >> uint(n)
	}else{
		return (val + 0x100000000) >> uint(n)
	}
}

func work_token(a int, seed string) int{
	runes := []rune(seed)
	for i := 0; i < len(runes); i += 3 {
		r := runes[i+2]
		fmt.Println(i, r, string(r))
	}
	return 0
}

func main() {
	work_token(1, SALT1)
	work_token(1, SALT2)
}
