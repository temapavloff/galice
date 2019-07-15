package galice

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Meta is an Alice API request metadata
type Meta struct {
	Locale     string      `json:"locale"`     // request locale
	Timezone   string      `json:"timezone"`   // request timezone
	ClientID   string      `json:"client_id"`  // request user client ID
	Interfaces interface{} `json:"interfaces"` // describes interfaces available on request user device
}

// Session is an Alice API session information
type Session struct {
	New       bool   `json:"new"`        // is it first request for current session or not
	MessageID uint   `json:"message_id"` // ID of current request
	SessionID string `json:"session_id"` // ID of current session
	SkillID   string `json:"skill_id"`   // ID of current skill
	UserID    string `json:"user_id"`    // ID of current user
}

// RequestType represents type of Alice API request: SimpleUtterance or ButtonPressed
type RequestType uint8

// MarshalJSON converts inner representation to values supported by Alice API
func (r RequestType) MarshalJSON() ([]byte, error) {
	if r == RequestTypeSimpleUtterance {
		return []byte("\"SimpleUtterance\""), nil
	}
	if r == RequestTypeButtonPressed {
		return []byte("\"ButtonPressed\""), nil
	}

	return []byte{}, fmt.Errorf("Unsupported RequestType value: %v", r)
}

// UnmarshalJSON converts Alice API type of request into internal RequestType value
func (r *RequestType) UnmarshalJSON(input []byte) error {
	str := string(input)
	if str == "\"SimpleUtterance\"" {
		*r = RequestTypeSimpleUtterance
		return nil
	}
	if str == "\"ButtonPressed\"" {
		*r = RequestTypeButtonPressed
		return nil
	}

	return fmt.Errorf("Unsupported RequestType value: %v", str)
}

const (
	// RequestTypeSimpleUtterance represents SimpleUtterance request type
	RequestTypeSimpleUtterance = RequestType(iota)
	// RequestTypeButtonPressed represents ButtonPressed request type
	RequestTypeButtonPressed
)

// RequestMarkup is Alice API request markup metadatas
type RequestMarkup struct {
	// true if current request message marked as dandegous (suicide, hate speech, threats) by Alice API
	DangerousContext bool `json:"dangerous_context"`
}

// EntityType is a type of ALice API request named entity:
// YANDEX.DATETIME, YANDEX.FIO, YANDEX.GEO, YANDEX.NUMBER
type EntityType uint8

const (
	// EntityTypeDateTime represents YANDEX.DATETIME
	EntityTypeDateTime = EntityType(iota)
	// EntityTypeFIO represents YANDEX.FIO
	EntityTypeFIO
	// EntityTypeGeo represents YANDEX.GEO
	EntityTypeGeo
	// EntityTypeNumber represents YANDEX.NUMBER
	EntityTypeNumber
)

// MarshalJSON converts inner representation to values supported by Alice API
func (e EntityType) MarshalJSON() ([]byte, error) {
	if e == EntityTypeDateTime {
		return []byte("\"YANDEX.DATETIME\""), nil
	}
	if e == EntityTypeFIO {
		return []byte("\"YANDEX.FIO\""), nil
	}
	if e == EntityTypeGeo {
		return []byte("\"YANDEX.GEO\""), nil
	}
	if e == EntityTypeNumber {
		return []byte("\"YANDEX.NUMBER\""), nil
	}

	return []byte{}, fmt.Errorf("Unsupported EntityType value: %v", e)
}

// UnmarshalJSON converts Alice API type of named entity into internal EntityType value
func (e *EntityType) UnmarshalJSON(input []byte) error {
	str := string(input)
	if str == "\"YANDEX.DATETIME\"" {
		*e = EntityTypeDateTime
		return nil
	}
	if str == "\"YANDEX.FIO\"" {
		*e = EntityTypeFIO
		return nil
	}
	if str == "\"YANDEX.GEO\"" {
		*e = EntityTypeGeo
		return nil
	}
	if str == "\"YANDEX.NUMBER\"" {
		*e = EntityTypeNumber
		return nil
	}

	return fmt.Errorf("Unsupported EntityType value: %v", str)
}

// ValueFIO is a value type for entities contains information of
// person firs, last, and patronymic  names
type ValueFIO struct {
	FirstName      string `json:"first_name"`
	PatronymicName string `json:"patronymic_name"`
	LastName       string `json:"last_name"`
}

// ValueGeo is a value type for entities contains information of
// some geografic object
type ValueGeo struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	Airport     string `json:"airport"`
}

