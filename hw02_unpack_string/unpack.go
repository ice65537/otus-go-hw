package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	numbersMap       = map[rune]int{'1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '0': 0}
)

func Unpack(inStr string) (string, error) {
	var (
		count               int
		isNum               bool
		outStr              strings.Builder
		isMultiplyReady     = false
		isBackSlashDetected = false
		runeToMulti         rune
	)
	fmt.Printf("Распаковка строки '%s'\n", inStr)
	for i, inRune := range inStr {
		fmt.Printf("Rune[%d]=%s\n", i, string(inRune))
		count, isNum = numbersMap[inRune]

		switch {
		case inRune == '\\':
			if isBackSlashDetected {
				fmt.Printf("Обнаружен второй бэкслэш\n")
				runeToMulti = '\\'
				isBackSlashDetected = false
				isMultiplyReady = true
			} else {
				fmt.Printf("Обнаружен первый бэкслэш\n")
				isBackSlashDetected = true
				// Выводим предыдущий символ, если такой есть
				if isMultiplyReady {
					outStr.WriteRune(runeToMulti)
				}
				isMultiplyReady = false
			}
		case isNum:
			switch {
			case isBackSlashDetected:
				fmt.Printf("Обнаружена экранированная бэкслэшем цифра\n")
				runeToMulti = inRune
				isBackSlashDetected = false
				isMultiplyReady = true
			case !isMultiplyReady:
				fmt.Printf("Ошибка: Обнаружена цифра при отсутствии символа для мультипликации!\n")
				return "", ErrInvalidString
			case isMultiplyReady:
				fmt.Printf("Символ %s будет повторен %d раз\n", string(runeToMulti), count)
				for j := 1; j <= count; j++ {
					outStr.WriteRune(runeToMulti)
				}
				isMultiplyReady = false
			default:
				fmt.Printf("Ошибка: Кажется случилось что-то, что мы не предусмотрели!\n")
				return "", ErrInvalidString
			}
		default:
			// Просто символ. выводим предыдущий, если такой есть, запоминаем текущий
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
	fmt.Printf("'%s'--->'%s'\n", inStr, outStr.String())
	return outStr.String(), nil
}
