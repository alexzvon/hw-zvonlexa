package internalhttp

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexzvon/hw12_13_14_15_calendar/internal/myutils"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	handler := &sHandler{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.root)
	mux.HandleFunc("/hello", handler.hello)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("test get root", func(t *testing.T) {
		cli := http.DefaultClient
		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)

		require.Nil(t, err)

		res, err := cli.Do(req)

		require.Nil(t, err)

		body, err := ioutil.ReadAll(res.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "Корень")

		res.Body.Close()
	})

	t.Run("test get hello", func(t *testing.T) {
		cli := http.DefaultClient
		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, myutils.ConCat(ts.URL, "/hello"), nil)

		require.Nil(t, err)

		res, err := cli.Do(req)

		require.Nil(t, err)

		body, err := ioutil.ReadAll(res.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "hello-world")

		res.Body.Close()
	})
}

func TestHandlerRoot(t *testing.T) {
	handler := &sHandler{}

	t.Run("test handler get root", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://site.ru/", nil)
		w := httptest.NewRecorder()

		handler.root(w, r)
		result := w.Result()

		require.Equal(t, http.StatusOK, result.StatusCode)

		body, err := ioutil.ReadAll(w.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "Корень")

		result.Body.Close()
	})

	t.Run("test handler post hello", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://site.ru", nil)
		w := httptest.NewRecorder()

		handler.root(w, r)
		result := w.Result()

		require.Equal(t, http.StatusOK, result.StatusCode)

		body, err := ioutil.ReadAll(w.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "Only GET allowed\n")

		result.Body.Close()
	})
}

func TestHandlerHello(t *testing.T) {
	handler := &sHandler{}

	t.Run("test handler get hello", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://site.ru/hello", nil)
		w := httptest.NewRecorder()

		handler.hello(w, r)
		result := w.Result()

		require.Equal(t, http.StatusOK, result.StatusCode)

		body, err := ioutil.ReadAll(w.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "hello-world")

		result.Body.Close()
	})

	t.Run("test handler post hello", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://site.ru/hello", nil)
		w := httptest.NewRecorder()

		handler.hello(w, r)
		result := w.Result()

		require.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)

		body, err := ioutil.ReadAll(w.Body)

		require.Nil(t, err)
		require.Equal(t, string(body), "Only GET allowed\n")

		result.Body.Close()
	})
}
