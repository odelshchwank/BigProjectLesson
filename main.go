package main

import (
	"fmt"
)

func main() {
	// Два вида блокировок -- оптимистичные и пессимистичные в PostgreSQL -- изучить инфу
	fmt.Printf("Введите число: ")
	var a int
	fmt.Scan(&a)

	if a >= 3 {
		fmt.Println("Больше или равно 3")
	} else {
		fmt.Println("Меньше 3!")
	}
}
