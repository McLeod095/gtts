package token

import (
	"strconv"
	"log"
)

const SALT_1 = "+-a^+6"
const SALT_2 = "+-3^+b+-f"

type Gtts struct {
	first_seed  int64
	second_seed int64
}

func (g *Gtts)calculate_token(text string) (int64, int64) {
	a := g.first_seed
	for i := 0; i < len(text); i++ {
		a += int64(text[i])
		a = work_token(a, SALT_1)
	}
	a = work_token(a, SALT_2)
	a = a ^ g.second_seed
	if 0 > a {
		a = (a & 2147483647) + 2147483648
	}
	a %= 1E6
	return a, a ^ g.first_seed
}

func rshift(val int64, n uint64) int64 {
	if val >= 0 {
		return val >> n
	}else {
		return (val + 0x100000000) >> n
	}
}

func work_token(a int64, seed string) int64 {
	for i := 0; i < len(seed)-2; i+=3 {
		var d int64
		if seed[i+2] >= 'a' {
			d = int64( seed[i+2]) - 87
		}else{
			dd, err := strconv.Atoi(string(seed[i+2]))
			if err != nil {
				log.Panicln(err)
			}
			d = int64(dd)
		}
		if seed[i+1] == '+' {
			d = rshift(a, uint64(d))
		}else{
			d = a << uint64(d)
		}
		if seed[i] == '+' {
			a = (a + d) & 4294967295
		}else{
			a = a ^ d
		}
	}
	return a
}

func New() *Gtts {

}