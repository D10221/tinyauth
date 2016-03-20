package main

import (
	"flag"
	"fmt"
	"github.com/D10221/tinyauth/credentials"
	"github.com/D10221/tinyauth/criptico"
	"github.com/D10221/tinyauth/config"
	"os"
	"path/filepath"
	"io"
)

type Options struct {
	Config *string
	Task   *string
	Left   *string
	Right  *string
	Word   *string
	Pwd    *string
	Secret *string
	AuthorizationKey *string
}

func (o *Options) NeedsConfigFile () bool {
	if *o.Secret == "" { return true }
	if *o.AuthorizationKey == "" { return true }
	return false
}

var options = &Options{}

func main() {
	e:= ParseCommandLine(os.Args[1:])
	if e != nil {
		fmt.Println(e.Error())
	}
	if options.NeedsConfigFile(){
		e:= LoadConfig()
		if e!= nil {
			fmt.Fprintln( os.Stdout,  )
			os.Exit(-1)
		}
	}
	e= SwitchTask(os.Stdout)
	if e != nil {
		fmt.Println(e.Error())
	}
}

func SwitchTask(writer io.Writer) error {

	if *options.Task == "" {
		return Bonkers("task required")
	}

	switch *options.Task {
	// ...
	case "encode":
		return encode(writer, *options.Left, *options.Right)
	case "decode":
		return decode(writer, *options.Word)
	case "encrypt":
		return encrypt(writer, *options.Word)
	case "decrypt":
		return decrypt(writer, *options.Word)
	default:
		return BonkersF("Unkown task: %s", *options.Task)
	}
	return nil
}

func encode(w io.Writer, left, right string) error {
	if left == "" && right == "" {
		return Bonkers("Nothing to encode: or uname and/or pwd required")
	}
	fmt.Fprintf(w,"%s \n", credentials.New(left, right).Encode())
	return nil
}

func decode(w io.Writer, word string) error {
	if word == "" {
		return Bonkers("Nothing to decode, give me a --word=<Word>")
	}
	fmt.Fprintf(w,"%s \n", credentials.ShouldDecode(word))
	return nil
}
func encrypt(w io.Writer,word string) error {

	if word == "" {
		return Bonkers("Nothing to encrypt, give me a --word=<Word>")
	}
	key := []byte(config.Current.Secret)

	fmt.Fprintf(w,"... \n %s \n ...\n",criptico.Encrypt(key, word))

	return nil
}

func decrypt(w io.Writer,word string) error {

	if word == "" {
		return Bonkers("Nothing to encrypt, give me a --word=<Word>")
	}
	key := []byte(config.Current.Secret)

	fmt.Fprintf(w,"... \n %s \n ... \n", criptico.Decrypt(key, word))

	return nil
}

func ParseCommandLine(args []string) error {

	flag.CommandLine.Parse(args)

	if *options.Secret != "" {
		config.Current.Secret = *options.Secret
	}

	if *options.AuthorizationKey != "" {
		config.Current.AuthorizationKey = *options.AuthorizationKey
	}
	return nil
}

func LoadConfig() error {
	// ... or load config
	dir, _ := os.Getwd()
	path := filepath.Join(dir, *options.Config)
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return &Bonk {"config not found: %s \n", 404 }
	}

	e = config.Current.LoadConfig(path)
	if e!= nil {
		return e
	}

	if e= config.Current.Validate() ; e != nil {
		return &Bonk{ fmt.Sprintf(":( bad config ... \n %s", e.Error()), 500 }
	}

	return nil
}

func init() {
	options.Task = 	flag.String("task", "", "<encode/decode>/encrypt/decrypt")
	options.Left = 	flag.String("left", "", "-task=encode -left=<username>")
	options.Right = flag.String("right", "", "-task=encode -right=<password>")
	options.Word = flag.String("word", "", "-task=decode -word=<word>")
	options.Config = flag.String("config", "config.json", "-config=<path_to_config_file>")
	options.Secret = flag.String("secret", "", "--secret=<16CHARLONGSUPERSECRETKEY>")
	options.AuthorizationKey = flag.String("authKey", "", "--authKey=<AuthorizationKey>")
}






