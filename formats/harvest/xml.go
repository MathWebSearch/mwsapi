package harvest

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/MathWebSearch/mwsapi/internal/xmlutil"
	"github.com/beevik/etree"
)

const MathWebSearchNamespace = "http://search.mathweb.org/ns"

// UnmarshalXMLFrom reads XML from reader and parses it into this harvest.
// Does not reset the harvest, and any existing Data or Expression nodes will be kept.
//
// This method does not validate the resulting harvest document, see .Validate() for that.
//
// Returns the number of bytes read from reader, and any error that occured during reading or parsing.
// The error will either be an underlying XML reading error, or an error of one of the ErrUnmarshal* types in this package.
func (h *Harvest) UnmarshalFrom(reader io.Reader) (n int64, err error) {
	// Read the document from the reader and parse it into xml
	doc := etree.NewDocument()
	n, err = doc.ReadFrom(reader)
	if err != nil {
		return n, err
	}

	return n, h.unmarshalHarvest(doc.Root())
}

// unmarshalDocument implements UnmarshalFrom.
func (h *Harvest) unmarshalHarvest(harvest *etree.Element) error {
	if harvest == nil {
		return ErrUnmarshalNilElement{Elements: []string{"mws:harvest"}}
	}

	// check that it is a 'mws:harvest' element
	if harvest.NamespaceURI() != MathWebSearchNamespace || harvest.Tag != "harvest" {
		return ErrUnmarshalInvalidRoot{Message: "Expected 'harvest' element in '" + MathWebSearchNamespace + "' namespace"}
	}

	for _, element := range harvest.ChildElements() {
		if element == nil {
			return ErrUnmarshalNilElement{Elements: []string{"mws:data", "mws:expr"}}
		}

		// mws:data element is found
		if element.NamespaceURI() == MathWebSearchNamespace && harvest.Tag == "data" {
			data, err := unmarshalDataElement(element)
			if err != nil {
				return err
			}
			h.Data = append(h.Data, data)
			continue
		}

		// mws:expr element is found
		if element.NamespaceURI() == MathWebSearchNamespace && harvest.Tag == "expr" {
			expr, err := unmarshalExprElement(element)
			if err != nil {
				return err
			}
			h.Expressions = append(h.Expressions, expr)
			continue
		}

		return ErrUnmarshalInvalidElement{Got: element.FullTag(), Expected: []string{"mws:data", "mws:expr"}}
	}

	return nil
}

func unmarshalDataElement(data *etree.Element) (d HarvestData, err error) {
	// Read the '<data>' attribute out of the element!
	var foundID bool
	for _, attr := range data.Attr {
		if attr.Key == "data_id" {
			d.ID = attr.Value
			foundID = true
			break
		}
	}
	if !foundID {
		err = ErrUnmarshalDataMissingDataID
		return
	}

	// read the contained xml!
	d.Content, err = xmlutil.InnerXML(data)
	return d, err
}

func unmarshalExprElement(expr *etree.Element) (e HarvestExpr, err error) {
	// Read the 'data_id' and 'local_id' / 'url' attribute out of the element!
	var foundID bool
	var foundDataID bool
	for _, attr := range expr.Attr {

		// found the 'id' of this element
		if attr.Key == "local_id" || attr.Key == "url" {
			if foundID && attr.Value != e.ID {
				err = ErrUnmarshalExprConflictingID{First: e.ID, Second: attr.Value}
				return
			}
			e.ID = attr.Value
			foundID = true
		}

		// found the 'data_id' of this element
		if !foundDataID && attr.Key == "data_id" {
			e.DataID = attr.Value
			foundDataID = true
		}
	}

	if !foundID {
		err = ErrUnmarshalExprMissingID
		return
	}

	if !foundDataID {
		err = ErrUnmarshalExprMissingDataID
		return
	}

	// read the contained xml, TODO: Normalize the namespace here!
	e.Math, err = xmlutil.InnerXML(expr)
	return e, err
}

// ErrUnmarshalNilElement indicates that an element was unexpectedly nil.
type ErrUnmarshalNilElement struct {
	Elements []string
}

func (e ErrUnmarshalNilElement) Error() string {
	return fmt.Sprintf("Received nil element, expected one of %s", strings.Join(e.Elements, ", "))
}

// ErrUnmarshalInvalidElement indicates that an invalid element was received.
type ErrUnmarshalInvalidElement struct {
	Got      string
	Expected []string
}

func (e ErrUnmarshalInvalidElement) Error() string {
	return fmt.Sprintf("Received %s element, expected one of %s", e.Got, strings.Join(e.Expected, ", "))
}

// ErrUnmarshalInvalidRoot indicates that the root element of the harvest is invalid
type ErrUnmarshalInvalidRoot struct {
	Message string
}

func (e ErrUnmarshalInvalidRoot) Error() string {
	return fmt.Sprintf("Invalid Root Element for Harvest: %s", e.Message)
}

// ErrUnmarshalDataMissingDataID indicates that an 'mws:data' element is missing a 'data_id' attribute
var ErrUnmarshalDataMissingDataID = errors.New("'mws:data' element is missing 'data_id' attribute")

// ErrUnmarshalExprMissingID indicates that an 'mws:expr' element has no ID
var ErrUnmarshalExprMissingID = errors.New("'mws:expr' element is missing 'local_id' and 'url' attribute")

// ErrUnmarshalExprMissingDataID indicates that an 'mws:expr' element is missing a 'data_id' attribute
var ErrUnmarshalExprMissingDataID = errors.New("'mws:expr' element is missing 'data_id' attribute")

// ErrUnmarshalExprConflictingID indicates that an 'mws:expr' element has two conflicting IDs
type ErrUnmarshalExprConflictingID struct {
	First, Second string
}

func (e ErrUnmarshalExprConflictingID) Error() string {
	return fmt.Sprintf("Found conflicting 'mws:expr' IDs: %s and %s", e.First, e.Second)
}
