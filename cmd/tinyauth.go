package main

import (
	"flag"
	"fmt"
	"github.com/D10221/tinyauth/credentials"
)

func main()  {
	task:= flag.String("task", "", "<encode/decode>")
	uname:= flag.String("uname", "", "-task=encode -uname=<username>")
	pwd:= flag.String("pwd", "", "-task=encode -pwd=<password>")
	word := flag.String("word", "", "-task=decode -word=<word>")
	flag.Parse()
	if *task == "" {
		fmt.Printf("task required")
		return
	}
	if *task == "encode" && *uname == "" && *pwd == "" {
		fmt.Println("or uname or pwd required")
		return
	}
	if *task =="decode" && *word == "" {
		fmt.Println("nothing to decode")
		return
	}
	if *task == "encode"{
		fmt.Printf("%s \n", credentials.New(*uname,*pwd).Encode())
	} else if *task =="decode" {
		fmt.Printf("%s \n", credentials.ShouldDecode(*word))
	} else {
		fmt.Printf("Unkown task: %s", *task)
	}

}



