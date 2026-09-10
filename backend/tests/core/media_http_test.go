package core_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/tests"
)

func TestMediaEndpointUploadsListsPreviewsAndDeletes(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, csrf := loginMediaTestAdmin(t, testCase)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="note.txt"`)
	header.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	_, err = part.Write([]byte("hello media"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	created, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		WithHeader("Content-Type", writer.FormDataContentType()).
		Post("/api/media", &body)
	require.NoError(t, err)
	created.AssertStatus(201)
	payload, err := created.Json()
	require.NoError(t, err)
	item := payload["data"].(map[string]any)
	id := item["id"].(string)

	list, err := testCase.Http(t).WithCookie(cookie).Get("/api/media")
	require.NoError(t, err)
	list.AssertOk()
	listPayload, err := list.Json()
	require.NoError(t, err)
	require.Len(t, listPayload["data"], 1)

	preview, err := testCase.Http(t).WithCookie(cookie).Get("/api/media/" + id + "/preview")
	require.NoError(t, err)
	preview.AssertOk()

	deleted, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Delete("/api/media/"+id, nil)
	require.NoError(t, err)
	deleted.AssertOk()
}

func loginMediaTestAdmin(t *testing.T, testCase *tests.TestCase) (cookie *http.Cookie, csrf string) {
	t.Helper()
	csrfResponse, err := testCase.Http(t).Get("/csrf")
	require.NoError(t, err)
	sessionCookie := csrfResponse.Cookie("goravel_session")
	require.NotNil(t, sessionCookie)

	login, err := testCase.Http(t).
		WithCookie(sessionCookie).
		WithHeader("X-CSRF-TOKEN", csrfResponse.Headers().Get("X-CSRF-TOKEN")).
		Post("/login", bytes.NewBufferString(`{"email":"admin@example.com","password":"test-only-password"}`))
	require.NoError(t, err)
	login.AssertOk()
	rotated := login.Cookie("goravel_session")
	require.NotNil(t, rotated)
	return rotated, login.Headers().Get("X-CSRF-TOKEN")
}
