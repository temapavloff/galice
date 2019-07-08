package galice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	pinbBody := `{
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
	pongBody := `{"version":"1.0","session":{"new":true,"message_id":4,"session_id":"2eac4854-fce721f3-b845abba-20d60","skill_id":"3ad36498-f5rd-4079-a14b-788652932056","user_id":"AC9WC3DF6FCE052E45A4566A48E6B7193774B84814CE49A922E163B8B29881DC"},"response":{"text":"pong","tts":"pong","end_session":false}}`
	cli := New(true, true)
	cli.SetLogger(func(err error) {
		t.Error(err)
	})
	h := cli.CreateHandler(func(i InputData) OutputData {
		return OutputData{}
	})

	req, err := http.NewRequest("POST", "/skill", bytes.NewReader([]byte(pinbBody)))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, pongBody, rr.Body.String())
}
