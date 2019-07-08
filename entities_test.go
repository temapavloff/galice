package galice

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestType(t *testing.T) {
	t1 := RequestTypeSimpleUtterance
	s1, err := json.Marshal(t1)
	require.NoError(t, err)
	require.Equal(t, "\"SimpleUtterance\"", string(s1))

	t2 := RequestTypeButtonPressed
	s2, err := json.Marshal(t2)
	require.NoError(t, err)
	require.Equal(t, "\"ButtonPressed\"", string(s2))

	s3 := []byte(`{"type1": "SimpleUtterance", "type2": "ButtonPressed"}`)
	var m map[string]RequestType
	err = json.Unmarshal(s3, &m)
	require.NoError(t, err)
	require.Equal(t, RequestTypeSimpleUtterance, m["type1"])
	require.Equal(t, RequestTypeButtonPressed, m["type2"])
}

func TestEntityType(t *testing.T) {
	e1 := EntityTypeDateTime
	s1, err := json.Marshal(e1)
	require.NoError(t, err)
	require.Equal(t, "\"YANDEX.DATETIME\"", string(s1))

	e2 := EntityTypeFIO
	s2, err := json.Marshal(e2)
	require.NoError(t, err)
	require.Equal(t, "\"YANDEX.FIO\"", string(s2))

	e3 := EntityTypeGeo
	s3, err := json.Marshal(e3)
	require.NoError(t, err)
	require.Equal(t, "\"YANDEX.GEO\"", string(s3))

	e4 := EntityTypeNumber
	s4, err := json.Marshal(e4)
	require.NoError(t, err)
	require.Equal(t, "\"YANDEX.NUMBER\"", string(s4))

	s5 := []byte(`{"e1": "YANDEX.DATETIME", "e2": "YANDEX.FIO", "e3": "YANDEX.GEO", "e4": "YANDEX.NUMBER"}`)
	var m map[string]EntityType
	err = json.Unmarshal(s5, &m)
	require.NoError(t, err)
	require.Equal(t, EntityTypeDateTime, m["e1"])
	require.Equal(t, EntityTypeFIO, m["e2"])
	require.Equal(t, EntityTypeGeo, m["e3"])
	require.Equal(t, EntityTypeNumber, m["e4"])
}

func TestInputData(t *testing.T) {
	input := `{
	"meta": {
		"locale": "ru-RU",
		"timezone": "Europe/Moscow",
		"client_id": "ru.yandex.searchplugin/5.80 (Samsung Galaxy; Android 4.4)",
		"interfaces": {
			"screen": { }
		}
	},
	"request": {
		"command": "закажи пиццу на улицу льва толстого 16 на завтра",
		"original_utterance": "закажи пиццу на улицу льва толстого, 16 на завтра",
		"type": "SimpleUtterance",
		"markup": {
		"dangerous_context": true
		},
		"payload": {},
		"nlu": {
		"tokens": [
			"закажи",
			"пиццу",
			"на",
			"льва",
			"толстого",
			"16",
			"на",
			"завтра"
		],
		"entities": [
			{
				"tokens": {
					"start": 2,
					"end": 6
				},
				"type": "YANDEX.GEO",
				"value": {
					"house_number": "16",
					"street": "льва толстого"
				}
			},
			{
				"tokens": {
					"start": 3,
					"end": 5
				},
				"type": "YANDEX.FIO",
				"value": {
					"first_name": "лев",
					"last_name": "толстой"
				}
			},
			{
				"tokens": {
					"start": 5,
					"end": 6
				},
				"type": "YANDEX.NUMBER",
				"value": 16
				},
			{
				"tokens": {
					"start": 6,
					"end": 8
				},
				"type": "YANDEX.DATETIME",
				"value": {
					"day": 1,
					"day_is_relative": true
				}
			}
		]
		}
	},
	"session": {
		"new": true,
		"message_id": 4,
		"session_id": "2eac4854-fce721f3-b845abba-20d60",
		"skill_id": "3ad36498-f5rd-4079-a14b-788652932056",
		"user_id": "AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"
	},
	"version": "1.0"
}`

	i := InputData{}
	err := json.Unmarshal([]byte(input), &i)
	require.NoError(t, err)
}

