package vibot

import "fmt"

// SayHello will print say hello
func (v *ViBot) SayHello(msg *message) {
	fmt.Println("Hello the world!")
}

// Help will give some info of vibot
func (v *ViBot) Help(msg *message) {
	fmt.Println("This is help")
}
