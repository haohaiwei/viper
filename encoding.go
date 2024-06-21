package viper

import (
	"github.com/spf13/viper/internal/encoding/dotenv"
	"github.com/spf13/viper/internal/encoding/hcl"
	"github.com/spf13/viper/internal/encoding/ini"
	"github.com/spf13/viper/internal/encoding/javaproperties"
	"github.com/spf13/viper/internal/encoding/json"
	"github.com/spf13/viper/internal/encoding/toml"
	"github.com/spf13/viper/internal/encoding/yaml"
)

// Encoder encodes Viper's internal data structures into a byte representation.
// It's primarily used for encoding a map[string]any into a file format.
type Encoder interface {
	Encode(v map[string]any) ([]byte, error)
}

// Decoder decodes the contents of a byte slice into Viper's internal data structures.
// It's primarily used for decoding contents of a file into a map[string]any.
type Decoder interface {
	Decode(b []byte, v map[string]any) error
}

// Codec combines [Encoder] and [Decoder] interfaces.
type Codec interface {
	Encoder
	Decoder
}

// EncoderRegistry returns an [Encoder] for a given format.
// The second return value is false if no [Encoder] is registered for the format.
type EncoderRegistry interface {
	Encoder(format string) (Encoder, bool)
}

// DecoderRegistry returns an [Decoder] for a given format.
// The second return value is false if no [Decoder] is registered for the format.
type DecoderRegistry interface {
	Decoder(format string) (Decoder, bool)
}

// [CodecRegistry] combines [EncoderRegistry] and [DecoderRegistry] interfaces.
type CodecRegistry interface {
	EncoderRegistry
	DecoderRegistry
}

// WithEncoderRegistry sets a custom [EncoderRegistry].
func WithEncoderRegistry(r EncoderRegistry) Option {
	return optionFunc(func(v *Viper) {
		v.encoderRegistry2 = r
	})
}

// WithDecoderRegistry sets a custom [DecoderRegistry].
func WithDecoderRegistry(r DecoderRegistry) Option {
	return optionFunc(func(v *Viper) {
		v.decoderRegistry2 = r
	})
}

// WithCodecRegistry sets a custom [EncoderRegistry] and [DecoderRegistry].
func WithCodecRegistry(r CodecRegistry) Option {
	return optionFunc(func(v *Viper) {
		v.encoderRegistry2 = r
		v.decoderRegistry2 = r
	})
}

type codecRegistry struct {
	v *Viper
}

func (r codecRegistry) Encoder(format string) (Encoder, bool) {
	return r.codec(format)
}

func (r codecRegistry) Decoder(format string) (Decoder, bool) {
	return r.codec(format)
}

func (r codecRegistry) codec(format string) (Codec, bool) {
	switch format {
	case "yaml", "yml":
		return yaml.Codec{}, true

	case "json":
		return json.Codec{}, true

	case "toml":
		return toml.Codec{}, true

	case "hcl", "tfvars":
		return hcl.Codec{}, true

	case "ini":
		return ini.Codec{
			KeyDelimiter: r.v.keyDelim,
			LoadOptions:  r.v.iniLoadOptions,
		}, true

	case "properties", "props", "prop": // Note: This breaks writing a properties file.
		return &javaproperties.Codec{
			KeyDelimiter: v.keyDelim,
		}, true

	case "dotenv", "env":
		return &dotenv.Codec{}, true
	}

	return nil, false
}