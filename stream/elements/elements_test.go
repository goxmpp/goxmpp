package elements_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/goxmpp/goxmpp/stream/elements"
)

type BasicXML struct {
	XMLName xml.Name `xml:"html"`
	*elements.InnerElements
}

func NewBasicXML() *BasicXML {
	return &BasicXML{InnerElements: elements.NewInnerElements(elements.NewFactory())}
}

var basicXMLSource = `
<html>
    <body>
        <div>test</div>
    </body>
    <something>really strange</something>
</html>
`

func TestSimple(t *testing.T) {
	bxml := NewBasicXML()
	if err := xml.Unmarshal([]byte(basicXMLSource), bxml); err != nil {
		t.Fatal(err)
	}

	t.Log(bxml)

	if bxml.XMLName.Local != "html" {
		t.Fatal("Wrong outer XML tag name unmarshaled")
	}
	if len(bxml.Elements()) > 0 {
		t.Fatal("Only RawXML should have been parsed")
	}

	if raw_xml, err := xml.MarshalIndent(bxml, "", "    "); err != nil {
		t.Fatal(err)
	} else if strings.TrimSpace(string(raw_xml)) != strings.TrimSpace(basicXMLSource) {
		t.Log("Got:", strings.TrimSpace(string(raw_xml)))
		t.Log("Expected:", strings.TrimSpace(basicXMLSource))
		t.Fatal("Marshaling back to xml failed")
	}
}
