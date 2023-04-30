package iso8583

import (
	"bytes"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"go.unistack.org/micro/v4/codec"
)

func TestFrame(t *testing.T) {
	s := &codec.Frame{Data: []byte("test")}

	buf, err := NewCodec().Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`test`)) {
		t.Fatalf("bytes not equal %s != %s", buf, `test`)
	}
}

func TestSpec(t *testing.T) {
	buf := []byte{48, 56, 48, 48, 130, 32, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 3, 17, 32, 72, 82, 37, 103, 32, 0, 1}
	data := &message{}
	if err := NewCodec().Unmarshal(buf, data, MessageSpec(newSpec())); err != nil {
		t.Fatal(err)
	}
	if data.F7.Value() != 311204852 {
		t.Fatalf("invalid data %#+v", data.F7.Value())
	}
}

func newSpec() *iso8583.MessageSpec {
	return &iso8583.MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "F0: Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Length:      8,
				Description: "F1: Bitmap",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.Fixed,
			}),
			7: field.NewNumeric(&field.Spec{
				Length:      10,
				Description: "F7: Transmission Date and Time",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
				Pad:         padding.Left('0'),
			}),
			11: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "F11: System trace audit number (STAN)",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
				Pad:         padding.Left('0'),
			}),
			70: field.NewNumeric(&field.Spec{
				Length:      3,
				Description: "F70: Network management Information code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}
}

type message struct {
	F7  *field.Numeric `index:"7"`
	F11 *field.Numeric `index:"11"`
	F70 *field.Numeric `index:"70"`
}
