package stream_test

import "encoding/xml"
import _ "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/decoder"
import "bytes"
import "testing"

var x = `<iq to="test@conference.jabber.ru" type="set" id="ab7ca">
    <query xmlns="http://jabber.org/protocol/muc#admin">
        <item affiliation="outcast" jid="test1@example.net">test</item>
        <item affiliation="outcast" jid="test2@example.net">test1</item>
    </query>
</iq>`

func TestUnmarshal(t *testing.T) {
	d := decoder.NewInnerDecoder()

	w := &stream.Wrapper{StreamDecoder: xml.NewDecoder(bytes.NewReader([]byte(x))), InnerDecoder: d}

	s := stream.NextStanza(w)

	var buffer []byte
	var err error
	if buffer, err = xml.MarshalIndent(s, "", "    "); err != nil {
		t.Error(err)
	}

	t.Log("Result (bytes):", buffer)
	t.Log("Source (bytes):", []byte(x))
	t.Log("Result:", string(buffer))
	t.Log("Source:", x)
	for index, b := range []byte(x) {
		if buffer[index] != b {
			t.Fatal("Source doesn't match to result in pos", index)
		}
	}
}
