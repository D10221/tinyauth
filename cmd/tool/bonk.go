package main

import "fmt"

type Bonk struct {
	Message string
	Code int
}

// error implementation
func (c *Bonk) Error() string {
	return c.Message
}

func Bonkers(message string) error {
	return &Bonk{Message: message}
}

func BonkersF(format string, i interface{}) error {
	return &Bonk{ Message: fmt.Sprintf(format, i), }
}

func BonkCode(e error) int {
	value, ok := e.(*Bonk)
	if !ok { panic("error is not Bok") }
	return value.Code
}