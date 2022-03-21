package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

func Top10(sss string) []string {
	asm := make(map[string]int)

	if sss == "" {
		return []string{}
	}

	var out []string

	re := regexp.MustCompile(`\s+`)
	sm := re.Split(sss, -1)

	for i := 0; i < len(sm)-1; i++ {
		asm[sm[i]]++
	}

	type ams struct {
		key   string
		value int
	}

	sam := make([]ams, len(asm))

	i := 0
	for k, v := range asm {
		sam[i] = ams{k, v}
		i++
	}

	sort.Slice(sam, func(i, j int) bool {
		return sam[i].value > sam[j].value
	})

	var val, count, start, end int

	for i := 0; i < 10; i++ {
		if val == 0 {
			val = sam[i].value
			start = i
		}

		if val != sam[i].value {
			end = i
			val = sam[i].value

			if count > 1 {
				sort.Slice(out[start:end], func(i, j int) bool {
					return out[i+start] < out[j+start]
				})
			}

			start = i
			count = 0
		}

		count++
		out = append(out, sam[i].key)
	}

	if count > 1 {
		end = 10
		sort.Slice(out[start:end], func(i, j int) bool {
			return out[i+start] < out[j+start]
		})
	}

	return out
}
