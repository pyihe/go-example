package testify

import "errors"

func Cal(symbol string, p1, p2 int) (int, error) {
	switch symbol {

	case "+":
		return p1 + p2, nil
	case "-":
		return p1 - p2, nil
	case "*":
		return p1 * p2, nil
	case "/":
		if p2 == 0 {
			return 0, errors.New("divide by zero")
		}
		return p1 / p2, nil
	default:
		return 0, errors.New("invalid symbol")
	}
}
