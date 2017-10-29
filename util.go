package vibot

import (
	"fmt"
	"time"

	"github.com/processone/mpg123"
)

const tempDir = "temp"

func (vb *ViBot) PlaySound(urlFile string) {
	p, err := mpg123.NewPlayer()
	if err != nil {
		fmt.Println(err)
	}
	p.Play(urlFile)
	time.Sleep(1 * time.Second)
}
