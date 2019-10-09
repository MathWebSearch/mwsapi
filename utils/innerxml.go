package utils

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

// InnerXML is an alias for string, which is read from / written to XML
// by using the innerXML of the target element
type InnerXML string

type innerHTMLStruct struct {
	Value string `xml:",innerxml"`
}

// UnmarshalXML unmarshals the content of an abitrary element
func (innerxml *InnerXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var inner innerHTMLStruct
	if err := d.DecodeElement(&inner, &start); err != nil {
		return errors.Wrap(err, "d.DecodeElement failed")
	}

	// set the value
	*innerxml = InnerXML(inner.Value)
	return nil
}

// MarshalXML marshals the content of an arbitrary element
func (innerxml *InnerXML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	inner := innerHTMLStruct{Value: string(*innerxml)}
	if err := e.EncodeElement(&inner, start); err != nil {
		return errors.Wrap(err, "e.EncodeElement failed")
	}

	return nil
}
