package password

import (
	"sync"
	"github.com/pkg/errors"
)

var (
	decoderStore      = make(map[string]Decoder)
	decoderMx         sync.RWMutex
	errDecoderUnknown = errors.New(`failed to load decoder (forgotten import?)`)
)

type Decoder interface {
	Encode(password []byte) ([]byte, error)
	Verify(password []byte, hash []byte) (bool, error)
	Initialize(params map[string]interface{}) (Decoder, error)
}

func RegisterDecoder(name string, decoder Decoder) {
	decoderMx.Lock()
	defer decoderMx.Unlock()
	decoderStore[name] = decoder
}

func GetDecoder(name string, params map[string]interface{}) (Decoder, error) {
	decoderMx.RLock()
	defer decoderMx.RUnlock()
	if d, ok := decoderStore[name]; ok {
		return d.Initialize(params)
	} else {
		return nil, errDecoderUnknown
	}
}
