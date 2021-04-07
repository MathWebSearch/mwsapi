package xmlutil

import "github.com/beevik/etree"

// InnerXML returns the inner xml source of element
func InnerXML(element *etree.Element) (string, error) {
	doc := etree.NewDocument()
	for _, ct := range element.Child {
		switch c := ct.(type) {
		case *etree.Element:
			doc.AddChild(CopySafe(c))
		default:
			doc.AddChild(c)
		}
	}
	return doc.WriteToString()
}
