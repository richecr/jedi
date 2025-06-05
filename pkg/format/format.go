package format

import "fmt"

func Printf(text string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf(text, args...)
	}

	fmt.Println(text)
}
