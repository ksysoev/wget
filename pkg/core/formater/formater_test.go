package formater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat_FormatMessage(t *testing.T) {
	formater := NewFormat()

	formattedTextMsg, err := formater.FormatMessage("Request", "TestFormat_FormatMessage")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTextMsg := "TestFormat_FormatMessage"
	if formattedTextMsg != expectedTextMsg {
		t.Errorf("Unexpected formatted message: %v", formattedTextMsg)
	}

	formattedJSONMsg, err := formater.FormatMessage("Response", `{"status": 200, "body": "TestFormat_FormatMessage"}`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedJSONMsg := "{\n  \"body\": \"TestFormat_FormatMessage\",\n  \"status\": 200\n}"
	if formattedJSONMsg != expectedJSONMsg {
		t.Errorf("Unexpected formatted message: %v, wanted %v", formattedJSONMsg, expectedJSONMsg)
	}

	testString := `{"status": 200, "body": "TestFormat_FormatMessage"`

	formattedInvalidJSONMsg, err := formater.FormatMessage("Response", testString)
	if err != nil {
		t.Errorf("Expected to get no error, but got %v", err)
	}

	if formattedInvalidJSONMsg != testString {
		t.Errorf("Expected formated plain string, but got %v", formattedInvalidJSONMsg)
	}

	formattedUnknownMsg, err := formater.FormatMessage("NotDefined", "unknown message type")
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if formattedUnknownMsg != "" {
		t.Errorf("Expected empty string, but got %v", formattedUnknownMsg)
	}
}

func TestFormat_FormatForFile(t *testing.T) {
	formater := NewFormat()

	formattedTextMsg, err := formater.FormatForFile("Request", "TestFormat_FormatForFile")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTextMsg := "TestFormat_FormatForFile"
	if formattedTextMsg != expectedTextMsg {
		t.Errorf("Unexpected formatted message: %v", formattedTextMsg)
	}

	formattedJSONMsg, err := formater.FormatForFile("Response", `{"status": 200, "body": "TestFormat_FormatForFile"}`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedJSONMsg := `{"body":"TestFormat_FormatForFile","status":200}`
	if formattedJSONMsg != expectedJSONMsg {
		t.Errorf("Unexpected formatted message: %v", formattedJSONMsg)
	}

	formattedInvalidJSONMsg, err := formater.FormatForFile("Response", `{"status": 200, "body": "TestFormat_FormatForFile"`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedInvalidJSONMsg := "{\"status\": 200, \"body\": \"TestFormat_FormatForFile\""
	if formattedInvalidJSONMsg != expectedInvalidJSONMsg {
		t.Errorf("Unexpected formatted message: %v", formattedInvalidJSONMsg)
	}
}

func TestFormat_formatTextMessage(t *testing.T) {
	formater := NewFormat()

	requestData := "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"

	formattedRequestMsg, err := formater.formatTextMessage("Request", requestData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedRequestMsg := "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"
	if formattedRequestMsg != expectedRequestMsg {
		t.Errorf("Unexpected formatted message: %v", formattedRequestMsg)
	}

	// Test response message formatting
	responseData := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, world!"

	formattedResponseMsg, err := formater.formatTextMessage("Response", responseData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponseMsg := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, world!"
	if formattedResponseMsg != expectedResponseMsg {
		t.Errorf("Unexpected formatted message: %v", formattedResponseMsg)
	}

	unknownData := "Hello, world!"

	formattedUnknownMsg, err := formater.formatTextMessage("NotDefined", unknownData)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if formattedUnknownMsg != "" {
		t.Errorf("Expected empty string, but got %v", formattedUnknownMsg)
	}
}

func TestFormat_formatJSONMessage(t *testing.T) {
	formater := NewFormat()

	// Test request message formatting as JSON
	requestData := map[string]interface{}{
		"status": "200",
		"body":   "Hello, world!",
	}

	formattedRequestMsg, err := formater.formatJSONMessage("Request", requestData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedRequestMsg := "{\n  \"body\": \"Hello, world!\",\n  \"status\": \"200\"\n}"

	if formattedRequestMsg != expectedRequestMsg {
		t.Errorf("Unexpected formatted message: %v", formattedRequestMsg)
	}

	// Test response message formatting as JSON
	responseData := map[string]interface{}{
		"status": "200",
		"body":   "Hello, world!",
	}

	formattedResponseMsg, err := formater.formatJSONMessage("Response", responseData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponseMsg := "{\n  \"body\": \"Hello, world!\",\n  \"status\": \"200\"\n}"
	if formattedResponseMsg != expectedResponseMsg {
		t.Errorf("Unexpected formatted message: %v", formattedResponseMsg)
	}

	// Test unknown message type
	unknownData := "Hello, world!"

	formattedUnknownMsg, err := formater.formatJSONMessage("NotDefined", unknownData)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if formattedUnknownMsg != "" {
		t.Errorf("Expected empty string, but got %v", formattedUnknownMsg)
	}
}

func TestFormat_parseJSON(t *testing.T) {
	formater := NewFormat()

	// Test valid JSON parsing
	validJSON := `{"status": 200, "body": "Hello, world!"}`

	parsedValidJSON, ok := formater.parseJSON(validJSON)
	if !ok {
		t.Errorf("Expected true, but got false")
	}

	expectedValidJSON := map[string]interface{}{
		"status": 200.0,
		"body":   "Hello, world!",
	}

	assert.Equal(t, expectedValidJSON, parsedValidJSON)

	// Test invalid JSON parsing
	invalidJSON := `{"status": 200, "body": "Hello, world!"`

	parsedInvalidJSON, ok := formater.parseJSON(invalidJSON)

	assert.False(t, ok)
	assert.Nil(t, parsedInvalidJSON)
}