// ValueDateTime is a value type for entities contains informtaion of
// relative or absolute data and time
type ValueDateTime struct {
	Year             int  `json:"year"`
	YearIsRelative   bool `json:"year_is_relative"`
	Month            int  `json:"month"`
	MonthIsRelative  bool `json:"month_is_relative"`
	Day              int  `json:"day"`
	DayIsRelative    bool `json:"day_is_relative"`
	Hour             int  `json:"hour"`
	HourIsRelative   bool `json:"hour_is_relative"`
	Minute           int  `json:"minute"`
	MinuteIsRelative bool `json:"minute_is_relative"`
}

// IsRelative return trus if ValueDateTime is in relative format, false otherwise
func (v *ValueDateTime) IsRelative() bool {
	return v.YearIsRelative || v.MonthIsRelative || v.DayIsRelative || v.HourIsRelative || v.MinuteIsRelative
}

// Time returns go time.Time based on ValueDateTime
func (v *ValueDateTime) Time(zoneID string) (time.Time, error) {
	location, err := time.LoadLocation(zoneID)
	if err != nil {
		return time.Now(), err
	}
	if v.IsRelative() {
		t := time.Now()
		return time.Date(
			t.Year()+v.Year,
			t.Month()+time.Month(v.Month),
			t.Day()+v.Day,
			t.Hour()+v.Hour,
			t.Minute()+v.Minute,
			0, 0, location,
		), nil
	}

	return time.Date(v.Year, time.Month(v.Month), v.Day, v.Hour, v.Minute, 0, 0, location), nil
}

// RequestEntity is a representation of Alice API request named entity
type RequestEntity struct {
	Tokens struct {
		Start uint `json:"start"`
		End   uint `json:"end"`
	} `json:"tokens"`
	Type  EntityType      `json:"type"`
	Value json.RawMessage `json:"value"`
}

// IsFIO checks if RequestEntity is YANDEX.FIO
func (e *RequestEntity) IsFIO() bool {
	return e.Type == EntityTypeFIO
}

// IsGeo checks if RequestEntity is YANDEX.GEO
func (e *RequestEntity) IsGeo() bool {
	return e.Type == EntityTypeGeo
}

// IsFloat checks if RequestEntity is floating point YANDEX.NUMBER
func (e *RequestEntity) IsFloat() bool {
	return e.Type == EntityTypeNumber && isJSONNumberIsFloat(e.Value)
}

// IsInt checks if RequestEntity is integer YANDEX.NUMBER
func (e *RequestEntity) IsInt() bool {
	return e.Type == EntityTypeNumber && !isJSONNumberIsFloat(e.Value)
}

// IsDateTime checks if RequestEntity is YANDEX.DATETIME
func (e *RequestEntity) IsDateTime() bool {
	return e.Type == EntityTypeDateTime
}

// FIOValue returns ValueFIO if RequestEntity is YANDEX.FIO or error otherwhise
func (e *RequestEntity) FIOValue() (ValueFIO, error) {
	var v ValueFIO

	if !e.IsFIO() {
		return v, fmt.Errorf("Cannot create ValueFIO for entity type %v", e.Type)
	}

	if err := json.Unmarshal(e.Value, &v); err != nil {
		return v, err
	}

	return v, nil
}

// GeoValue returns ValueGeo if RequestEntity is YANDEX.GEO or error otherwhise
func (e *RequestEntity) GeoValue() (ValueGeo, error) {
	var v ValueGeo

	if !e.IsGeo() {
		return v, fmt.Errorf("Cannot create ValueGeo for entity type %v", e.Type)
	}

	if err := json.Unmarshal(e.Value, &v); err != nil {
		return v, err
	}

	return v, nil
}

// FloatValue returns float if RequestEntity is floating point YANDEX.NUMBER or error otherwhise
func (e *RequestEntity) FloatValue() (float64, error) {
	var v float64

	if !e.IsFloat() {
		var add string
		if e.IsInt() {
			add = ", integer"
		}
		return v, fmt.Errorf("Cannot create float for entity type %v%v", e.Type, add)
	}

	if err := json.Unmarshal(e.Value, &v); err != nil {
		return v, err
	}

	return v, nil
}

// IntValue returns integer if RequestEntity is integer YANDEX.NUMBER or error otherwhise
func (e *RequestEntity) IntValue() (int, error) {
	var v int

	if !e.IsInt() {
		var add string
		if e.IsFloat() {
			add = ", float"
		}
		return v, fmt.Errorf("Cannot create integer for entity type %v%v", e.Type, add)
	}

	if err := json.Unmarshal(e.Value, &v); err != nil {
		return v, err
	}

	return v, nil
}

