package beem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Beem struct {
	ApiKey    string
	SecretKey string
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var BeemURL = "https://apisms.beem.africa/v1/send"

type Recipient struct {
	RecipientID int    `json:"recipient_id"`
	DestAddr    string `json:"dest_addr"`
}

type successResponse struct {
	Successful bool   `json:"successful"`
	RequestID  int    `json:"request_id"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Valid      int    `json:"valid"`
	Invalid    int    `json:"invalid"`
	Duplicates int    `json:"duplicates"`
}

func New(apiKey string, secretKey string) *Beem {
	return &Beem{
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}
}

func (b *Beem) Send(msg string, recipient string) error {

	phone := "255" + recipient[1:]

	apiKey := b.ApiKey
	secretKey := b.SecretKey
	sourceAddr := "Mangi"

	recipients := []Recipient{
		{
			RecipientID: 1,
			DestAddr:    phone,
		},
	}

	payload := map[string]interface{}{
		"source_addr": sourceAddr,
		"encoding":    "0",
		"message":     msg,
		"recipients":  recipients,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {

		slog.Error(fmt.Sprintf("Error encoding JSON payload: %s", err))

		return err
	}

	req, err := http.NewRequest("POST", BeemURL, bytes.NewBuffer(jsonPayload))

	if err != nil {

		slog.Error(fmt.Sprintf("Error creating HTTP request: %s", err))

		return err
	}

	req.SetBasicAuth(apiKey, secretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {

		slog.Error(fmt.Sprintf("Error sending request: %s", err))

		return err
	}

	defer res.Body.Close()

	c := context.Background()

	if res.StatusCode != http.StatusOK {
		var errResponse errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errResponse); err != nil {
			slog.Error(fmt.Sprintf("Error decoding  error response: %s", err))
			return err
		}

		slog.LogAttrs(c,
			slog.LevelError,
			"Error response from Beem",

			slog.Int("http-status", res.StatusCode),
			slog.Int("beem-code", errResponse.Code),
			slog.String("message", errResponse.Message),
		)

		return errors.New("error response from Beem")

	}

	var successResponse successResponse

	if err = json.NewDecoder(res.Body).Decode(&successResponse); err != nil {
		slog.Error(fmt.Sprintf("Error decoding success response: %s", err))
		return err
	}

	slog.LogAttrs(c,
		slog.LevelInfo,
		"Success response from Beem",
		slog.Int("http-status", res.StatusCode),
		slog.Bool("successful", successResponse.Successful),
		slog.Int("request-id", successResponse.RequestID),
		slog.Int("beem-code", successResponse.Code),
		slog.String("message", successResponse.Message),
		slog.Int("valid", successResponse.Valid),
		slog.Int("invalid", successResponse.Invalid),
		slog.Int("duplicates", successResponse.Duplicates),
	)
	return nil
}
