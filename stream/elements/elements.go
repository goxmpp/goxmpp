package elements

import "encoding/xml"

type Element interface{}

type InnerXML struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

type InnerElements struct {
	Elements       []Element
	ElementFactory `xml:"-"`
	RawXML         []*InnerXML
}

func NewInnerElements(factory ElementFactory) *InnerElements {
	return &InnerElements{
		Elements:       make([]Element, 0),
		RawXML:         make([]*InnerXML, 0),
		ElementFactory: factory,
	}
}

func (c *InnerElements) AddElement(e Element) {
	c.Elements = append(c.Elements, e)
}

func (c *InnerElements) HandlerInnerElements(d *xml.Decoder, final xml.EndElement) error {
	var err error
	for token, err := d.Token(); err == nil; token, err = d.Token() {
		// TODO: Add logic to handler inner elements with same name as our start element
		switch element := token.(type) {
		case xml.EndElement:
			if element.Name.Local == final.Name.Local {
				break
			}
		case xml.StartElement:
			elementObject, err := c.DecodeElement(d, &element)
			if err != nil {
				return err
			}

			if innerXML, ok := elementObject.(*InnerXML); ok {
				c.RawXML = append(c.RawXML, innerXML)
			} else {
				c.AddElement(elementObject)
			}
		}
	}

	return err
}

func (ie *InnerElements) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return ie.HandlerInnerElements(d, start.End())
}
