package iso8583

import (
	"github.com/moov-io/iso8583"
	codec "go.unistack.org/micro/v4/codec"
)

type specKey struct{}

func MessageSpec(spec *iso8583.MessageSpec) codec.Option {
	return codec.SetOption(specKey{}, spec)
}
