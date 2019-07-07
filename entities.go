package galice

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Meta - Alise request metadata
type Meta struct {
	Locale     string      `json:"locale"`
	Timezone   string      `json:"timezone"`
	ClientID   string      `json:"client_id"`
	Interfaces interface{} `json:"interfaces"`
}

// Session - Alise session information
type Session struct {
	New       bool   `json:"new"`
	MessageID uint   `json:"message_id"`
	SessionID string `json:"session_id"`
	SkillID   string `json:"skill_id"`
	UserID    string `json:"user_id"`
}

// RequestType - type of Alice request: SimpleUtterance or ButtonPressed
type RequestType uint8

// MarshalJSON converts inner representation to values supported by Alise API
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

// RequestMarkup - Alice request markup metadatas
type RequestMarkup struct {
	DangerousContext bool `json:"dangerous_context"`
}

// EntityType - type of ALice requst named entity: YANDEX.DATETIME, YANDEX.FIO, YANDEX.GEO, YANDEX.NUMBER
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

// MarshalJSON converts inner representation to values supported by Alise API
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

// ValueFIO - value type for entities YANDEX.FIO
type ValueFIO struct {
	FirstName      string `json:"first_name"`
	PatronymicName string `json:"patronymic_name"`
	LastName       string `json:"last_name"`
}

// ValueGeo - value type for entities YANDEX.GEO
type ValueGeo struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	Airport     string `json:"airport"`
}

// RequestEntity - machine representation of Alice requst named entity
type RequestEntity struct {
	Tokens struct {
		Start uint `json:"start"`
		End   uint `json:"end"`
	} `json:"tokens"`
	Type  EntityType      `json:"type"`
	Value json.RawMessage `json:"value"` // TODO !!!
}

// IsFIO - checks if RequestEntity is YANDEX.FIO
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

// TODO: Parse Alice request properly!!!
// DateTimeValue returns time.Time if RequestEntity is YANDEX.DATETIME or error otherwhise
func (e *RequestEntity) DateTimeValue() (time.Time, error) {
	if !e.IsDateTime() {
		return time.Now(), fmt.Errorf("Cannot create time.Time for entity type %v", e.Type)
	}

	return time.Now(), nil
}

// RequestNLU - words and names entities of Alice request
type RequestNLU struct {
	Tokens   []string        `json:"tokens"`
	Entities []RequestEntity `json:"entities"`
}

// Request - Alise request
type Request struct {
	Command           string          `json:"command"`
	OriginalUtterance string          `json:"original_utterance"`
	Type              RequestType     `json:"type"`
	Markup            RequestMarkup   `json:"markup"`
	Payload           json.RawMessage `json:"payload"`
}

// DecodePayload decodes current request payload into provied variable
func (r *Request) DecodePayload(p interface{}) error {
	if err := json.Unmarshal(r.Payload, p); err != nil {
		return fmt.Errorf("Unable to decode request payload: %v", err)
	}
	return nil
}

// IsPing returns `true` if it is Yandex healthcheck
func (r *Request) IsPing() bool {
	return r.OriginalUtterance == "ping"
}

// InputData - incoming data from Alice API
type InputData struct {
	Version string  `json:"version"`
	Meta    Meta    `json:"meta"`
	Session Session `json:"session"`
	Request Request `json:"request"`
}

// ResponseButton - Alice API representation of Button for response
type ResponseButton struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload"`
	URL     string      `json:"url"`
	Hide    string      `json:"hide"`
}

// Response - Alice response
type Response struct {
	Text       string         `json:"text"`
	TTS        string         `json:"tts"`
	Buttons    ResponseButton `json:"buttons"`
	EndSession bool           `json:"end_session"`
}

// OutputData - outcoming data for Alice API
type OutputData struct {
	Version  string   `json:"version"`
	Session  Session  `json:"session"`
	Response Response `json:"response"`
}

func isJSONNumberIsFloat(v json.RawMessage) bool {
	return strings.Contains(string(v), ".")
}
