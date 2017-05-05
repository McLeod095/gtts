package token

import (
	"strconv"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"fmt"
)

var SALT_1 = "+-a^+6"
var SALT_2 = "+-3^+b+-f"

var GttsUrl = "https://translate.google.com/"

type Gtts struct {
	first_seed  int64
	second_seed int64
}

//func (g *Gtts)calculate_token(text string) (int64, int64) {
func (g *Gtts) Calculate_token(text string) string {
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
	//return a, a ^ g.first_seed
	return strings.Join([]string{strconv.FormatInt(a,10), strconv.FormatInt(a ^ g.first_seed, 10)}, ".")
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

func New() Gtts {
	g := Gtts{}
	//if err := g.get_token(); err != nil {
	//	return nil
	//}
	return g
}

func (g *Gtts) get_token() error {
	g.first_seed = int64(time.Now().Unix() / 3600)

	resp, err := http.Get(GttsUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodystr := string(body)
	a_start := strings.Index(bodystr, "a\\x3d")
	if a_start < 0 {
		return fmt.Errorf("Cannot find a\\x3d")
	}
	b_start := strings.Index(bodystr, "b\\x3d")
	if b_start < 0 {
		return fmt.Errorf("Cannot find b\\x3d")
	}
	a_stop := strings.Index(bodystr[a_start:], ";")
	b_stop := strings.Index(bodystr[b_start:], ";")
	a_str := bodystr[a_start+5:a_start+a_stop]
	b_str := bodystr[b_start+5:b_start+b_stop]
	a, err := strconv.Atoi(a_str)
	if err != nil {
		return err
	}
	b, err := strconv.Atoi(b_str)
	if err != nil {
		return err
	}
	g.second_seed = int64(a+b)
	return nil
}