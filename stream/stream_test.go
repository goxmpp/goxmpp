package stream_test

import "encoding/xml"
import (
	_ "github.com/dotdoom/goxmpp"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas"
)

import "bytes"
import "testing"
import "log"

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
</iq>
sdfdf`

var iqExpect = `<iq to="test@conference.jabber.ru" type="set" id="ab7ca">
    <query xmlns="http://jabber.org/protocol/muc#admin">
        <item affiliation="outcast" jid="test1@example.net">test</item>
        <item affiliation="outcast" jid="test2@example.net">test1</item>
    </query>
    <unknown>test</unknown>
</iq>`

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
	st := stream.NewStream(bytes.NewBuffer(source))
	st.ElementFactory = stanzas.Factory
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
	unmarshalTester(t, []byte(iqSource), []byte(iqExpect))
}

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

func TestMessageElementUnmarshal(t *testing.T) {
	unmarshalTester(t, []byte(messageSource), []byte(messageExpect))
}

type TestState struct {
	name string
}

func (self *TestState) Name() string {
	return self.name
}

func (self *TestState) SetName(value string) {
	self.name = value
}

func TestSimple(t *testing.T) {
	var state stream.State
	ts_w := TestState{}
	ts_w.SetName("test")

	if ts_w.Name() != "test" {
		t.Fatal("Test basic setter/getter failure.")
	}

	state.Push(&ts_w)

	var ts_r *TestState

	if err := state.Get(ts_r); err == nil {
		t.Fatal("Should fail.")
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
