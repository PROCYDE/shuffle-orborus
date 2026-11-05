package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	shuffle "github.com/shuffle/shuffle-shared"
)

func TestExecutionRerunAcknowledgement(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		wantErr bool
	}{
		{"queued", 200, `{"success":true}`, false},
		{"terminal", 200, `{"success":true,"reason":"Already finished"}`, false},
		{"rejected", 200, `{"success":false}`, true},
		{"missing acknowledgement", 200, `{}`, true},
		{"malformed", 200, `not JSON`, true},
		{"unavailable", 503, `{"success":false}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if request.Method != http.MethodPost || request.URL.Path != "/api/v1/workflows/workflow/executions/execution/rerun" || request.Header.Get("Authorization") != "Bearer fixture-auth" {
					t.Error("incorrect recovery request")
				}
				response.WriteHeader(tc.status)
				fmt.Fprint(response, tc.body)
			}))
			defer server.Close()
			err := requestExecutionRerun(server.Client(), server.URL, shuffle.ExecutionRequest{ExecutionId: "execution", WorkflowId: "workflow", Authorization: "fixture-auth"})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}
