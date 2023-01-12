package ical

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/fx"
	"gopkg.in/yaml.v3"
)

type Examples struct {
	Events map[string]struct {
		Event struct {
			UID        string `yaml:"UID"`
			DTSTAMP    string `yaml:"DTSTAMP"`
			DTSTART    string `yaml:"DTSTART"`
			DTEND      string `yaml:"DTEND"`
			SUMMARY    string `yaml:"SUMMARY"`
			COMMENT    string `yaml:"COMMENT"`
			CLASS      string `yaml:"CLASS"`
			CATEGORIES string `yaml:"CATEGORIES"`
			RRULE      string `yaml:"RRULE"`
			TRANSP     string `yaml:"TRANSP"`
		} `yaml:"event"`
		Data string
	} `yaml:"events"`
}

func readExamples(t *testing.T) *Examples {
	f, err := os.Open("fixtures/ical_data.yaml")
	require.NoError(t, err)
	defer f.Close()

	examples := new(Examples)
	require.NoError(t, yaml.NewDecoder(f).Decode(examples))

	return examples
}

func TestEventEncode(t *testing.T) {
	examples := readExamples(t)

	type args struct {
		name string
	}

	tests := [...]struct {
		name string
		args args
	}{
		{"example1", args{"example1"}},
		{"example2", args{"example2"}},
		{"example3", args{"example3"}},
		{"example4", args{"example4"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			example := examples.Events[tt.args.name]
			want := &VEvent{
				UID:        example.Event.UID,
				DtStamp:    mustParseDateTime(example.Event.DTSTAMP),
				DtStart:    mustParseDateTime(example.Event.DTSTART),
				DtEnd:      mustParseDateTime(example.Event.DTEND),
				Summary:    example.Event.SUMMARY,
				Comment:    example.Event.COMMENT,
				Class:      example.Event.CLASS,
				Categories: fx.Ternary(example.Event.CATEGORIES != "", strings.Split(example.Event.CATEGORIES, ","), nil),
				RRule:      example.Event.RRULE,
				Transp:     example.Event.TRANSP,
			}

			buf := new(bytes.Buffer)
			err := NewEventEncoder(buf).Encode(want)
			require.NoError(t, err)

			s := bufio.NewScanner(bytes.NewReader(buf.Bytes()))
			for s.Scan() {
				require.LessOrEqual(t, len(s.Text()), 80, "TEXT: %s", s.Text())
			}

			got := new(VEvent)
			err = NewEventDecoder(bytes.NewReader(buf.Bytes())).Decode(got)
			require.NoError(t, err)

			require.Equal(t, want, got)
		})
	}
}

func FuzzEventEncode(f *testing.F) {
	f.Add("summary", "description", time.Now().Unix())
	f.Fuzz(func(t *testing.T, summary, description string, dtStamp int64) {
		evt := &VEvent{
			Summary:     summary,
			Description: description,
			DtStart:     DateTime{Time: time.Unix(dtStamp, 0)},
		}
		buf := new(bytes.Buffer)
		require.NoError(t, NewEventEncoder(buf).Encode(evt))

		got := new(VEvent)
		require.NoError(t, NewEventDecoder(buf).Decode(got))
		require.Equal(t, evt, got)
	})
}
