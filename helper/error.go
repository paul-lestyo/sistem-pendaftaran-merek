package helper

import "fmt"

func PanicIfError(err error) {
	defer HandlePanic()

	if err != nil {
		panic(err)
	}
}

func HandlePanic() {
	r := recover()

	if r != nil {
		fmt.Println("RECOVER", r)
	}
}
