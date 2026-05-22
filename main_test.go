// Package main contains the test suite for the go-web-app HTTP server.
//
// Tests use Go's built-in net/http/httptest package to run handlers directly
// without opening a real network connection. This keeps tests fast, isolated,
// and free of external dependencies.
//
// All route tests are written using the table-driven pattern, which is the
// standard Go style. Each test case is a row in a table — adding a new route
// to test only requires adding one new row.
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// ---------------------------------------------------------------------------
// Route Tests
// ---------------------------------------------------------------------------

// TestRoutes checks that each HTTP handler returns the correct status code
// and Content-Type header.
func TestRoutes(t *testing.T) {
	cases := []struct {
		name            string
		handler         http.HandlerFunc
		path            string
		wantStatus      int
		wantContentType string
	}{
		{
			name:            "home page returns 200 with HTML content-type",
			handler:         homePage,
			path:            "/home",
			wantStatus:      http.StatusOK,
			wantContentType: "text/html; charset=utf-8",
		},
		{
			name:            "liveness probe returns 200 with plain text",
			handler:         healthzHandler,
			path:            "/healthz",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name:            "readiness probe returns 200 with plain text",
			handler:         readyzHandler,
			path:            "/readyz",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range cases {
		// t.Run creates an isolated sub-test for each case.
		// Failed cases are reported individually, and you can run a single
		// case with: go test -run TestRoutes/liveness_probe
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// ResponseRecorder captures the handler's output in memory.
			// No real TCP connection is made.
			rr := httptest.NewRecorder()
			tc.handler.ServeHTTP(rr, req)

			// Check the HTTP status code.
			if got := rr.Code; got != tc.wantStatus {
				t.Errorf("status code: got %d, want %d", got, tc.wantStatus)
			}

			// Check the Content-Type header.
			// A wrong Content-Type can cause browsers to misinterpret the
			// response or trigger MIME-sniffing security issues.
			if got := rr.Header().Get("Content-Type"); got != tc.wantContentType {
				t.Errorf("Content-Type: got %q, want %q", got, tc.wantContentType)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Health Endpoint Body Tests
// ---------------------------------------------------------------------------

// TestHealthzBody checks that the liveness probe returns exactly "ok".
// Some monitoring tools (like Prometheus blackbox exporter) parse the
// response body, so the content is part of the public contract.
func TestHealthzBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	healthzHandler(rr, req)

	want := "ok\n"
	if got := rr.Body.String(); got != want {
		t.Errorf("healthz body: got %q, want %q", got, want)
	}
}

// TestReadyzBody checks that the readiness probe returns exactly "ready".
func TestReadyzBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	readyzHandler(rr, req)

	want := "ready\n"
	if got := rr.Body.String(); got != want {
		t.Errorf("readyz body: got %q, want %q", got, want)
	}
}
