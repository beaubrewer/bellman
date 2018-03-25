package manager

import (
	"fmt"
	"io"
	"os"
	"sync"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

// Audio holds the object/sound catalog map and
// provides the mutex for access
type Audio struct {
	// AudioCatalog maps sound file names to keys
	// keys represent real world objects
	// i.e. front, side - for doors
	Catalog map[string][]string
	// TODO set defaults that can be configured on
	// initial setup
	mutex sync.Mutex
}

// UpdateAudio replaces the object/sound catalog
func (m *Audio) UpdateAudio(catalog map[string][]string) {
	m.mutex.Lock()
	m.Catalog = catalog
	m.mutex.Unlock()
}

// GetAudio returns the sound files for a provided key
func (m *Audio) GetAudio(key string) []string {
	m.mutex.Lock()
	i, ok := m.Catalog[key]
	if !ok {
		// TODO default value could be provided
		return []string{""}
	}
	m.mutex.Unlock()
	return i
}

// NewAudioManager initializes the Audio struct
func NewAudioManager() *Audio {
	return &Audio{
		Catalog: make(map[string][]string),
	}
}

// Play a file
func (m *Audio) Play(file string) error {
	f, err := os.Open(fmt.Sprintf("audio/%s", file))
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

	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(p, d, buf); err != nil {
		return err
	}
	return nil
}
