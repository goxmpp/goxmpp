package stream_test

import "encoding/xml"
import _ "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/decoder"
import "bytes"
import "testing"

var iqSource = `<iq to="test@conference.jabber.ru" id="ab7ca" type="set">
sdfsdf
    <query xmlns="http://jabber.org/protocol/muc#admin">
    sfsdf
        <item affiliation="outcast" jid="test1@example.net">test</item>
        sdfsdf
        <item affiliation="outcast" jid="test2@example.net">test1</item>
        sfsdf
    </query>
    sdfsdf
</iq>
sdfdf`

var iqExpect = `<iq to="test@conference.jabber.ru" type="set" id="ab7ca">
    <query xmlns="http://jabber.org/protocol/muc#admin">
        <item affiliation="outcast" jid="test1@example.net">test</item>
        <item affiliation="outcast" jid="test2@example.net">test1</item>
    </query>
</iq>`

func getWrapper(source []byte) *stream.Wrapper {
	return &stream.Wrapper{StreamDecoder: xml.NewDecoder(bytes.NewReader([]byte(iqSource))), InnerDecoder: decoder.NewInnerDecoder()}
}

func TestIQUnmarshal(t *testing.T) {
	s := stream.NextStanza(getWrapper([]byte(iqSource)))

	var buffer []byte
	var err error
	if buffer, err = xml.MarshalIndent(s, "", "    "); err != nil {
		t.Error(err)
	}

	t.Log("Result (bytes):", buffer)
	t.Log("Source (bytes):", []byte(iqSource))
	t.Log("Result:", string(buffer))
	t.Log("Source:", iqSource)
	for index, b := range buffer {
		if iqExpect[index] != b {
			t.Fatal("Source doesn't match to result in pos", index)
		}
	}
}
