package store

import (
	"sync"
	"github.com/armon/go-socks5"
	"github.com/pkg/errors"
)

var (
	storeRegistry   = make(map[string]Store)
	storeMx         sync.RWMutex
	errUnknownStore = errors.New(`store is unknown (import is forgotten?)`)
)

type Store interface {
	socks5.CredentialStore
	Add(user, password string) error
	Initialize(params map[string]interface{}) (Store, error)
}

func Register(name string, impl Store) {
	storeMx.Lock()
	defer storeMx.Unlock()
	storeRegistry[name] = impl
}

func Open(name string, params map[string]interface{}) (Store, error) {
	storeMx.RLock()
	defer storeMx.RUnlock()
	if s, ok := storeRegistry[name]; ok {
		return s.Initialize(params)
	} else {
		return nil, errUnknownStore
	}
}
