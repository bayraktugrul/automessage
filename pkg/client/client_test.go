package client

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Client_should_send_message_successfully_status_ok(t *testing.T) {
	// given
	mockResponse := Response{
		Message:   "Message sent successfully",
		MessageID: "msg-123456",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "+905551112233", req.To)
		assert.Equal(t, "Test message", req.Content)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()
	client := New(server.URL)

	// when
	response, err := client.SendMessage(context.Background(), Request{
		To:      "+905551112233",
		Content: "Test message",
	})

	// then
	assert.NoError(t, err)
	assert.Equal(t, mockResponse.Message, response.Message)
	assert.Equal(t, mockResponse.MessageID, response.MessageID)
}

func Test_Client_should_send_message_with_status_accepted_status(t *testing.T) {
	// given
	mockResponse := Response{
		Message:   "Message accepted for processing",
		MessageID: "msg-123456",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()
	client := New(server.URL)

	// when
	response, err := client.SendMessage(context.Background(), Request{
		To:      "+905551112233",
		Content: "Test message",
	})

	// then
	assert.NoError(t, err)
	assert.Equal(t, mockResponse.Message, response.Message)
	assert.Equal(t, mockResponse.MessageID, response.MessageID)
}

func Test_Client_should_send_message_and_handle_error_status(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}))
	defer server.Close()
	client := New(server.URL)

	// when
	_, err := client.SendMessage(context.Background(), Request{
		To:      "+905551112233",
		Content: "Test message",
	})

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed with status: 500")
}

func Test_Client_should_send_messages_and_handle_invalid_response(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	client := New(server.URL)

	// when
	_, err := client.SendMessage(context.Background(), Request{
		To:      "+905551112233",
		Content: "Test message",
	})

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding response")
}

func Test_Client_should_send_message_and_handle_network_error(t *testing.T) {
	// given
	client := New("http://non-existent-url")

	// when
	_, err := client.SendMessage(context.Background(), Request{
		To:      "+905551112233",
		Content: "Test message",
	})

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error sending request")
}
