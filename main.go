package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Получение ввода от пользователя
	input := getInput()
	fmt.Println(input)

	// Валидация и парсинг входных данных
	left, operator, right, err := validateInput(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Левый операнд", left)
	fmt.Println("Оператор", operator)
	fmt.Println("Правый оператор", right)

	// Валидация и парсинг операндов
	a, b, isRoman, err := validateOperands(left, right, operator)
	if err != nil {
		panic(err)
	}
	fmt.Println("a :=", a, "b :=", b, "isRoman :=", isRoman)

	// Выполнение операции
	result, err := calculate(a, operator, b)
	if err != nil {
		panic(err)
	}

	// Вывод результата согласно введенным данным
	err = outputResult(result, isRoman)
	if err != nil {
		panic(err)
	}
}

// getInput - получает ввод от пользователя
func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите выражение: ")

	if scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка ввода:", err)
	}
	return ""
}

// validateInput - проверяет формат ввода и возвращает операнды и оператор
func validateInput(input string) (left, operator, right string, err error) {
	parts := strings.Fields(input)

	if len(parts) != 3 {
		return "", "", "", errors.New("неправильный формат ввода: ожидалось 'число оператор число'")
	}

	left = parts[0]
	operator = parts[1]
	right = parts[2]

	if operator != "+" && operator != "-" && operator != "*" && operator != "/" {
		return "", "", "", errors.New("некорректный оператор")
	}
	return left, operator, right, nil
}

// validateOperands - проверяет операнды и возвращает их числовые значения и тип
func validateOperands(left, right string, operator string) (int, int, bool, error) {
	leftRoman := isRoman(strings.ToUpper(left))
	rightRoman := isRoman(strings.ToUpper(right))

	// Проверяем римские
	if leftRoman && rightRoman {
		aRoman, err1 := romanToArabic(left)
		bRoman, err2 := romanToArabic(right)
		if err1 != nil || err2 != nil {
			return 0, 0, false, errors.New("некорректный римский операнд")
		}
		if aRoman < 1 || aRoman > 10 {
			return 0, 0, false, errors.New("левый римский операнд должен быть от I до X включительно")
		}
		if bRoman < 1 || bRoman > 10 {
			return 0, 0, false, errors.New("правый римский операнд должен быть от I до X включительно")
		}
		return aRoman, bRoman, true, nil
	}

	// Проверяем арабские
	if !leftRoman && !rightRoman {
		aArabic, err1 := strconv.Atoi(left)
		bArabic, err2 := strconv.Atoi(right)

		if err1 != nil || err2 != nil {
			return 0, 0, false, errors.New("некорректный арабский операнд")
		}
		if operator == "/" && bArabic == 0 {
			return 0, 0, false, errors.New("деление на ноль запрещено")
		}
		if aArabic < 1 || aArabic > 10 {
			return 0, 0, false, errors.New("левый операнд должен быть от 1 до 10 включительно")
		}
		if bArabic < 1 || bArabic > 10 {
			return 0, 0, false, errors.New("правый операнд должен быть от 1 до 10 включительно")
		}
		return aArabic, bArabic, false, nil
	}
	// Смешанные системы
	return 0, 0, false, errors.New("нельзя смешивать арабские и римские числа")
}

// calculate - выполняет арифметическую операцию
func calculate(a int, operator string, b int) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("ошибка, деление на ноль")
		}
		return a / b, nil
	default:
		return 0, errors.New("неподдерживаемый оператор")
	}
}

func isRoman(s string) bool {
	romanNumerals := "IVXLCDM"
	for _, char := range s {
		if !strings.ContainsRune(romanNumerals, char) {
			return false
		}
	}
	return true
}

// outputResult - выводит результат в соответствующей системе счисления
func outputResult(result int, isRoman bool) error {
	if isRoman {
		roman, err := arabicToRoman(result)
		if err != nil {
			return err
		}
		fmt.Println("Результат:", roman)
	} else {
		fmt.Println("Результат:", result)
	}
	return nil
}

// romanToArabic - преобразует римское число в арабское
func romanToArabic(roman string) (int, error) {
	validRomans := map[string]int{
		"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5,
		"VI": 6, "VII": 7, "VIII": 8, "IX": 9, "X": 10,
	}

	upper := strings.ToUpper(roman)
	val, ok := validRomans[upper]
	if !ok {
		return 0, errors.New("некорректная запись римского числа (от I до X)")
	}
	return val, nil
}

// arabicToRoman - преобразует арабское число в римское
func arabicToRoman(num int) (string, error) {
	if num <= 0 {
		return "", errors.New("римские числа должны быть положительными")
	}

	values := []int{100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	var result strings.Builder
	for i, val := range values {
		for num >= val {
			num -= val
			result.WriteString(symbols[i])
		}
	}

	return result.String(), nil
}
