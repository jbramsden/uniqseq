package main

import (
	"fmt"

	"github.com/jbramsden/uniqseq"
)

func main() {
	//Create returns the uniqseq structure with default values
	us := uniqseq.Create()
	//Init resets the counter value in this simple instance
	us.Init()
	//Next get the next sequence of the unique sequence
	UniqueString, Numb := us.Next()
	//This should output ___A, 0
	fmt.Printf("%s, %d", UniqueString, Numb)

}
