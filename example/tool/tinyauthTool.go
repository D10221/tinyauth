package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"github.com/D10221/tinyauth"
	"github.com/D10221/tinyauth/config"
)

type Options struct {
	Config           *string
	Task             *string
	Left             *string
	Right            *string
	Word             *string
	Pwd              *string
	Secret           *string
	AuthorizationKey *string
}

func (o *Options) NeedsConfigFile() bool {
	if *o.Secret == "" {
		return true
	}
	if *o.AuthorizationKey == "" {
		return true
	}
	return false
}

var options = &Options{}

type TinyApp struct {
	Auth *tinyauth.TinyAuth
}


func main() {

	var app = &TinyApp{Auth : tinyauth.NewTinyAuth(&config.TinyAuthConfig{
		Secret: "", // ABRACADABRA12345
		AuthorizationKey: "Authorization",
		BasicScheme: "Basic ",
	})}

	e := app.ParseCommandLine(os.Args[1:])
	if e != nil {
		fmt.Println(e.Error())
	}

	if options.NeedsConfigFile() {
		e := app.LoadConfig()
		if e != nil {
			fmt.Fprintln(os.Stdout, )
			os.Exit(-1)
		}
	}
	e = app.SwitchTask(os.Stdout)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func (app *TinyApp)SwitchTask(writer io.Writer) error {

	if *options.Task == "" {
		return Bonkers("task required")
	}

	switch *options.Task {
	// ...
	case "encode":
		return app.encode(writer, *options.Left, *options.Right)
	case "decode":
		return app.decode(writer, *options.Word)
	case "encrypt":
		return app.encrypt(writer, *options.Word)
	case "decrypt":
		return app.decrypt(writer, *options.Word)
	default:
		return BonkersF("Unkown task: %s", *options.Task)
	}
	return nil
}

func (app *TinyApp) encode(w io.Writer, left, right string) error {
	if left == "" && right == "" {
		return Bonkers("Nothing to encode: or uname and/or pwd required")
	}
	fmt.Fprintf(w, "%s \n", app.Auth.Encoder.Encode(left, right))
	return nil
}

func (app *TinyApp) decode(w io.Writer, word string) error {
	if word == "" {
		return Bonkers("Nothing to decode, give me a --word=<Word>")
	}
	s, e := app.Auth.Encoder.Decode(word)
	if e != nil {
		return e
	} else {
		fmt.Fprintf(w, "%s \n", s)
	}

	return nil
}
func (app *TinyApp) encrypt(w io.Writer, word string) error {

	if word == "" {
		return Bonkers("Nothing to encrypt, give me a --word=<Word>")
	}

	value, e := app.Auth.Criptico.Encrypt(word)
	if e != nil {
		panic(e)
	}
	fmt.Fprintf(w, "... \n %s \n ...\n", value)

	return nil
}

func (app TinyApp) decrypt(w io.Writer, word string) error {

	if word == "" {
		return Bonkers("Nothing to encrypt, give me a --word=<Word>")
	}
	decrypted, err := app.Auth.Criptico.Decrypt(word)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "... \n %s \n ... \n", decrypted)

	return nil
}

func (app *TinyApp)ParseCommandLine(args []string) error {

	flag.CommandLine.Parse(args)

	if *options.Secret != "" {
		app.Auth.Config.Secret = *options.Secret
	}

	if *options.AuthorizationKey != "" {
		app.Auth.Config.AuthorizationKey = *options.AuthorizationKey
	}
	return nil
}

func (app *TinyApp) LoadConfig() error {
	// ... or load config
	dir, _ := os.Getwd()
	path := filepath.Join(dir, *options.Config)
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return &Bonk{"config not found: %s \n", 404 }
	}

	e = app.Auth.Config.LoadConfig(path)
	if e != nil {
		return e
	}

	if e = app.Auth.Config.Validate(); e != nil {
		return &Bonk{fmt.Sprintf(":( bad config ... \n %s", e.Error()), 500 }
	}

	return nil
}

func init() {

	options.Task = flag.String("task", "", "<encode/decode>/encrypt/decrypt")
	options.Left = flag.String("left", "", "-task=encode -left=<username>")
	options.Right = flag.String("right", "", "-task=encode -right=<password>")
	options.Word = flag.String("word", "", "-task=decode -word=<word>")
	options.Config = flag.String("config", "config.json", "-config=<path_to_config_file>")
	options.Secret = flag.String("secret", "", "--secret=<16CHARLONGSUPERSECRETKEY>")
	options.AuthorizationKey = flag.String("authKey", "", "--authKey=<AuthorizationKey>")
}






