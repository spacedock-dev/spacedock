// ABOUTME: Agent-facing --json serializer — per-command envelopes of ordered,
// ABOUTME: string-valued objects, compact + newline-terminated, HTML escaping off.
package status

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// jsonValue is one of the shapes a --json document is built from: a leaf value,
// an ordered object, or an array of ordered objects. Existing status records keep
// their string-valued fields unless a newer public schema explicitly opts into a
// native JSON type, and insertion order keeps the bytes stable run-to-run.
type jsonValue interface{ writeJSON(b *bytes.Buffer) }

// jsonStr is a string leaf. score "0.38", team_state.present "true", and most
// existing status fields flow through here to preserve the historical string-valued
// status contract. New self-describing envelopes may opt into native JSON booleans
// or numbers where their public schema requires them.
type jsonStr string

type jsonBool bool

type jsonInt int

// jsonObj is an ordered key->value object. Insertion order IS emission order, so
// the byte output is reproducible (no Go map iteration).
type jsonObj struct {
	keys []string
	vals map[string]jsonValue
}

// jsonArr is an ordered array of objects.
type jsonArr []*jsonObj

// jsonStrArr is an ordered array of strings (the mods-per-point lists in --boot).
type jsonStrArr []string

func newJSONObj() *jsonObj {
	return &jsonObj{vals: map[string]jsonValue{}}
}

// set appends key with a string value, preserving first-set order.
func (o *jsonObj) set(key, val string) *jsonObj {
	o.setValue(key, jsonStr(val))
	return o
}

// setValue appends key with an arbitrary jsonValue (nested object/array).
func (o *jsonObj) setValue(key string, val jsonValue) *jsonObj {
	if _, ok := o.vals[key]; !ok {
		o.keys = append(o.keys, key)
	}
	o.vals[key] = val
	return o
}

func (s jsonStr) writeJSON(b *bytes.Buffer) { writeJSONString(b, string(s)) }

func (v jsonBool) writeJSON(b *bytes.Buffer) {
	if v {
		b.WriteString("true")
		return
	}
	b.WriteString("false")
}

func (v jsonInt) writeJSON(b *bytes.Buffer) { b.WriteString(strconv.Itoa(int(v))) }

func (o *jsonObj) writeJSON(b *bytes.Buffer) {
	b.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		writeJSONString(b, k)
		b.WriteByte(':')
		o.vals[k].writeJSON(b)
	}
	b.WriteByte('}')
}

func (a jsonArr) writeJSON(b *bytes.Buffer) {
	b.WriteByte('[')
	for i, o := range a {
		if i > 0 {
			b.WriteByte(',')
		}
		o.writeJSON(b)
	}
	b.WriteByte(']')
}

func (a jsonStrArr) writeJSON(b *bytes.Buffer) {
	b.WriteByte('[')
	for i, s := range a {
		if i > 0 {
			b.WriteByte(',')
		}
		writeJSONString(b, s)
	}
	b.WriteByte(']')
}

// writeJSONString encodes s as a JSON string with HTML escaping OFF, so `&`,
// `<`, `>` survive as written (matching the table) and the bytes never carry
// &-style sequences a token proxy could mangle.
func writeJSONString(b *bytes.Buffer, s string) {
	var tmp bytes.Buffer
	enc := json.NewEncoder(&tmp)
	enc.SetEscapeHTML(false)
	// Encode the bare string; the encoder appends a newline we trim.
	_ = enc.Encode(s)
	b.WriteString(strings.TrimRight(tmp.String(), "\n"))
}

// emitJSON writes doc as one compact, newline-terminated JSON document.
func emitJSON(w io.Writer, doc jsonValue) {
	var b bytes.Buffer
	doc.writeJSON(&b)
	b.WriteByte('\n')
	w.Write(b.Bytes())
}
