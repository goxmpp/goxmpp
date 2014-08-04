package stream_test

import (
	"encoding/xml"

	_ "github.com/goxmpp/goxmpp"
	"github.com/goxmpp/goxmpp/stream"
)

import "bytes"
import "testing"
import "log"

type BytesBuffer struct {
	*bytes.Buffer
}

func (b BytesBuffer) Close() error {
	return nil
}

type ReadonlyBytesBuffer struct {
	*BytesBuffer
}

func (rb ReadonlyBytesBuffer) Write(data []byte) (int, error) {
	return len(data), nil
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
	t.Log("Result (bytes)  :", got)
	t.Log("Expected (bytes):", expect)
	t.Log("Source  :", string(source))
	t.Log("Result  :", string(got))
	t.Log("Expected:", string(expect))
}

func unmarshalTester(t *testing.T, source, expect []byte) {
	st := stream.NewStream(BytesBuffer{bytes.NewBuffer(source)})
	s, err := st.ReadElement()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", s)
	buffer, err := xml.MarshalIndent(s, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	logEpectations(t, buffer, expect, source)

	if !is(buffer, expect) {
		t.Fatal("Source doesn't match the result")
	}
}

func TestIQElementUnmarshal(t *testing.T) {
	var iqSource = `<iq to="test@conference.jabber.ru" id="ab7ca" type="set">
sdfsdf
    <query xmlns="http://jabber.org/protocol/muc#admin">
    sfsdf
        <item affiliation="outcast" jid="test1@example.net">test</item>
        sdfsdf
        <item affiliation="outcast" jid="test2@example.net">test1</item>
        sfsdf
    </query>
    <unknown>test</unknown>
    sdfsdf
    <unknown2>test</unknown2>
</iq>
sdfdf`

	var iqExpect = `<iq to="test@conference.jabber.ru" type="set" id="ab7ca">
    <query xmlns="http://jabber.org/protocol/muc#admin">
        <item affiliation="outcast" jid="test1@example.net">test</item>
        <item affiliation="outcast" jid="test2@example.net">test1</item>
    </query>
    <unknown>test</unknown>
    <unknown2>test</unknown2>
</iq>`

	unmarshalTester(t, []byte(iqSource), []byte(iqExpect))
}

func TestMessageElementUnmarshal(t *testing.T) {
	var messageSource = `<message>
    <body>hi!<some inner="xml">test</some></body>
    <html xmlns="http://jabber.org/protocol/xhtml-im">
        <body xmlns="http://www.w3.org/1999/xhtml">
            <p style='font-weight:bold'>hi!</p>
        </body>
        <some-unknown-xml><with inner="xml">and with data</with></some-unknown-xml>
    </html>
</message>`

	var messageExpect = `<message>
    <body>hi!<some inner="xml">test</some></body>
    <html xmlns="http://jabber.org/protocol/xhtml-im">
        <body xmlns="http://www.w3.org/1999/xhtml">
            <p style='font-weight:bold'>hi!</p>
        </body>
        <some-unknown-xml><with inner="xml">and with data</with></some-unknown-xml>
    </html>
</message>`

	unmarshalTester(t, []byte(messageSource), []byte(messageExpect))
}

type NameState struct {
	name string
}

func (self *NameState) Name() string {
	return self.name
}

func (self *NameState) SetName(value string) {
	self.name = value
}

func TestState(t *testing.T) {
	var state stream.State
	ts_w := NameState{}
	ts_w.SetName("test")

	if ts_w.Name() != "test" {
		t.Fatal("Test basic setter/getter failure.")
	}

	var ts_r *NameState

	if err := state.Get(&ts_r); err == nil {
		t.Fatal("Should not get state which is not saved.")
	}

	state.Push(&ts_w)

	if err := state.Get(ts_r); err == nil {
		t.Fatal("Should not get state without a pointer.")
	}

	if err := state.Get(&ts_r); err != nil {
		t.Fatal(err)
	}

	if ts_r.Name() != "test" {
		t.Fatal("Test state getter failure.")
	}
	ts_r.SetName("test2")
	if ts_r.Name() != "test2" {
		t.Fatal("Test state setter failure.")
	}
}

func TestAttributes(t *testing.T) {
	var streamOpen = `
<?xml version="1.0"?>
<stream:stream
	xmlns:stream="http://etherx.jabber.org/streams"
	version="1.0"
	xmlns="jabber:client"
	to="localhost"
	xml:lang="en"
	xmlns:xml="http://www.w3.org/XML/1998/namespace">
`
	var st = stream.NewStream(ReadonlyBytesBuffer{&BytesBuffer{bytes.NewBuffer([]byte(streamOpen))}})

	st.SetTo("test")
	if st.To() != "test" {
		t.Fatal("Cannot set 'to' field")
	}

	if err := st.Open(); err != nil {
		t.Fatal("Cannot open stream")
	}
	if st.From() != "localhost" {
		t.Fatal("Unknown 'to' attribute value")
	}
	if st.Version() != "1.0" {
		t.Fatal("Unknown 'version' attribute value")
	}
}
