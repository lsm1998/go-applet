package main

import (
	"fmt"
	"github.com/kyokomi/emoji/v2"
)

func printEmoji() {
	fmt.Println("Hello World Emoji!")
	// Hello World Emoji!

	_, _ = emoji.Println(":beer: Beer!!!")
	// 🍺  Beer!!!

	pizzaMessage := emoji.Sprint("I like a :pizza: and :sushi:!!")
	fmt.Println(pizzaMessage)
	// I like a 🍕  and 🍣 !!
}
