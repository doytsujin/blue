package blue

import (
	"fmt"
	"testing"
	"time"
)

func TestUnit(t *testing.T) {
	sample :=
		[]struct {
			key    string
			value  interface{}
			expect string
		}{
			{"value", 1, "value=1i"},
			{"value", 1.0, "value=1"},
			{"value", 1.2, "value=1.2"},
			{"value", true, "value=true"},
			{"value", "logged out", `value=logged\ out`},
		}

	for _, v := range sample {
		line := newUnit(v.key, v.value).escape().line()
		if line != v.expect {
			t.Errorf("expected %s got %s", v.expect, line)
		}
	}
}

func TestField(t *testing.T) {
	sample :=
		[]struct {
			key    string
			value  interface{}
			expect string
		}{
			{"value", 1, "value=1i"},
			{"value", 1.0, "value=1"},
			{"value", 1.2, "value=1.2"},
			{"value", true, "value=true"},
			{"value", "logged out", `value="logged out"`},
		}

	for _, v := range sample {
		f := Field{Key: v.key, Value: v.value}
		line := f.Line()
		if line != v.expect {
			t.Errorf("expected %s got %s", v.expect, line)
		}
	}
}

func TestMeasurment(t *testing.T) {
	t0 := Tags{
		{"host", "serverA"},
		{"region", "us-west"},
	}
	t1 := Tags{
		{"host", "server A"},
		{"region", "us west"},
	}
	t3 := Tags{
		{"host", "server01"},
		{"region", "uswest"},
	}
	f0 := Fields{
		{"value", 1},
	}
	f1 := Fields{
		{"value", 1.0},
	}
	f2 := Fields{
		{"value", true},
	}
	f3 := Fields{
		{"value", "logged out"},
	}
	f4 := Fields{
		{"load", 10},
		{"alert", true},
		{"reason", "value above maximum threshold"},
	}
	f5 := Fields{
		{"value", 1.0},
	}
	ts := time.Unix(0, 1434055562000000000)
	sample := []struct {
		measure   string
		tags      Tags
		fields    Fields
		timestamp time.Time
		expect    string
	}{
		{"cpu", nil, nil, time.Time{}, "cpu"},
		{"cpu", t0, nil, time.Time{}, "cpu,host=serverA,region=us-west"},
		{"cpu,01", t0, nil, time.Time{}, `cpu\,01,host=serverA,region=us-west`},
		{"cpu", t1, nil, time.Time{}, `cpu,host=server\ A,region=us\ west`},
		{"cpu", nil, f0, time.Time{}, "cpu value=1i"},
		{"cpu", nil, f1, time.Time{}, "cpu value=1"},
		{"cpu", nil, f2, time.Time{}, "cpu value=true"},
		{"cpu", nil, f3, time.Time{}, `cpu value="logged out"`},
		{"cpu", nil, f4, time.Time{}, `cpu alert=true,load=10i,reason="value above maximum threshold"`},
		{"cpu", t3, f5, ts, "cpu,host=server01,region=uswest value=1 1434055562000000000"},
	}

	for _, v := range sample {
		m :=
			&Measurement{
				Name:      v.measure,
				Tags:      v.tags,
				Fields:    v.fields,
				Timestamp: v.timestamp,
			}
		line := m.line()
		if line != v.expect {
			fmt.Println(line)
			t.Errorf("expected %s got %s", v.expect, line)
		}
	}
}
