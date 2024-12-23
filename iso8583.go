// Package iso8583 provides a iso8583 codec
package iso8583 // import "go.unistack.org/micro-codec-iso8583/v3"

import (
	"fmt"

	"github.com/moov-io/iso8583"
	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
)

type iso8583Codec struct {
	opts codec.Options
}

var _ codec.Codec = &iso8583Codec{}

func (c *iso8583Codec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	}

	return nil, nil
}

func (c *iso8583Codec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], b...)
		return nil
	case codec.RawMessage:
		copy(m, b)
	}

	var spec *iso8583.MessageSpec
	if options.Context != nil {
		if v, ok := options.Context.Value(specKey{}).(*iso8583.MessageSpec); ok {
			spec = v
		}
	}
	if spec == nil {
		return fmt.Errorf("missing spec option")
	}

	message := iso8583.NewMessage(spec)
	err := message.Unpack(b)
	if err != nil {
		return err
	}

	return message.Unmarshal(v)
}

func (c *iso8583Codec) String() string {
	return "iso8583"
}

func NewCodec(opts ...codec.Option) *iso8583Codec {
	return &iso8583Codec{opts: codec.NewOptions(opts...)}
}
