// Package iso8583 provides a iso8583 codec
package iso8583 // import "go.unistack.org/micro-codec-iso8583/v3"

import (
	"fmt"
	"io"

	"github.com/moov-io/iso8583"
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

	if m, ok := v.(*codec.Frame); ok {
		return m.Data, nil
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

	if m, ok := v.(*codec.Frame); ok {
		m.Data = b
		return nil
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

func (c *iso8583Codec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *iso8583Codec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	return c.Unmarshal(buf, v)
}

func (c *iso8583Codec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	return err
}

func (c *iso8583Codec) String() string {
	return "iso8583"
}

func NewCodec(opts ...codec.Option) *iso8583Codec {
	return &iso8583Codec{opts: codec.NewOptions(opts...)}
}
