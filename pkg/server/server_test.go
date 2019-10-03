package server

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ms-choudhary/slackup/pkg/store"
	"github.com/ms-choudhary/slackup/pkg/util"
)

var (
	testDBFile    = "test.db"
	migrationFile = "../../scripts/migration.sql"
)

func expectNoError(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServeHTTP(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{
			name:   "InvalidPath-Empty",
			method: "GET",
			path:   "https://example.com/",
			status: 404,
		},
		{
			name:   "InvalidPath-Less",
			method: "GET",
			path:   "https://example.com/asdf",
			status: 404,
		},
		{
			name:   "InvalidPath-More",
			method: "GET",
			path:   "https://example.com/asdf/asdf/asdf",
			status: 404,
		},
		{
			name:   "InvalidMethod",
			method: "POST",
			path:   "https://example.com/asdf/asdf",
			status: 404,
		},
		{
			name:   "ValidPathWithoutQuery",
			method: "GET",
			path:   "https://example.com/asdf/asdf",
			status: 200,
		},
		{
			name:   "ValidPathWithQuery",
			method: "GET",
			path:   "https://example.com/asdf/asdf?user=mohit&text=hello",
			status: 200,
		},
	}

	err := util.SetupDatabaseFrom(testDBFile, migrationFile)
	expectNoError(err, t)

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			store, err := store.Init(testDBFile)
			expectNoError(err, t)

			server := &Server{
				Store: store,
			}

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			response := w.Result()
			if response.StatusCode != tc.status {
				t.Fatalf("expected status %v got %v", tc.status, response.StatusCode)
			}
		})
	}

	os.Remove(testDBFile)
}
