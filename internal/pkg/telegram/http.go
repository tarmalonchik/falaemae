package telegram

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

const (
	documentField = "document"
	photoField    = "photo"
	videoField    = "video"
	mediaField    = "media"
	chatIDField   = "chat_id"
	clientName    = "tg-bot-client"
	tokenTemp     = "/bot%s"

	mediaPhotoType = "photo"
)

func (c *Client) SendDocument(ctx context.Context, req SendDocumentRequest) error {
	if !req.ChatID.Valid {
		return nil
	}

	const (
		method        = "/sendDocument"
		requestMethod = http.MethodPost
	)

	u, err := url.Parse(c.conf.GetTgBotBaseURL())
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] parse url '%s' error %w",
			requestMethod, clientName, method, c.conf.GetTgBotBaseURL(), err)
	}
	u.Path += fmt.Sprintf(tokenTemp, c.conf.GetTgBotToken()) + method

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer func() { _ = writer.Close() }()

	fileID := c.resolveFileID(req.FilePath)

	if fileID == "" {
		part, err := writer.CreateFormFile(documentField, req.FileName)
		if err != nil {
			return err
		}
		if _, err := part.Write(req.File); err != nil {
			return err
		}
	} else {
		if err = writer.WriteField(documentField, fileID); err != nil {
			return err
		}
	}

	if err = writer.WriteField(chatIDField, strconv.Itoa(int(req.ChatID.Int64))); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), body)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] create request error: %w",
			requestMethod, clientName, method, err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	respHTTP, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] requesting error: %w",
			requestMethod, clientName, method, err)
	}
	defer func() { _ = respHTTP.Close }()

	if respHTTP.StatusCode >= 400 {
		return nil
	}

	return c.saveData(req.FilePath, respHTTP.Body)
}

func (c *Client) resolveFileID(path sql.NullString) string {
	if !path.Valid {
		return ""
	}
	if path.String == "" {
		return ""
	}
	id, ok := c.documentsMap.Get(path.String)
	if ok {
		return id
	}
	return ""
}

func (c *Client) saveData(filePath sql.NullString, reader io.Reader) error {
	if !filePath.Valid || filePath.String == "" {
		return nil
	}

	var responseData sendDocResponseBody

	respDataByte, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(respDataByte, &responseData); err != nil {
		return err
	}

	c.documentsMap.AddData(filePath.String, responseData.Result.Document.FileID)
	return nil
}

func (c *Client) SendVideoByID(ctx context.Context, chatID int64, fileID string) error {
	const (
		method        = "/sendVideo"
		requestMethod = http.MethodPost
	)

	u, err := url.Parse(c.conf.GetTgBotBaseURL())
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] parse url '%s' error %w",
			requestMethod, clientName, method, c.conf.GetTgBotBaseURL(), err)
	}
	u.Path += fmt.Sprintf(tokenTemp, c.conf.GetTgBotToken()) + method

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer func() { _ = writer.Close() }()

	if err = writer.WriteField(chatIDField, strconv.Itoa(int(chatID))); err != nil {
		return err
	}
	if err = writer.WriteField(videoField, fileID); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), body)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] create request error: %w",
			requestMethod, clientName, method, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	respHTTP, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] requesting error: %w",
			requestMethod, clientName, method, err)
	}
	defer func() { _ = respHTTP.Close }()

	if respHTTP.StatusCode >= 400 {
		return fmt.Errorf("[%s][client][%s][%s] bad status code: %d",
			requestMethod, clientName, method, respHTTP.StatusCode)
	}
	return nil
}

func (c *Client) SendPhotos(ctx context.Context, chatID int64, files ...string) error {
	if len(files) == 1 {
		return c.sendPhotoByID(ctx, chatID, files[0])
	}

	filesBatch := make([]struct {
		Type  string `json:"type"`
		Media string `json:"media"`
	}, len(files))

	for i := range files {
		filesBatch[i].Media = files[i]
		filesBatch[i].Type = mediaPhotoType
	}

	filesBatchBytes, err := json.Marshal(filesBatch)
	if err != nil {
		return fmt.Errorf("telegram SendPhotosGroup error marshal: %w", err)
	}
	return c.sendMediaGroup(ctx, chatID, string(filesBatchBytes))
}

func (c *Client) sendMediaGroup(ctx context.Context, chatID int64, filesBatchJSON string) error {
	const (
		method        = "/sendMediaGroup"
		requestMethod = http.MethodPost
	)

	u, err := url.Parse(c.conf.GetTgBotBaseURL())
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] parse url '%s' error %w",
			requestMethod, clientName, method, c.conf.GetTgBotBaseURL(), err)
	}
	u.Path += fmt.Sprintf(tokenTemp, c.conf.GetTgBotToken()) + method

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer func() { _ = writer.Close() }()

	if err = writer.WriteField(chatIDField, strconv.Itoa(int(chatID))); err != nil {
		return err
	}
	if err = writer.WriteField(mediaField, filesBatchJSON); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), body)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] create request error: %w",
			requestMethod, clientName, method, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	respHTTP, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] requesting error: %w",
			requestMethod, clientName, method, err)
	}
	defer func() { _ = respHTTP.Close }()

	if respHTTP.StatusCode >= 400 {
		return fmt.Errorf("[%s][client][%s][%s] bad status code: %d",
			requestMethod, clientName, method, respHTTP.StatusCode)
	}
	return nil
}

func (c *Client) sendPhotoByID(ctx context.Context, chatID int64, fileID string) error {
	const (
		method        = "/sendPhoto"
		requestMethod = http.MethodPost
	)

	u, err := url.Parse(c.conf.GetTgBotBaseURL())
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] parse url '%s' error %w",
			requestMethod, clientName, method, c.conf.GetTgBotBaseURL(), err)
	}
	u.Path += fmt.Sprintf(tokenTemp, c.conf.GetTgBotToken()) + method

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer func() { _ = writer.Close() }()

	if err = writer.WriteField(chatIDField, strconv.Itoa(int(chatID))); err != nil {
		return err
	}
	if err = writer.WriteField(photoField, fileID); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, requestMethod, u.String(), body)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] create request error: %w",
			requestMethod, clientName, method, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	respHTTP, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[%s][client][%s][%s] requesting error: %w",
			requestMethod, clientName, method, err)
	}
	defer func() { _ = respHTTP.Close }()

	if respHTTP.StatusCode >= 400 {
		return fmt.Errorf("[%s][client][%s][%s] bad status code: %d",
			requestMethod, clientName, method, respHTTP.StatusCode)
	}
	return nil
}
