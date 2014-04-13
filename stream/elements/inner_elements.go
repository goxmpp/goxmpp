package elements

import "encoding/xml"

type Element interface{}

type InnerXML struct {
	XMLName xml.Name
	XML     string `xml:",innerxml"`
}

type InnerElements struct {
	Inner *elements `xml:",any"`
}

func NewInnerElements(factory Factory) *InnerElements {
	return &InnerElements{Inner: NewElements(factory)}
}

func (ie *InnerElements) AddElement(element Element) {
	ie.Inner.AddElement(element)
}

func (ie *InnerElements) Elements() []Element {
	return ie.Inner.Elements()
}

type elements struct {
	elements []Element
	rawXML   []*InnerXML
	Factory  `xml:"-"`
}

func NewElements(factory Factory) *elements {
	return &elements{
		elements: make([]Element, 0),
		rawXML:   make([]*InnerXML, 0),
		Factory:  factory,
	}
}

func (es *elements) AddElement(e Element) {
	es.elements = append(es.elements, e)
}

func (es *elements) Elements() []Element {
	return es.elements
}

func (es *elements) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	elementObject, err := es.DecodeElement(d, &start)
	if err != nil {
		return err
	}

	if innerXML, ok := elementObject.(*InnerXML); ok {
		es.rawXML = append(es.rawXML, innerXML)
	} else {
		es.AddElement(elementObject)
	}

	return nil
}

func (es *elements) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.Encode(es.elements); err != nil {
		return err
	}

	return e.Encode(es.rawXML)
}
