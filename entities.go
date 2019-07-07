package galice

import (
	"encoding/json"
	"fmt"
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
		return []byte("SimpleUtterance"), nil
	}
	if r == RequestTypeButtonPressed {
		return []byte("ButtonPressed"), nil
	}

	return []byte{}, fmt.Errorf("Unsupported RequestType value: %v", r)
}

// UnmarshalJSON converts Alice API type of request into internal RequestType value
func (r RequestType) UnmarshalJSON(input []byte) error {
	str := string(input)
	if str == "SimpleUtterance" {
		r = RequestTypeSimpleUtterance
		return nil
	}
	if str == "ButtonPressed" {
		r = RequestTypeButtonPressed
		return nil
	}

	return fmt.Errorf("Unsupported RequestType value: %v", input)
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
		return []byte("YANDEX.DATETIME"), nil
	}
	if e == EntityTypeFIO {
		return []byte("YANDEX.FIO"), nil
	}
	if e == EntityTypeGeo {
		return []byte("YANDEX.GEO"), nil
	}
	if e == EntityTypeNumber {
		return []byte("YANDEX.NUMBER"), nil
	}

	return []byte{}, fmt.Errorf("Unsupported EntityType value: %v", e)
}

// UnmarshalJSON converts Alice API type of named entity into internal EntityType value
func (e EntityType) UnmarshalJSON(input []byte) error {
	str := string(input)
	if str == "YANDEX.DATETIME" {
		e = EntityTypeDateTime
		return nil
	}
	if str == "YANDEX.FIO" {
		e = EntityTypeFIO
		return nil
	}
	if str == "YANDEX.GEO" {
		e = EntityTypeGeo
		return nil
	}
	if str == "YANDEX.NUMBER" {
		e = EntityTypeNumber
		return nil
	}

	return fmt.Errorf("Unsupported EntityType value: %v", input)
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
