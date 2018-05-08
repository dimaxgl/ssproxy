package sha

import (
	"github.com/dimaxgl/ssproxy/password"
	"hash"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"crypto/sha256"
	"crypto/sha512"
	"bytes"
	"sync"
)

var (
	errInvalidAlgo = errors.New(`invalid hashing algorithm`)
)

const (
	sha256Algo = `sha256`
	sha512Algo = `sha512`
)

type decoder struct {
	h      hash.Hash
	hashMx sync.Mutex
	s      []byte
}

type decoderOpts struct {
	Algo string
	Salt string
}

// TODO make it thread safe
func (d *decoder) Encode(password []byte) ([]byte, error) {
	return d.hash(password), nil
}

// TODO make it thread safe
func (d *decoder) Verify(password []byte, hash []byte) (bool, error) {
	return bytes.Equal(d.hash(password), hash), nil
}

func (d *decoder) hash(data []byte) []byte {
	d.hashMx.Lock()
	defer d.hashMx.Unlock()
	defer d.h.Reset()
	hashData := d.h.Sum(append(data, d.s...))
	return hashData
}

func (d *decoder) Initialize(params map[string]interface{}) (password.Decoder, error) {
	var opts decoderOpts
	if err := mapstructure.Decode(params, &opts); err != nil {
		return nil, errors.Wrap(err, `failed to parse sha decoder params`)
	}

	var h hash.Hash
	switch opts.Algo {
	case sha256Algo:
		h = sha256.New()
	case sha512Algo:
		h = sha512.New()
	default:
		return nil, errInvalidAlgo
	}

	if opts.Salt != `` {
		return &decoder{h: h, s: []byte(opts.Salt)}, nil
	}

	return &decoder{h: h, s: nil}, nil
}
