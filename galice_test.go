package galice

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	pingBody := `{
	"meta": {
		"locale": "ru-RU",
		"timezone": "Europe/Moscow",
		"client_id": "ru.yandex.searchplugin/5.80 (Samsung Galaxy; Android 4.4)",
		"interfaces": {
			"screen": { }
		}
	},
	"request": {
		"command": "ping",
		"original_utterance": "ping",
		"type": "SimpleUtterance",
		"markup": {
			"dangerous_context": false
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
	pongBody := `{"version":"1.0","session":{"new":true,"message_id":4,"session_id":"2eac4854-fce721f3-b845abba-20d60","skill_id":"3ad36498-f5rd-4079-a14b-788652932056","user_id":"AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"},"response":{"text":"pong","tts":"pong","end_session":false}}`
	cli := New(true, true)
	h := cli.CreateHandler(func(i InputData) (OutputData, error) {
		return OutputData{}, nil
	})

	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte(pingBody)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, pongBody, rr.Body.String())
	require.Equal(t, "application/json", rr.Header().Get("Content-type"))
}

func TestDangerousContext(t *testing.T) {
	dangerousBody := `{
	"meta": {
		"locale": "ru-RU",
		"timezone": "Europe/Moscow",
		"client_id": "ru.yandex.searchplugin/5.80 (Samsung Galaxy; Android 4.4)",
		"interfaces": {
			"screen": { }
		}
	},
	"request": {
		"command": "test",
		"original_utterance": "test",
		"type": "SimpleUtterance",
		"markup": {
			"dangerous_context": true
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
	respBody := `{"version":"1.0","session":{"new":true,"message_id":4,"session_id":"2eac4854-fce721f3-b845abba-20d60","skill_id":"3ad36498-f5rd-4079-a14b-788652932056","user_id":"AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"},"response":{"text":"Не понимаю, о чем вы. Пожалуйста, переформулируйте вопрос.","tts":"Не понимаю, о чем вы. Пожалуйста, переформулируйте вопрос.","end_session":false}}`
	cli := New(true, true)
	h := cli.CreateHandler(func(i InputData) (OutputData, error) {
		return OutputData{}, nil
	})

	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte(dangerousBody)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, respBody, rr.Body.String())
	require.Equal(t, "application/json", rr.Header().Get("Content-type"))
}

func TestHandlingExpectedError(t *testing.T) {
	cli := New(true, true)
	errStr := ""
	cli.SetLogger(func(err error) {
		errStr = err.Error()
	})
	h := cli.CreateHandler(func(i InputData) (OutputData, error) {
		return OutputData{}, errors.New("test")
	})
	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte("{}")))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "test", errStr)
}

func TestHandlingUnexpectedError(t *testing.T) {
	cli := New(true, true)
	errStr := ""
	cli.SetLogger(func(err error) {
		errStr = err.Error()
	})
	h := cli.CreateHandler(func(i InputData) (OutputData, error) {
		panic(errors.New("test"))
		return OutputData{}, errors.New("test")
	})
	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte("{}")))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Unexpected error: test", errStr)
}

func TestGeneralResponse(t *testing.T) {
	body := `{
	"meta": {
		"locale": "ru-RU",
		"timezone": "Europe/Moscow",
		"client_id": "ru.yandex.searchplugin/5.80 (Samsung Galaxy; Android 4.4)",
		"interfaces": {
			"screen": { }
		}
	},
	"request": {
		"command": "hi there",
		"original_utterance": "hi there",
		"type": "SimpleUtterance",
		"markup": {
			"dangerous_context": false
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
	resp := `{"version":"1.0","session":{"new":true,"message_id":4,"session_id":"2eac4854-fce721f3-b845abba-20d60","skill_id":"3ad36498-f5rd-4079-a14b-788652932056","user_id":"AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"},"response":{"text":"test","tts":"test","buttons":[{"title":"button 1","hide":true,"url":"https://ya.ru","payload":123},{"title":"button 2","hide":false}],"end_session":true}}`
	cli := New(true, true)
	h := cli.CreateHandler(func(i InputData) (OutputData, error) {
		r := NewResponse("test", "", true)
		r.AddButton("button 1", true, "https://ya.ru", 123)
		r.AddButton("button 2", false, "", nil)
		o := NewOutput(i, r)
		return o, nil
	})
	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte(body)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, resp, rr.Body.String())
}
