package web

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chriswalker/simplesrv/model"
)

// testItemService is a (currently bare-bones) test implementation of the
// item service.
type testItemService struct {
	// If set, methods will return this error.
	err error
}

func (s testItemService) GetItems() ([]model.Item, error) {
	return nil, s.err
}

func newTestServer(serviceErr error) (*app, error) {
	a := &app{
		svc:    testItemService{err: serviceErr},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	t, err := a.loadTemplates()
	if err != nil {
		return nil, err
	}

	a.templates = t

	return a, nil
}

func TestIndex(t *testing.T) {
	testCases := map[string]struct {
		expectedStatus int
		serviceErr     error
	}{
		"index-basic": {
			expectedStatus: http.StatusOK,
		},
		"index-error": {
			expectedStatus: http.StatusInternalServerError,
			serviceErr:     errors.New("error returning items list"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			s, err := newTestServer(tc.serviceErr)
			if err != nil {
				t.Fatal(err)
			}

			s.index(rr, r)
			if rr.Result().StatusCode != tc.expectedStatus {
				t.Errorf("got status code '%d', want '%d'",
					rr.Result().StatusCode, tc.expectedStatus)
			}
		})
	}
}
