package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	"vibot"

	"github.com/kardianos/osext"
)

func main() {
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("error getting executable folder: %s", err)
	}

	configJSON, err := ioutil.ReadFile(path.Join(pwd, "vibot_config.json"))
	if err != nil {
		log.Fatalf("error reading config file! %s", err)
	}

	logger := log.New(os.Stdout, "[vibot] ", 0)

	vb := vibot.InitVibot(configJSON, logger)
	vb.AddFunction("/help", vb.Help)

	wordSearch := flag.String("w", "", "Text to search. (Required)")
	flag.Parse()

	if *wordSearch == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	vb.SearchDictionary(*wordSearch)

	flag.Parse()
}
