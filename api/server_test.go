package api

import (
	"file-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestFileAuthorizationPolicy(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddFiles(nil)

	t.Run("metadata by id requires authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/files/1", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
		}
	})

	t.Run("metadata list allows internal calls without authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/files?id=1", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code == http.StatusUnauthorized {
			t.Fatalf("status = %d, route must stay open for internal service calls", response.Code)
		}
	})
}