// DateTimeValue returns time.Time if RequestEntity is YANDEX.DATETIME or error otherwhise
func (e *RequestEntity) DateTimeValue() (ValueDateTime, error) {
	var v ValueDateTime

	if !e.IsDateTime() {
		return v, fmt.Errorf("Cannot create time.Time for entity type %v", e.Type)
	}

	if err := json.Unmarshal(e.Value, &v); err != nil {
		return v, err
	}

	return v, nil
}

// RequestNLU is a struct contains words and names entities of Alice API request
type RequestNLU struct {
	Tokens   []string        `json:"tokens"`
	Entities []RequestEntity `json:"entities"`
}

// Request is an Alice request
type Request struct {
	// User request converted for internal processing of Alice.
	// The text is cleared of punctuation marks, and the numerals are converted to numbers.
	Command string `json:"command"`
	// The full text of the user request, a maximum of 1024 characters.
	OriginalUtterance string `json:"original_utterance"`
	// The type of input, required field. May be RequestTypeSimpleUtterance or RequestTypeButtonPressed
	Type RequestType `json:"type"`
	// Replica Formal Characteristics
	Markup RequestMarkup `json:"markup"`
	// JSON, received with the button pressed from the skill handler (in response to the previous request)
	Payload json.RawMessage `json:"payload"`
	// Words and named entities retrieved from user request
	NLU RequestNLU `json:"nlu"`
}

// DecodePayload decodes current request payload into provied variable
func (r *Request) DecodePayload(p interface{}) error {
	if err := json.Unmarshal(r.Payload, p); err != nil {
		return fmt.Errorf("Unable to decode request payload: %v", err)
	}
	return nil
}

// IsPing checks if current request is Yandex healthcheck
func (r *Request) IsPing() bool {
	return r.OriginalUtterance == "ping"
}

// IsDangerousContext checks if current request has dangerous context
func (r *Request) IsDangerousContext() bool {
	return r.Markup.DangerousContext
}

// InputData is an incoming data from Alice API
type InputData struct {
	Version string  `json:"version"`
	Meta    Meta    `json:"meta"`
	Session Session `json:"session"`
	Request Request `json:"request"`
}

// ResponseButton is an Alice API representation of Button for response
type ResponseButton struct {
	Title   string      `json:"title"`
	Hide    bool        `json:"hide"`
	URL     string      `json:"url,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// Response is an Alice response
type Response struct {
	Text       string           `json:"text"`
	TTS        string           `json:"tts"`
	Buttons    []ResponseButton `json:"buttons,omitempty"`
	EndSession bool             `json:"end_session"`
}

// AddButton adds new button into current response
func (r *Response) AddButton(title string, hide bool, URL string, payload interface{}) {
	if r.Buttons == nil {
		r.Buttons = []ResponseButton{}
	}

	r.Buttons = append(r.Buttons, ResponseButton{title, hide, URL, payload})
}

// NewResponse creates new response. Use text variable to set response text message.
// Use tts variable to specify text to speach markup, if empty text value will be used.
// Use endSession flag to specify that current message is the last one in current session.
func NewResponse(text, tts string, endSession bool) Response {
	if tts == "" {
		tts = text
	}
	return Response{
		text,
		tts,
		nil,
		endSession,
	}
}

// OutputData is an outcoming data for Alice API
type OutputData struct {
	Version  string   `json:"version"`
	Session  Session  `json:"session"`
	Response Response `json:"response"`
}

// NewOutput creates new OutputData. Use i variable to provide InputDate to setup
// Alice API version and session data from it. Use r variable to set response
func NewOutput(i InputData, r Response) OutputData {
	return OutputData{
		Version:  i.Version,
		Session:  i.Session,
		Response: r,
	}
}

func pong(i InputData) OutputData {
	return OutputData{
		Version: i.Version,
		Session: i.Session,
		Response: Response{
			Text: "pong",
			TTS:  "pong",
		},
	}
}

func dangerous(i InputData) OutputData {
	return OutputData{
		Version: i.Version,
		Session: i.Session,
		Response: Response{
			Text: "Не понимаю, о чем вы. Пожалуйста, переформулируйте вопрос.",
			TTS:  "Не понимаю, о чем вы. Пожалуйста, переформулируйте вопрос.",
		},
	}
}

func isJSONNumberIsFloat(v json.RawMessage) bool {
	return strings.Contains(string(v), ".")
}
