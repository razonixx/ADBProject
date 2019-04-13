package main

import (
	"fmt"

	sentiment "gopkg.in/vmarkovtsev/BiDiSentiment.v1"
)

func main() {
	session, _ := sentiment.OpenSession()
	defer session.Close()
	result, _ := sentiment.Evaluate([]string{"This is the worst movie I have ever seen, it sucks balls!", "That was a great movie"}, session)
	fmt.Println(result[0])
	fmt.Println(result[1])
}