func TestRequestIsPing(t *testing.T) {
	requestPing := []byte(`{"original_utterance": "ping"}`)
	requestNotPing := []byte(`{"original_utterance": "blah blah"}`)

	var r Request
	var err error

	err = json.Unmarshal(requestPing, &r)
	require.NoError(t, err)
	require.True(t, r.IsPing())

	err = json.Unmarshal(requestNotPing, &r)
	require.NoError(t, err)
	require.False(t, r.IsPing())
}

func TestRequestParsePayload(t *testing.T) {
	var r Request
	var p []int
	var err error

	r.Payload = json.RawMessage(`[1,2,3,4,5]`)
	err = r.DecodePayload(&p)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4, 5}, p)

	r.Payload = json.RawMessage(`5q46hdafbsgL:^5wbgdgd`)
	err = r.DecodePayload(&p)
	require.Error(t, err)
}

func TestEntityValues(t *testing.T) {
	geoStr := []byte(`{
		"tokens": {
		  "start": 2,
		  "end": 6
		},
		"type": "YANDEX.GEO",
		"value": {
		  "house_number": "16",
		  "street": "льва толстого"
		}
	}`)
	fioStr := []byte(`{
		"tokens": {
		  "start": 3,
		  "end": 5
		},
		"type": "YANDEX.FIO",
		"value": {
		  "first_name": "лев",
		  "last_name": "толстой"
		}
	}`)
	floatStr := []byte(`{
		"tokens": {
		  "start": 5,
		  "end": 6
		},
		"type": "YANDEX.NUMBER",
		"value": 5.5
	}`)
	intStr := []byte(`{
		"tokens": {
		  "start": 5,
		  "end": 6
		},
		"type": "YANDEX.NUMBER",
		"value": 16
	}`)

	var ge RequestEntity
	var fioe RequestEntity
	var fe RequestEntity
	var ie RequestEntity
	var err error

	err = json.Unmarshal(geoStr, &ge)
	require.NoError(t, err)
	err = json.Unmarshal(fioStr, &fioe)
	require.NoError(t, err)
	err = json.Unmarshal(floatStr, &fe)
	require.NoError(t, err)
	err = json.Unmarshal(intStr, &ie)
	require.NoError(t, err)

	_, err = ge.GeoValue()
	require.NoError(t, err)
	_, err = ge.FIOValue()
	require.Error(t, err)

	_, err = fioe.FIOValue()
	require.NoError(t, err)
	_, err = fioe.GeoValue()
	require.Error(t, err)

	fv, err := fe.FloatValue()
	require.Equal(t, 5.5, fv)
	require.NoError(t, err)
	_, err = fe.GeoValue()
	require.Error(t, err)

	iv, err := ie.IntValue()
	require.Equal(t, 16, iv)
	require.NoError(t, err)
	_, err = ie.GeoValue()
	require.Error(t, err)
}

func TestValueDateTimeAbsolute(t *testing.T) {
	strAbs := []byte(`{
		"year": 1982,
		"month": 9,
		"day": 15,
		"hour": 22,
		"minute": 30
	}`)
	locAbs, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)
	timeAbs := time.Date(1982, time.Month(9), 15, 22, 30, 0, 0, locAbs)
	var vl ValueDateTime
	err = json.Unmarshal(strAbs, &vl)
	require.NoError(t, err)
	require.False(t, vl.IsRelative())
	tv, err := vl.Time("Europe/Moscow")
	require.NoError(t, err)
	require.True(t, tv.Equal(timeAbs))
	_, err = vl.Time("NO_EXISTING_ZONE_ID")
	require.Error(t, err)
}

func TestValueDateTimeRelative(t *testing.T) {
	strRel := []byte(`{
		"year": 3,
		"year_is_relative": true,
		"month": 2,
		"month_is_relative": true,
		"day": -8,
		"day_is_relative": true,
		"hour": -4,
		"hour_is_relative": true,
		"minute": -30,
		"minute_is_relative": true
	}`)
	locRel, err := time.LoadLocation("Europe/Moscow")
	require.NoError(t, err)
	now := time.Now()
	timeRel := time.Date(now.Year()+3, now.Month()+time.Month(2),
		now.Day()-8, now.Hour()-4, now.Minute()-30, 0, 0, locRel)
	var vl ValueDateTime
	err = json.Unmarshal(strRel, &vl)
	require.NoError(t, err)
	require.True(t, vl.IsRelative())
	tv, err := vl.Time("Europe/Moscow")
	require.NoError(t, err)
	require.True(t, tv.Equal(timeRel))
}
