package project7

import "fmt"

func read() int {
	var num int
	fmt.Print("Введите номер последнего нужного числа Фибоначчи:")
	for {
		_, err := fmt.Scanln(&num)
		if err == nil {
			break
		}
		fmt.Println("Введено некорректное значение, пожалуйста повторите ввод!")
	}
	return num
}

func fibGenerator(num int) {
	var first, second int = 1, 1
	var buf int
	switch num {
	case 1:
		fmt.Println(0)
		return
	case 2:
		fmt.Println(1)
		return
	case 3:
		fmt.Println(1)
		return
	}
	fmt.Println("Число под номером  1  равно:  0")
	fmt.Println("Число под номером  2  равно:  1")
	fmt.Println("Число под номером  3  равно:  1")
	for i := 3; i < num; i++ {
		fmt.Println("Число под номером ", i+1, " равно: ", first+second)
		buf = first
		first += second
		second = buf
	}
}

func Project7() {
	var num int = read()
	fibGenerator(num)
}
