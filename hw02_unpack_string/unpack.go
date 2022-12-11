package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")
var numbersMap = map[rune]int{'1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '0': 0}

func Unpack(inStr string) (string, error) {
	var (
		count               int
		isNum               bool
		outStr              strings.Builder
		isMultiplyReady     bool = false
		isBackSlashDetected bool = false
		runeToMulti         rune
	)

	for i, inRune := range inStr {
		fmt.Printf("Rune[%d]=%s\n", i, string(inRune))
		count, isNum = numbersMap[inRune]

		if inRune == '\\' {
			if isBackSlashDetected {
				fmt.Printf("Обнаружен второй бэкслэш\n")
				runeToMulti = '\\'
				isBackSlashDetected = false
				isMultiplyReady = true
			} else {
				fmt.Printf("Обнаружен первый бэкслэш\n")
				isBackSlashDetected = true
				//Выводим предыдущий символ, если такой есть
				if isMultiplyReady {
					outStr.WriteRune(runeToMulti)
				}
				isMultiplyReady = false
			}
		} else if isNum {
			if isBackSlashDetected {
				fmt.Printf("Обнаружена экранированная бэкслэшем цифра\n")
				runeToMulti = inRune
				isBackSlashDetected = false
				isMultiplyReady = true
			} else if !isMultiplyReady {
				fmt.Printf("Ошибка: Обнаружена цифра при отсутствии символа для мультипликации!\n")
				return "", ErrInvalidString
			} else if isMultiplyReady {
				fmt.Printf("Символ %s будет повторен %d раз\n", string(runeToMulti), count)
				for j := 1; j <= count; j++ {
					outStr.WriteRune(runeToMulti)
				}
				isMultiplyReady = false
			}
		} else {
			//Просто символ. выводим предыдущий, если такой есть, запоминаем текущий
			if isMultiplyReady {
				outStr.WriteRune(runeToMulti)
			}
			runeToMulti = inRune
			isMultiplyReady = true
		}
	}
	if isMultiplyReady {
		outStr.WriteRune(runeToMulti)
	}
	return outStr.String(), nil
}

func main() {
	const testStr = "\\\\\\89a\\32\\\\ы4c6"
	unpackedStr, _ := Unpack(testStr)
	fmt.Printf("%s ---> %s\n", testStr, unpackedStr)
}
