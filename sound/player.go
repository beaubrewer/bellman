package sound

import (
	"fmt"
	"io"
	"os"
	"time"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

//PlayAt queues a sound file to play
func PlayAt(time time.Duration, file string) {
	Play(file)
}

//Play a sound file
func Play(file string) error {
	//TODO.. check to make sure the file exists

	f, err := os.Open(fmt.Sprintf("chimes/%s", file))
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())
	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(p, d, buf); err != nil {
		return err
	}
	return nil
}
