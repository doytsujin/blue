package blue

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

//Measurement represent influxdb metric.
type Measurement struct {
	Name      string
	Tags      Tags
	Fields    Fields
	Timestamp time.Time
}

func (m *Measurement) String() string {
	return m.line()
}

func (m *Measurement) line() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString(escape(m.Name))
	if m.Tags != nil {
		_, _ = buf.WriteRune(',')
		_, _ = buf.WriteString(m.Tags.Line())
	}
	if m.Fields != nil {
		_, _ = buf.WriteRune(' ')
		_, _ = buf.WriteString(m.Fields.Line())
	}
	if !m.Timestamp.IsZero() {
		_, _ = buf.WriteRune(' ')
		_, _ = buf.WriteString(fmt.Sprint(m.Timestamp.UnixNano()))
	}
	return buf.String()
}

//Unit is a  key value pair.
type Unit struct {
	Key   string
	Value interface{}
}

func newUnit(key string, value interface{}) *Unit {
	return &Unit{Key: key, Value: value}
}

func (u *Unit) escape() *Unit {
	n := newUnit(u.Key, u.Value)
	n.Key = escape(n.Key)
	if k, ok := n.Value.(string); ok {
		n.Value = escape(k)
	}
	return n
}

func (u Unit) line() string {
	switch u.Value.(type) {
	case int, int64, int32:
		u.Value = fmt.Sprintf("%vi", u.Value)
	}
	return u.Key + "=" + fmt.Sprint(u.Value)
}

// escapes the string to suit the influxdb line protocol. It returns the string
// with space,comma and equal sign escaped by the \ character
func escape(src string) string {
	var buf bytes.Buffer
	for _, v := range src {
		switch v {
		case ' ', ',', '=', '"':
			_, _ = buf.WriteString(`\`)
			_, _ = buf.WriteRune(v)
		default:
			_, _ = buf.WriteRune(v)
		}
	}
	return buf.String()
}

//Field is a Unit that represent influxbd field.
type Field Unit

//Line returns string representation of the underlying Field.
func (f *Field) Line() string {
	f.Key = escape(f.Key)
	switch f.Value.(type) {
	case int, int64, int32:
		f.Value = fmt.Sprintf("%vi", f.Value)
	case string:
		return f.Key + "=" + fmt.Sprintf("\"%s\"", f.Value)
	}
	return f.Key + "=" + fmt.Sprint(f.Value)
}

//Fields is a list of Fields.
type Fields []*Field

func (f Fields) Len() int {
	return len(f)
}
func (f Fields) Less(i, j int) bool {
	return f[i].Key < f[j].Key
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//Line returns the string representation of the underlying fields that complies
//with influxdb line protocol.
func (f Fields) Line() string {
	sort.Sort(f)
	var buf bytes.Buffer
	for _, v := range f {
		if buf.Len() == 0 {
			_, _ = buf.WriteString(v.Line())
			continue
		}
		_, _ = buf.WriteRune(',')
		_, _ = buf.WriteString(v.Line())
	}
	return buf.String()
}

//Tag is a Unit representing influxdb field.
type Tag Unit

//Line returns influxdb line protocl compliant string representation of the
//field.
func (t Tag) Line() string {
	u := Unit(t)
	return u.escape().line()
}

//Tags is a list of Tag
type Tags []*Tag

func (t Tags) Len() int {
	return len(t)
}
func (t Tags) Less(i, j int) bool {
	return t[i].Key < t[j].Key
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

//Line returns influxdb line protocol compliant string representation of the
//underlying tags.
func (t Tags) Line() string {
	sort.Sort(t)
	var buf bytes.Buffer
	for _, v := range t {
		if buf.Len() == 0 {
			_, _ = buf.WriteString(v.Line())
			continue
		}
		_, _ = buf.WriteRune(',')
		_, _ = buf.WriteString(v.Line())
	}
	return buf.String()
}
