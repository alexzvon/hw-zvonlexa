package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

var re = regexp.MustCompile(`\s+`)

func Top10(sss string) []string {
	asm := make(map[string]int)

	if sss == "" {
		return []string{}
	}

	var out []string

	sm := re.Split(sss, -1)

	for i := 0; i < len(sm); i++ {
		asm[sm[i]]++
	}

	type ams struct {
		key   string
		value int
	}

	sam := make([]ams, len(asm))

	i := 0
	for k, v := range asm {
		sam[i] = ams{key: k, value: v}
		i++
	}

	sort.Slice(sam, func(i, j int) bool {
		if sam[i].value > sam[j].value {
			return true
		}

		if sam[i].value == sam[j].value {
			return sam[i].key < sam[j].key
		}

		return false
	})

	for i := 0; i < 10; i++ {
		out = append(out, sam[i].key)
	}

	return out
}
