package ical

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mustParseDateTime(s string) DateTime {
	if s == "" {
		return DateTime{}
	}

	tm, err := parseDateTime(s)
	if err != nil {
		panic(err)
	}
	return tm
}

func TestEscapeText(t *testing.T) {
	type args struct {
		s string
	}
	tests := [...]struct {
		name string
		args args
		enc  string
	}{
		{`valid`, args{"\\"}, `\\`},
		{`valid`, args{"\\n"}, `\\n`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := escape(tt.args.s)
			require.Equal(t, tt.enc, enc)
			got := unescape(enc)
			require.Equal(t, tt.args.s, got)
		})
	}
}

func TestDateTime(t *testing.T) {
	tzEastern, _ := time.LoadLocation("US/Eastern")
	_ = tzEastern

	type args struct {
		s string
	}
	tests := [...]struct {
		name    string
		args    args
		want    DateTime
		wantStr string
		wantErr bool
	}{
		{`valid date`, args{"19970714"},
			DateTime{
				Time:   time.Date(1997, 7, 14, 0, 0, 0, 0, time.UTC),
				isDate: true},
			"19970714", false},
		{`valid local date-time`, args{"19980118T230000"},
			DateTime{Time: time.Date(1998, 1, 18, 23, 0, 0, 0, time.Local)},
			"19980118T230000", false},
		{`valid utc date-time`, args{"19980119T070000Z"},
			DateTime{Time: time.Date(1998, 1, 19, 7, 0, 0, 0, time.UTC)},
			"19980119T070000Z", false},
		{`valid tzid date-time`, args{"TZID=US-Eastern:19970714T133000"},
			DateTime{Time: time.Date(1997, 7, 14, 13, 30, 0, 0, tzEastern), withTZID: true},
			"TZID=US-Eastern:19970714T133000", false},
		{`valid date-time`, args{"19970901T1300Z"},
			DateTime{Time: time.Date(1997, 9, 1, 13, 0, 0, 0, time.UTC)},
			"19970901T130000Z", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateTime(tt.args.s)
			require.Truef(t, (err != nil) == tt.wantErr, `parseDateTime() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, got, "want %v but got %v", tt.want.Time.String(), got.Time.String())
			require.Equal(t, tt.wantStr, got.String())
		})
	}
}

func TestDuration(t *testing.T) {
	type args struct {
		s string
	}
	tests := [...]struct {
		name    string
		args    args
		want    time.Duration
		wantStr string
		wantErr bool
	}{
		{`valid`, args{"P15DT5H0M20S"}, 15*Day + 5*time.Hour + 20*time.Second, "P2W1DT5H0S", false},
		{`valid`, args{"-P15DT5H0M20S"}, -(15*Day + 5*time.Hour + 20*time.Second), "-P2W1DT5H0S", false},
		{`valid`, args{"P7W"}, 7 * Week, "P7W", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.args.s)
			require.Truef(t, (err != nil) == tt.wantErr, `parseDuration() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, got.Duration)
			require.Equal(t, tt.wantStr, got.String())
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		s string
	}
	tests := [...]struct {
		name    string
		args    args
		want    Time
		wantErr bool
	}{
		{`valid`, args{"230000Z"}, Time{Time: time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC)}, false},
		{`valid`, args{"230000"}, Time{Time: time.Date(0, 1, 1, 23, 0, 0, 0, time.Local)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.args.s)
			require.Truef(t, (err != nil) == tt.wantErr, `parseTime() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
			require.Equalf(t, tt.want, got, "%s %s", tt.want.String(), got.Time.String())
		})
	}
}
