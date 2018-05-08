package bcrypt

import (
	"github.com/dimaxgl/ssproxy/password"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const ModuleName = `bcrypt`

type decoder struct {
	cost int
}

type decoderOpts struct {
	Cost int
}

func (d decoder) Encode(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, d.cost)
}

func (d decoder) Verify(password []byte, hash []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d decoder) Initialize(params map[string]interface{}) (password.Decoder, error) {
	var opts decoderOpts

	if err := mapstructure.Decode(params, &opts); err != nil {
		return nil, errors.Wrap(err, `failed to parse bcrypt params`)
	}

	if opts.Cost == 0 {
		opts.Cost = bcrypt.DefaultCost
	}

	if opts.Cost < bcrypt.MinCost || opts.Cost > bcrypt.MaxCost {
		return nil, errors.Errorf("invalid bcrypt cost value: must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	return decoder{cost: opts.Cost}, nil
}

func init() {
	password.RegisterDecoder(ModuleName, &decoder{})
}
