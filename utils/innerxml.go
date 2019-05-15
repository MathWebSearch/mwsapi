package utils

import (
	"encoding/xml"
)

// InnerXML is an alias for string, which is read from / written to XML
// by using the innerXML of the target element
type InnerXML string

// UnmarshalXML unmarshals the content of an abitrary element
func (innerxml *InnerXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var inner struct {
		Value string `xml:",innerxml"`
	}
	err := d.DecodeElement(&inner, &start)
	if err != nil {
		return err
	}

	*innerxml = InnerXML(inner.Value)
	return nil
}

// MarshalXML marshals the content of an arbitrary element
func (innerxml *InnerXML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	inner := struct {
		Value string `xml:",innerxml"`
	}{Value: string(*innerxml)}
	return e.EncodeElement(&inner, start)
}
