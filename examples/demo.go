package demo

import (
	"fmt"
)

func Welcome() {
	fmt.Println("welcome to mo2")
	connectMongoDB()
	//test()

	Find()
	defer disconnectMongoDB()

}
