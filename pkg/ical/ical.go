package ical

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/whitekid/goxp/log"
)

// SPEC
// https://icalendar.org/
// https://www.ietf.org/rfc/rfc2445.txt
type VCalendar struct {
}

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

// Component: VEVENT, VTODO, VJOURNAL, VFREEBUS?y, VTIMEZONE, x-name, iana-token

func escape(s string) string {
	replacer := strings.NewReplacer(
		"\n", `\n`,
		"\\", `\\`,
		";", `\;`,
		",", `\,`,
	)
	return replacer.Replace(s)
}

func unescape(s string) string {
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\\`, "\\",
		`\;`, `;`,
		`\,`, `,`,
	)
	return replacer.Replace(s)
}

func parseDateTime(value string) (tm DateTime, err error) {
	// Date
	if tm.Time, err = time.Parse("20060102", value); err == nil {
		tm.isDate = true
		return
	}

	loc := time.Local
	withTZID := false
	if strings.HasPrefix(value, "TZID=") {
		withTZID = true
		idx := strings.Index(value, ":")
		if idx == -1 {
			return tm, fmt.Errorf("invalid date format: %s", value)
		}
		tzid := strings.SplitN(value[0:idx], "=", 2)
		if tzid[1] != "" {
			loc, err = time.LoadLocation(strings.Replace(tzid[1], "-", "/", 1))
			if err != nil {
				return tm, fmt.Errorf("invalid timezone: %s", tzid[1])
			}
		}

		value = value[idx+1:]
	}

	// UTC
	if tm.Time, err = timeParse([]string{"20060102T150405Z", "20060102T1504Z"}, value, time.UTC); err == nil {
		return
	}

	// local, with timezone
	if tm.Time, err = timeParse([]string{"20060102T150405", "20060102T1504"}, value, loc); err == nil {
		tm.withTZID = withTZID
		return tm, nil
	}

	return tm, err
}

func timeParse(layouts []string, s string, loc *time.Location) (tm time.Time, err error) {
	for _, layout := range layouts {
		tm, err = time.ParseInLocation(layout, s, loc)
		if err == nil {
			return
		}
	}
	return
}

// DateTime represents DATE-TIME or DATE
// - 4.3.4 Date
// - 4.3.5 DateTime
//
// TODO: TZID
type DateTime struct {
	time.Time

	isDate   bool
	withTZID bool
}

func (dt *DateTime) String() string {
	if dt.isDate {
		return dt.Format("20060102")
	}

	if dt.withTZID {
		tzid := strings.Replace(dt.Location().String(), "/", "-", 1)
		return fmt.Sprintf("TZID=%s:%s", tzid, dt.Format("20060102T150405"))
	}

	if dt.Time.Location() == time.Local {
		return dt.Format("20060102T150405")
	}

	return dt.Format("20060102T150405Z07:00")
}

// Duration 4.3.6 Duration
type Duration struct {
	time.Duration
}

const (
	Day  = time.Hour * 24
	Week = Day * 7
)

// Examples:
//
//	P15DT5H0M20S
//	P7W
func parseDuration(s string) (d Duration, err error) {
	orig := s
	neg := false
	if s != "" && (s[0] == '-' || s[0] == '+') {
		neg = s[0] == '-'
		s = s[1:]
	}

	if len(s) == 0 && s[0] != 'P' {
		return d, fmt.Errorf("invalid duration: P required %s", orig)
	}
	s = s[1:]

	for s != "" {
		if s[0] == 'T' {
			s = s[1:]
			continue
		}

		if !('0' <= s[0] && s[0] <= '9') {
			log.Debugf("digit required: %s, %s, %v", s, orig)
			return d, fmt.Errorf("invalid duration: %s", orig)
		}
		v := 0
		v, s = leadingInt(s)

		if len(s) == 0 {
			log.Debugf("unit required: %s, %s", s, orig)
			return d, fmt.Errorf("invalid duration: %s", orig)
		}

		switch s[0] {
		case 'W':
			d.Duration += time.Duration(v) * Week
		case 'D':
			d.Duration += time.Duration(v) * Day
		case 'H':
			d.Duration += time.Duration(v) * time.Hour
		case 'M':
			d.Duration += time.Duration(v) * time.Minute
		case 'S':
			d.Duration += time.Duration(v) * time.Second
		default:
			log.Debugf("invalid unit: %s, %s", s, orig)
			return d, fmt.Errorf("invalid duration: %s", orig)
		}
		s = s[1:]
	}

	if neg {
		d.Duration = -d.Duration
	}

	return d, nil
}

func leadingInt(s string) (int, string) {
	i := 0
	for ; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			continue
		}
		break
	}

	v, _ := strconv.Atoi(s[0:i])
	return v, s[i:]
}

func (d *Duration) String() string {
	du := d.Duration
	s := new(strings.Builder)

	if du < 0 {
		s.WriteString("-")
		du = -du
	}

	s.WriteString("P")

	if du > Week {
		weeks := du / Week
		du = du % Week

		fmt.Fprintf(s, "%dW", weeks)
	}

	if du > Day {
		days := du / Day
		du = du % Day

		fmt.Fprintf(s, "%dD", days)
	}

	if du > 0 {
		s.WriteRune('T')
	}

	if du > time.Hour {
		hours := du / time.Hour
		du = du % time.Hour

		fmt.Fprintf(s, "%dH", hours)
	}

	if du > time.Minute {
		minutes := du / time.Minute
		du = du % time.Minute

		fmt.Fprintf(s, "%dM", minutes)
	}

	if du > time.Second {
		seconds := du / time.Minute

		fmt.Fprintf(s, "%dS", seconds)
	}

	return s.String()
}

// Time
// 4.3.12 Time
type Time struct {
	time.Time
}

func parseTime(s string) (tm Time, err error) {
	if tm.Time, err = time.Parse("150405Z", s); err == nil {
		return
	}

	if tm.Time, err = time.ParseInLocation("150405", s, time.Local); err == nil {
		return
	}

	return
}
