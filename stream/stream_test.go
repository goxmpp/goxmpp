package stream_test

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/decoder"
import "bytes"
import "testing"

var x = `<iq to="test@conference.jabber.ru" id="ab7ca" type="set">
    <query xmlns="http://jabber.org/protocol/muc#admin">
        <item affiliation="outcast" jid="test1@example.net"></item>
        <item affiliation="outcast" jid="test2@example.net"></item>
    </query>
</iq>`

type Iq struct {
	XMLName xml.Name `xml:"iq"`
	To      string   `xml:"to,attr,omitempty"`
	Id      string   `xml:"id,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	stream.InnerXML
}

type Query struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	stream.InnerXML
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
}

func TestUnmarshal(t *testing.T) {
	d := decoder.NewInnerDecoder()

	w := &stream.Wrapper{StreamDecoder: xml.NewDecoder(bytes.NewReader([]byte(x))), InnerDecoder: d}
	iqr := stream.NewElementHandlerRegistrator()
	qr := stream.NewElementHandlerRegistrator()

	iqr.Register("http://jabber.org/protocol/muc#admin query", func() stream.Element {
		return &Query{InnerXML: stream.InnerXML{Registrator: qr}}
	})

	qr.Register(" item", func() stream.Element {
		return &Item{}
	})

	stream.HandlerRegistrator.Register(" iq", func() stream.Element {
		return &Iq{InnerXML: stream.InnerXML{Registrator: iqr}}
	})

	s := stream.NextStanza(w)

	var buffer []byte
	var err error
	if buffer, err = xml.MarshalIndent(s, "", "    "); err != nil {
		t.Error(err)
	}

	t.Log(string(buffer))
	t.Log([]byte(x))
	t.Log(string(buffer))
	t.Log(x)
	for index, b := range []byte(x) {
		if buffer[index] != b {
			t.Fatal("Source doesn't match to result in pos", index)
		}
	}
}
