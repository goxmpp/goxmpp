package elements

import "encoding/xml"

type Element interface{}

type InnerXML struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

type InnerElements struct {
	Elements       []interface{}
	ElementFactory `xml:"-"`
	RawXML         []*InnerXML
}

type InnerElementAppender interface {
	AddElement(Element)
}

func (c *InnerElements) AddElement(e Element) {
	c.Elements = append(c.Elements, e)
}

func (c *InnerElements) HandlerInnerElements(d *xml.Decoder, finalName string) error {
	var err error
	for token, err := d.Token(); err == nil; token, err = d.Token() {
		// TODO: Add logic to handler inner elements with same name as our start element
		switch element := token.(type) {
		case xml.EndElement:
			if element.Name.Local == finalName {
				break
			}
		case xml.StartElement:
			elementObject, err := c.DecodeElemenet(d, &element)
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

func (factory ElementFactory) DecodeElemenet(d *xml.Decoder, element *xml.StartElement) (interface{}, error) {
	elementObject, err := factory.Get(*element)
	if err != nil {
		return nil, err
	}
	if err := d.DecodeElement(elementObject, element); err != nil {
		return nil, err
	}

	return elementObject, nil
}

func (ie *InnerElements) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return ie.HandlerInnerElements(d, start.Name.Local)
}
