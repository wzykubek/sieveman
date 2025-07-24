package client

import "testing"

func TestParseInlineResponse(t *testing.T) {
	tests := []struct {
		input         string
		expectedResp  Response
		expectedBytes int
		expectedErr   bool
	}{
		{
			input: "OK",
			expectedResp: Response{
				Name:    "OK",
				Code:    ResponseCode{Name: "", Message: ""},
				Message: "",
			},
			expectedBytes: 0,
			expectedErr:   false,
		},
		{
			input: `OK "LISTSCRIPTS completed"`,
			expectedResp: Response{
				Name:    "OK",
				Code:    ResponseCode{Name: "", Message: ""},
				Message: "LISTSCRIPTS completed",
			},
			expectedBytes: 0,
			expectedErr:   false,
		},
		{
			input: `OK (TAG "string with spaces")`,
			expectedResp: Response{
				Name:    "OK",
				Code:    ResponseCode{Name: "TAG", Message: "string with spaces"},
				Message: "",
			},
			expectedBytes: 0,
			expectedErr:   false,
		},
		{
			input: `NO (QUOTA/MAXSIZE) "Quota exceeded"`,
			expectedResp: Response{
				Name:    "NO",
				Code:    ResponseCode{Name: "QUOTA/MAXSIZE", Message: ""},
				Message: "Quota exceeded",
			},
			expectedBytes: 0,
			expectedErr:   false,
		},
		{
			input: `No (NONEXISTENT) {31}`,
			expectedResp: Response{
				Name:    "NO",
				Code:    ResponseCode{Name: "NONEXISTENT", Message: ""},
				Message: "",
			},
			expectedBytes: 31,
			expectedErr:   false,
		},
		{
			input: `NO {131}`,
			expectedResp: Response{
				Name:    "NO",
				Code:    ResponseCode{Name: "", Message: ""},
				Message: "",
			},
			expectedBytes: 131,
			expectedErr:   false,
		},
		{
			input: `NO {131}`,
			expectedResp: Response{
				Name:    "NO",
				Code:    ResponseCode{Name: "", Message: ""},
				Message: "",
			},
			expectedBytes: 131,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			resp, bytes, err := parseInlineResponse(tt.input)

			if tt.expectedErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if resp.Name != tt.expectedResp.Name {
				t.Errorf("Expected response name %q, got %q", tt.expectedResp.Name, resp.Name)
			}
			if resp.Code.Name != tt.expectedResp.Code.Name {
				t.Errorf("Expected response code name %q, got %q", tt.expectedResp.Code.Name, resp.Code.Name)
			}
			if resp.Code.Message != tt.expectedResp.Code.Message {
				t.Errorf("Expected response code message %q, got %q", tt.expectedResp.Code.Message, resp.Code.Message)
			}
			if resp.Message != tt.expectedResp.Message {
				t.Errorf("Expected response message %q, got %q", tt.expectedResp.Message, resp.Message)
			}

			if bytes != tt.expectedBytes {
				t.Errorf("Expected bytes %d, got %d", tt.expectedBytes, bytes)
			}
		})
	}
}
