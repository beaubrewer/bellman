// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package audio

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/beaubrewer/bellman/config"
	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

// Audio holds the object/sound catalog map and
// provides the mutex for access
// type Audio struct {
// 	// AudioCatalog maps sound file names to keys
// 	// keys represent real world objects
// 	// i.e. front, side - for doors
// 	Catalog map[string][]string
// 	// TODO set defaults that can be configured on
// 	// initial setup
// 	mutex sync.Mutex
// }

var Catalog map[string][]string
var catelogMutex sync.RWMutex

// UpdateCatalog replaces the object/sound catalog
func UpdateCatalog(catalog map[string][]string) {
	catelogMutex.Lock()
	Catalog = catalog
	catelogMutex.Unlock()
}

// GetAudio returns the sound files for a provided key
func GetAudio(key string) string {
	catelogMutex.RLock()
	i, ok := Catalog[key]
	catelogMutex.RUnlock()
	if !ok {
		i = []string{config.GetDefaultAudio(key)}
	}
	rand.Seed(time.Now().Unix())
	return i[rand.Intn(len(i))]
}

// Play a file
func Play(file string) {
	go func() {
		f, err := os.Open(fmt.Sprintf("audio/%s.mp3", file))
		if err != nil {
			return
		}
		defer f.Close()

		d, err := mp3.NewDecoder(f)
		if err != nil {
			return
		}
		defer d.Close()

		p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
		if err != nil {
			return
		}
		defer p.Close()

		buf := make([]byte, 1024)
		if _, err := io.CopyBuffer(p, d, buf); err != nil {
			return
		}
	}()
}
