package ical

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/whitekid/goxp"
	"github.com/whitekid/iter"
)

// VEvent rfc2455 4.6.1 Event Component
type VEvent struct {
	Class        string   `validate:"max=100"` // classification, PUBLIC(*), PRIVATE, CONFIDENTIAL
	Created      DateTime // 4.8.7.1 Date/Time Created, 서버 리소스 생성 시점
	Description  string   `validate:"max=500"`
	DtStart      DateTime // or Date
	Geo          string   `validate:"max=100"`
	LastModied   DateTime
	Location     string   `validate:"max=100"`
	Organizer    string   `validate:"max=100"`
	Priority     string   `validate:"max=100"`
	DtStamp      DateTime `validate:"required"` // 4.8.7.2 Date/Time Stamp; as UTC,
	Seq          int      // 4.8.7.4 Sequence Number, revision number
	Status       string   `validate:"max=100"`
	Summary      string   `validate:"max=100"`
	Transp       string   `validate:"max=100"`
	UID          string   `validate:"max=100"`
	URL          string   `validate:"max=100"`
	RecurrenceID string   `validate:"max=100"`
	DtEnd        DateTime // or Date
	Duration     Duration
	Attatch      string   `validate:"max=100"`
	Attendee     string   `validate:"max=100"`
	Categories   []string `validate:"max=100"`
	Comment      string   `validate:"max=100"`
	Contact      string   `validate:"max=100"`
	ExDate       DateTime // or DateTime
	ExRule       string   `validate:"max=100"`
	RStatus      string   `validate:"max=100"`
	Related      string   `validate:"max=100"`
	Resources    string   `validate:"max=100"`
	RDate        DateTime // or DateTime
	RRule        string   `validate:"max=100"` // RecurrenceRule
	XProp        string   `validate:"max=100"` // TODO
}

type eventDecoder struct {
	r io.Reader
}

var _ Decoder = (*eventDecoder)(nil)

func NewEventDecoder(r io.Reader) Decoder {
	return &eventDecoder{r: r}
}

func (d *eventDecoder) Decode(v any) (err error) {
	e, ok := v.(*VEvent)
	if !ok {
		return fmt.Errorf("v must be pointer of VEvent")
	}

	feed := func(s string) error {
		p := strings.SplitN(s, ":", 2)
		if len(p) < 2 {
			panic(s)
		}
		key, value := p[0], p[1]

		switch key {
		case "BEGIN":
		case "END":
			return nil
		case "CLASS":
			e.Class = unescape(value)
		case "CREATED":
			e.Created, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: CREATED")
			}
		case "DESCRIPTION":
			e.Description = unescape(value)
		case "DTSTART":
			e.DtStart, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: DTSTART")
			}
		case "GEO":
			e.Geo = unescape(value)
		case "LAST-MODIFIED":
			e.LastModied, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: LAST-MODIFIED")
			}
		case "LOCATION":
			e.Location = unescape(value)
		case "ORGANIZER":
			e.Organizer = unescape(value)
		case "PRIORITY":
			e.Priority = unescape(value)
		case "DTSTAMP":
			e.DtStamp, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: DISTAMP")
			}
		case "SEQ":
			e.Seq, _ = strconv.Atoi(value)
		case "STATUS":
			e.Status = unescape(value)
		case "SUMMARY":
			e.Summary = unescape(value)
		case "TRANSP":
			e.Transp = unescape(value)
		case "UID":
			e.UID = unescape(value)
		case "URL":
			e.URL = unescape(value)
		case "RECURRENCE-ID":
			e.RecurrenceID = unescape(value)
		case "DTEND":
			e.DtEnd, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: DTEND")
			}
		case "DURATION":
			e.Duration, err = parseDuration(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: DURATION")
			}
		case "ATTACH":
			e.Attatch = unescape(value)
		case "ATTENDEE":
			e.Attendee = unescape(value)
		case "CATEGORIES":
			if len(value) != 0 {
				e.Categories = strings.Split(value, ",")
			}
		case "COMMENT":
			e.Comment = unescape(value)
		case "CONTACT":
			e.Contact = unescape(value)
		case "EXDATE":
			e.ExDate, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: EXDATE")
			}
		case "EXRULE":
			e.ExRule = unescape(value)
		case "RSTATUS":
			e.RStatus = unescape(value)
		case "RELATED":
			e.Related = unescape(value)
		case "RESOURCES":
			e.Resources = unescape(value)
		case "RDATE":
			e.RDate, err = parseDateTime(value)
			if err != nil {
				return errors.Wrap(err, "invalid format: RDATE")
			}
		case "RRULE":
			e.RRule = unescape(value)
		case "X-PROP":
			e.XProp = unescape(value)

		default:
			return fmt.Errorf("unsupported field %s", key)
		}

		return nil
	}

	s := bufio.NewScanner(d.r)

	inFolding := false
	line := ""
	for s.Scan() {
		text := s.Text()
		if strings.HasPrefix(text, "  ") {
			line += text[2:]
			inFolding = true
			continue
		}

		if inFolding {
			if err := feed(line); err != nil {
				return err
			}
			if err := feed(text); err != nil {
				return err
			}

			line = ""
			inFolding = false
			continue
		}

		if line != "" {
			if err := feed(line); err != nil {
				return err
			}
		}

		line = text
	}

	return nil
}

