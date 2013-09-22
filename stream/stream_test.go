package stream_test

import "encoding/xml"
import _ "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/elements/stanzas"
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
	return &stream.Wrapper{
		StreamDecoder:  xml.NewDecoder(bytes.NewReader([]byte(source))),
		InnerDecoder:   decoder.NewInnerDecoder(),
		ElementFactory: stanzas.Factory,
	}
}

func is(got, expect []byte) bool {
	got = bytes.TrimSpace(got)
	if len(got) != len(expect) {
		return false
	}

	for index, b := range got {
		if expect[index] != b {
			return false
		}
	}

	return true
}

func logEpectations(t *testing.T, got, expect, source []byte) {
	t.Log("Result (bytes):", got)
	t.Log("Expected (bytes):", expect)
	t.Log("Source:", string(source))
	t.Log("Result:", string(got))
	t.Log("Expected:", string(expect))
}

func unmarshalTester(t *testing.T, source, expect []byte) {
	s := getWrapper(source).ReadElement()

	buffer, err := xml.MarshalIndent(s, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	logEpectations(t, buffer, expect, source)

	if !is(buffer, expect) {
		t.Fatal("Source doesn't match the result")
	}
}

func TestIQUnmarshal(t *testing.T) {
	unmarshalTester(t, []byte(iqSource), []byte(iqExpect))
}

var messageSource = `<message>
    <body>hi!</body>
    <html xmlns="http://jabber.org/protocol/xhtml-im">
        <body xmlns="http://www.w3.org/1999/xhtml">
            <p style='font-weight:bold'>hi!</p>
        </body>
    </html>
</message>`

var messageExpect = `<message>
    <body>hi!</body>
    <html xmlns="http://jabber.org/protocol/xhtml-im">
        <body xmlns="http://www.w3.org/1999/xhtml">
            <p style='font-weight:bold'>hi!</p>
        </body>
    </html>
</message>`

func TestMessageUnmarshal(t *testing.T) {
	unmarshalTester(t, []byte(messageSource), []byte(messageExpect))
}
