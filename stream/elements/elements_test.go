package elements_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/dotdoom/goxmpp/stream/elements"
)

type BasicXML struct {
	XMLName xml.Name `xml:"html"`
	*elements.InnerElements
}

func NewBasicXML() *BasicXML {
	return &BasicXML{InnerElements: elements.NewInnerElements(elements.NewFactory())}
}

func (bxml *BasicXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	bxml.XMLName = start.Name
	return bxml.UnmarshalInnerElements(d, start.End())
}

var basicXMLSource = `
<html>
    <body>
        <div>test</div>
    </body>
</html>
`

func TestSimple(t *testing.T) {
	bxml := NewBasicXML()
	if err := xml.Unmarshal([]byte(basicXMLSource), bxml); err != nil {
		t.Fatal(err)
	}

	t.Log(bxml)
	t.Log(bxml.InnerElements)
	t.Log(bxml.InnerElements.RawXML[0])

	if bxml.XMLName.Local != "html" {
		t.Fatal("Wrong outer XML tag name unmarshaled")
	}

	if bxml.RawXML[0].XMLName.Local != "body" {
		t.Fatal("Wrong xml InnerXML outer tag")
	}

	if strings.TrimSpace(bxml.RawXML[0].XML) != `<div>test</div>` {
		t.Fatal("Wrong InnerXML parsed")
	}

	if raw_xml, err := xml.MarshalIndent(bxml, "", "    "); err != nil {
		t.Fatal(err)
	} else if strings.TrimSpace(string(raw_xml)) != strings.TrimSpace(basicXMLSource) {
		t.Log("Got:", strings.TrimSpace(string(raw_xml)))
		t.Log("Expected:", strings.TrimSpace(basicXMLSource))
		t.Fatal("Marshaling back to xml failed")
	}
}
