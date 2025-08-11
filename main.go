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
	// ввод данных
	input := getInput()
	fmt.Println("Вы ввели:", input)
	if input == "" {
		panic("пустой ввод, завершение программы")
	}

	// парсинг и валидация ввода
	left, operator, right, err := validateInput(input)
	if err != nil {
		panic(err)
	}

	// валидация и парсинг операндов
	a, b, isRoman, err := validateOperands(left, right, operator)
	if err != nil {
		panic(err)
	}
	fmt.Println("Operator:", operator)
	fmt.Println("a =", a, "b =", b, "isRoman =", isRoman)

	// выполнение операции
	result, err := calculate(a, operator, b)
	if err != nil {
		panic(err)
	}

	// вывод результата согласно введенным данным
	err = outputResult(result, isRoman)
	if err != nil {
		panic(err)
	}
}

func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите данные:")

	if scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка ввода", err)
	}

	return ""
}

func validateInput(input string) (string, string, string, error) {
	parts := strings.Fields(input)

	if len(parts) != 3 {
		return "", "", "", errors.New("неправильный формат ввода: ожидалось 'число оператор число'")
	}

	leftOperand := parts[0]
	operator := parts[1]
	rightOperand := parts[2]

	if operator != "+" && operator != "-" && operator != "*" && operator != "/" {
		return "", "", "", errors.New("неподдерживаемый оператор")
	}

	return leftOperand, operator, rightOperand, nil
}

// валидация операндов
func validateOperands(left, right string, operator string) (int, int, bool, error) {
	// пробуем как арабские числа
	aArabic, err1 := strconv.Atoi(left)
	bArabic, err2 := strconv.Atoi(right)

	if err1 == nil && err2 == nil {
		if aArabic < 1 || aArabic > 10 {
			return 0, 0, false, errors.New("левый операнд должен быть от 1 до 10")
		}
		if operator == "/" && bArabic == 0 {
			return 0, 0, false, errors.New("деление на ноль запрещено")
		}
		if operator != "/" && (bArabic < 1 || bArabic > 10) {
			return 0, 0, false, errors.New("правый операнд должен быть от 1 до 10")
		}
		return aArabic, bArabic, false, nil
	}

	if !isValidRoman(left) || !isValidRoman(right) {
		return 0, 0, false, errors.New("некорректный синтаксис римских чисел")
	}

	// пробуем как римские числа
	aRoman, errA := romanToArabic(left)
	bRoman, errB := romanToArabic(right)

	if errA == nil && errB == nil {
		// проверка:
		if aRoman > 10 || bRoman > 10 {
			return 0, 0, false, errors.New("римские числа должны быть от I до X")
		}
		return aRoman, bRoman, true, nil
	}

	return 0, 0, false, errors.New("операнды должны быть либо арабскими, либо корректными римскими числами от I до X")
}

// логика калькулятора
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
			return 0, errors.New("invalid input")
		}
		return a / b, nil
	default:
		return 0, errors.New("invalid operator")
	}
}

// логика перевода арабских в римские
func arabicToRoman(num int) (string, error) {
	if num <= 0 {
		return "", errors.New("римские числа не могут быть <= 0")
	}
	values := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	symbols := []string{"X", "IX", "VIII", "VII", "VI", "V", "IV", "III", "II", "I"}

	result := ""
	for i := 0; i < len(values); i++ {
		for num >= values[i] {
			num -= values[i]
			result += symbols[i]
		}
	}
	return result, nil
}

// перевод римских в арабские
func romanToArabic(s string) (int, error) {
	romanNumerals := map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
	}

	result := 0
	prev := 0

	for i := len(s) - 1; i >= 0; i-- {
		value := romanNumerals[s[i]]
		if value < prev {
			result -= value
		} else {
			result += value
		}
		prev = value
	}

	return result, nil
}

// исключение ввода римских чисел неверного формата
func isValidRoman(s string) bool {
	invalidPatterns := []string{
		"IIII", "VV", "XXXX", "LL", "CCCC", "DD", "MMMM",
	}

	for _, pattern := range invalidPatterns {
		if strings.Contains(s, pattern) {
			return false
		}
	}
	return true
}

// результат перевода в римские
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