type eventEncoder struct {
	w io.Writer
}

var _ Encoder = (*eventEncoder)(nil)

func NewEventEncoder(w io.Writer) Encoder {
	return &eventEncoder{w: w}
}

func (enc *eventEncoder) Encode(v any) error {
	evt, ok := v.(*VEvent)
	if !ok {
		return fmt.Errorf("want *VEvent ")
	}

	fmt.Fprintf(enc.w, "BEGIN:VEVENT\r\n")
	if err := enc.writeFields([]goxp.Tuple2[string, any]{
		{"CLASS", evt.Class},
		{"CREATED", evt.Created},
		{"DESCRIPTION", evt.Description},
		{"DTSTART", evt.DtStart},
		{"GEO", evt.Geo},
		{"LAST-MODIFIED", evt.LastModied},
		{"LOCATION", evt.Location},
		{"ORGANIZER", evt.Organizer},
		{"ORGANIZER", evt.Organizer},
		{"PRIORITY", evt.Priority},
		{"DTSTAMP", evt.DtStamp},
		{"SEQ", evt.Seq},
		{"STATUS", evt.Status},
		{"SUMMARY", evt.Summary},
		{"TRANSP", evt.Transp},
		{"UID", evt.UID},
		{"URL", evt.URL},
		{"RECURRENCE-ID", evt.RecurrenceID},
		{"DTEND", evt.DtEnd},
		{"DURATION", evt.Duration},
		{"ATTACH", evt.Attatch},
		{"ATTENDEE", evt.Attendee},
		{"CATEGORIES", evt.Categories},
		{"COMMENT", evt.Comment},
		{"CONTACT", evt.Contact},
		{"EXDATE", evt.ExDate},
		{"EXRULE", evt.ExRule},
		{"RSTATUS", evt.RStatus},
		{"RELATED", evt.Related},
		{"RESOURCES", evt.Resources},
		{"RDATE", evt.RDate},
		{"RRULE", evt.RRule},
		{"X-PROP", evt.XProp},
	}); err != nil {
		return err
	}
	fmt.Fprintf(enc.w, "END:VEVENT\r\n")
	return nil
}

func (enc *eventEncoder) writeFields(values []goxp.Tuple2[string, any]) (err error) {
	for _, value := range values {
		if err := enc.writeField(value.V1, value.V2); err != nil {
			return err
		}
	}

	return nil
}

func (enc *eventEncoder) writeField(field string, value any) (err error) {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			break
		}

		v = escape(v)

		if len(v)+len(field) < 78 {
			_, err = fmt.Fprintf(enc.w, "%s:%s\r\n", field, v)
			break
		}

		i := 0
		iter.Chunk(iter.Of([]rune(fmt.Sprintf("%s:%s", field, v))...), 78).Each(
			func(r []rune) {
				if i != 0 {
					fmt.Fprintf(enc.w, "  ")
				}
				fmt.Fprintf(enc.w, "%s\r\n", string(r))
				i++
			})

	case []string:
		if len(v) != 0 {
			_, err = fmt.Fprintf(enc.w, "%s:%s\r\n", field, strings.Join(v, ","))
		}
	case DateTime:
		if !v.IsZero() {
			_, err = fmt.Fprintf(enc.w, "%s:%s\r\n", field, v.String())
		}
	case Duration:
		if v.Duration != 0 {
			_, err = fmt.Fprintf(enc.w, "%s:%s\r\n", field, v.String())
		}
	case int:
		if v != 0 {
			_, err = fmt.Fprintf(enc.w, "%s:%d\r\n", field, v)
		}
	default:
		err = fmt.Errorf("unsupported data type: %T", v)
	}

	return err
}
