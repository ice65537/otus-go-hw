package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordStat struct {
	Word string
	Freq int64
}

const WordStatGaps string = " \n\t,.:;!?\"'()[]"

type WordStatMap map[string]int64

type WordStatSlice []WordStat

func (sliceWS WordStatSlice) Less(i, j int) bool {
	if sliceWS[i].Freq == sliceWS[j].Freq {
		return sliceWS[i].Word < sliceWS[j].Word
	}
	return sliceWS[i].Freq > sliceWS[j].Freq
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

func (mapWS *WordStatMap) Build(inText string, inGapSet string, inCaseSensitive bool) {
	var buffer strings.Builder

	mapGapSet := make(map[rune]struct{})
	for _, symbol := range inGapSet {
		mapGapSet[symbol] = struct{}{}
	}

	if !inCaseSensitive {
		inText = strings.ToLower(inText)
	}

	arrText := []rune(inText)
	for i, v := range arrText {
		_, okg := mapGapSet[v]
		if !okg {
			buffer.WriteRune(v)
		}
		if okg || i == len(arrText)-1 {
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
		}
	}
	delete((*mapWS), "-")
}

func TopN(inText string, inGapSet string, inCaseSensitive bool, inN int) WordStatSlice {
	mapWS := make(WordStatMap)
	mapWS.Build(inText, inGapSet, inCaseSensitive)

	sliceWS := mapWS.ToSlice()
	sort.Sort(sliceWS)

	if inN > 0 && len(sliceWS) > inN {
		return sliceWS[:inN]
	}
	return sliceWS
}

func Top10(inText string) []string {
	tmpArr := TopN(inText, WordStatGaps, false, 10)
	outArr := make([]string, len(tmpArr))
	for i := range tmpArr {
		outArr[i] = tmpArr[i].Word
	}
	return outArr
}
