package main //hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(_ string) []string {
	// Place your code here.

	return nil
}

type RuneSet map[rune]struct{}

type WordStat struct {
	Word string
	Freq int64
}

type WordStatMap map[string]int64

type WordStatSlice []WordStat

func (sliceWS WordStatSlice) Less(i, j int) bool {
	if sliceWS[i].Freq == sliceWS[j].Freq {
		return sliceWS[i].Word < sliceWS[j].Word
	} else {
		return sliceWS[i].Freq < sliceWS[j].Freq
	}
}

func (sliceWS WordStatSlice) Len() int {
	return len(sliceWS)
}

func (sliceWS WordStatSlice) Swap(i, j int) {
	sliceWS[i], sliceWS[j] = sliceWS[j], sliceWS[i]
}

func (mapWS *WordStatMap) ToSlice() WordStatSlice {
	var outSlice = make([]WordStat, len(*mapWS))
	i := 0
	for idx, value := range *mapWS {
		outSlice[i] = WordStat{Word: idx, Freq: value}
		i++
	}
	return outSlice
}

func (mapWS *WordStatMap) Build(inText string, inGapSet RuneSet, inCaseSensitive bool) {
	var buffer strings.Builder

	if *mapWS == nil {
		newmapWS := make(WordStatMap)
		mapWS = &newmapWS
	}

	if !inCaseSensitive {
		inText = strings.ToLower(inText)
	}
	for _, v := range inText {
		_, okg := inGapSet[v]
		if okg {
			if buffer.Len() > 0 {
				word := buffer.String()
				_, okw := (*mapWS)[word]
				if okw {
					(*mapWS)[word]++
				} else {
					(*mapWS)[word] = 1
				}
				buffer.Reset()
			}
			break
		}
		buffer.WriteRune(v)
	}
}

func TopN(inText string, inGapSet RuneSet, inCaseSensitive bool, inN int) WordStatSlice {
	var mapWS WordStatMap
	mapWS.Build(inText, inGapSet, inCaseSensitive)
	sliceWS := mapWS.ToSlice()
	sort.Sort(sliceWS)
	return nil
}

func main() {
	text := `One, one, one, one, one, and one 3 3 3`
	fmt.Println(Top10(text))
}
