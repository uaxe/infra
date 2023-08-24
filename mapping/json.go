package mapping

import (
	"io"

	"github.com/uaxe/infra/internal/json"
)

var (
	EnableDecoderUseNumber = false

	EnableDecoderDisallowUnknownFields = false
)

func DecodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
