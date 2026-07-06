package domain

import (
	"encoding/json"
	"testing"
)

func TestGeneralResponseOmitsEmptyData(t *testing.T) {
	response := GeneralResponse[any]{
		Success: false,
		Message: "failed",
	}

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	expected := `{"success":false,"message":"failed"}`
	if string(body) != expected {
		t.Fatalf("expected %s, got %s", expected, body)
	}
}

func TestListUsersResponseIncludesTotalAndData(t *testing.T) {
	response := ListUsersResponse{
		Success: true,
		Message: "users fetched successfully",
		Data: []User{
			{
				UUID:        "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
				Fullname:    "Jane Doe",
				Email:       "jane@example.com",
				Description: "Example user",
			},
		},
		Total: 1,
	}

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if decoded["total"] != float64(1) {
		t.Fatalf("expected total 1, got %#v", decoded["total"])
	}

	if _, ok := decoded["data"]; !ok {
		t.Fatal("expected data field to be present")
	}
}
