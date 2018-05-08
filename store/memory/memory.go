package memory

import (
	"github.com/dimaxgl/ssproxy/store"
	"sync"
	"github.com/dimaxgl/ssproxy/password"
	"github.com/pkg/errors"
)

type memoryStory struct {
	pd       password.Decoder
	values   map[string][]byte
	valuesMx sync.Mutex
}

func (s *memoryStory) Verify(user string, password []byte) (bool, error) {
	s.valuesMx.Lock()
	defer s.valuesMx.Unlock()
	if passwordHash, ok := s.values[user]; ok {
		return s.pd.Verify(password, passwordHash)
	}
	return false, nil
}

func (s *memoryStory) Add(user string, password []byte) error {
	if passwordHash, err := s.pd.Encode(password); err != nil {
		return errors.Wrap(err, `failed to encode password`)
	} else {
		s.valuesMx.Lock()
		defer s.valuesMx.Unlock()
		s.values[user] = passwordHash
	}
	return nil
}

func (s *memoryStory) Initialize(passwordDecoder password.Decoder, params map[string]interface{}) (store.Store, error) {
	return &memoryStory{pd: passwordDecoder, values: make(map[string][]byte)}, nil
}
