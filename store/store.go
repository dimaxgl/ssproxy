package store

import (
	"sync"
	"github.com/pkg/errors"
	"github.com/dimaxgl/ssproxy/password"
)

var (
	storeRegistry   = make(map[string]Store)
	storeMx         sync.RWMutex
	errUnknownStore = errors.New(`store is unknown (import is forgotten?)`)
)

type Store interface {
	Add(user string, password []byte) error
	Verify(user string, password []byte) (bool, error)
	Initialize(passwordDecoder password.Decoder, params map[string]interface{}) (Store, error)
}

func Register(name string, impl Store) {
	storeMx.Lock()
	defer storeMx.Unlock()
	storeRegistry[name] = impl
}

func Open(name string, pd password.Decoder, params map[string]interface{}) (Store, error) {
	storeMx.RLock()
	defer storeMx.RUnlock()
	if s, ok := storeRegistry[name]; ok {
		return s.Initialize(pd, params)
	} else {
		return nil, errUnknownStore
	}
}
