package vibot

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kardianos/osext"
)

type ViBot struct {
	log  *log.Logger
	fmap FuncMap
}

// config for setting up vibot
type config struct {
	OxfordAPIKey string `json:"oxford_api_key"`
}

// message wrap a message
type message struct {
	Cmd  string
	Args []string
}

// FuncMap is a map of command string to reponse functions.
type FuncMap map[string]ResponseFunc

// ResponseFunc is a handler for bot command
type ResponseFunc func(m *message)

// GetArgString print all argument in one string
func (m message) GetArgString() string {
	argString := ""
	for _, s := range m.Args {
		argString = argString + s + " "
	}
	return strings.TrimSpace(argString)
}

// InitVi will initalise a Vibot
func InitVibot(configJSON []byte, lg *log.Logger) *ViBot {
	// We'll random numbers throughout ViBot
	rand.Seed(time.Now().UTC().UnixNano())

	if lg == nil {
		lg = log.New(os.Stdout, "[vibot] ", 0)
	}

	var cfg config
	err := json.Unmarshal(configJSON, &cfg)
	if err != nil {
		lg.Fatalf("cannot unmarshal config json: %s", err)
	}

	// keyChannel := make(chan string)
	v := &ViBot{log: lg}
	v.fmap = v.getDefaultFuncMap()

	// Get current executing folder
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		lg.Fatalf("cannot retrive persent working directory: %s", err)
	}

	// Ensure temp directory is created.
	// This is used to store media temporarily.
	tmpDirPath := filepath.Join(pwd, tempDir)
	if _, err := os.Stat(tmpDirPath); os.IsNotExist(err) {
		v.log.Printf("[%s] creating temporary directory", time.Now().Format(time.RFC3339))
		mkErr := os.Mkdir(tmpDirPath, 0775)
		if mkErr != nil {
			v.log.Printf("[%s] error creating temporary directory\n%s", time.Now().Format(time.RFC3339), err)
		}
	}

	return v
}

func (v *ViBot) getDefaultFuncMap() FuncMap {
	return FuncMap{
		"/hello": v.SayHello,
	}
}

//func (v *ViBot) Router(msg string) {

//}

// AddFunction add a repsonse to funcmap
func (v *ViBot) AddFunction(command string, resp ResponseFunc) error {
	if !strings.HasPrefix(command, "/") {
		return fmt.Errorf("not a valid command string - it should be of the format / something")
	}
	v.fmap[command] = resp
	return nil
}

// GoSafely will prevents program if some thing cashed or panic
func (v *ViBot) GoSafely(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]

				v.log.Printf("PANIC: %s\n%s", err, stack)
			}
		}()

		fn()
	}()
}
