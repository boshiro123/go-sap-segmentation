package sap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go-test/model"
	"go-test/pkg/config"
)

// Client представляет клиент для работы с SAP API
type Client struct {
	httpClient *http.Client
	baseURL    string
	authHeader string
	userAgent  string
	batchSize  int
	interval   time.Duration
	logger     *slog.Logger
}

// Response представляет ответ от SAP API
type Response struct {
	Items []model.Segmentation `json:"items"`
}

// NewClient создает новый клиент для работы с SAP API
func NewClient(cfg *config.Config, logger *slog.Logger) *Client {
	// Кодируем логин:пароль в base64 для Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(cfg.Connection.AuthLoginPwd))

	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Connection.Timeout,
		},
		baseURL:    cfg.Connection.URI,
		authHeader: fmt.Sprintf("Basic %s", auth),
		userAgent:  cfg.Connection.UserAgent,
		batchSize:  cfg.Import.BatchSize,
		interval:   cfg.Connection.Interval,
		logger:     logger,
	}
}

// FetchSegmentation получает данные о сегментации из SAP API
func (c *Client) FetchSegmentation() ([]*model.Segmentation, error) {
	var allSegments []*model.Segmentation
	offset := 0

	for {
		c.logger.Info("fetching data from SAP API",
			"url", c.baseURL,
			"offset", offset,
			"limit", c.batchSize,
		)

		// Формируем URL с параметрами
		url := fmt.Sprintf("%s?p_limit=%d&p_offset=%d", c.baseURL, c.batchSize, offset)

		// Создаем запрос
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		// Устанавливаем заголовки
		req.Header.Set("Authorization", c.authHeader)
		req.Header.Set("User-Agent", c.userAgent)

		// Выполняем запрос
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		// Закрываем тело ответа после обработки
		defer resp.Body.Close()

		// Проверяем код ответа
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("error response from SAP API: status=%d, body=%s",
				resp.StatusCode, string(bodyBytes))
		}

		// Читаем тело ответа
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		// Если ответ пустой, значит достигли конца данных
		if len(bodyBytes) == 0 || string(bodyBytes) == "[]" || string(bodyBytes) == "{}" {
			break
		}

		// Декодируем JSON
		var response Response
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		// Если нет элементов, значит достигли конца данных
		if len(response.Items) == 0 {
			break
		}

		// Преобразуем элементы ответа в указатели на структуры Segmentation
		segments := make([]*model.Segmentation, len(response.Items))
		for i, item := range response.Items {
			segments[i] = &model.Segmentation{
				AddressSapID: item.AddressSapID,
				AdrSegment:   item.AdrSegment,
				SegmentID:    item.SegmentID,
			}
		}

		// Добавляем полученные сегменты к общему списку
		allSegments = append(allSegments, segments...)

		// Увеличиваем смещение для следующего запроса
		offset += c.batchSize

		// Делаем паузу перед следующим запросом
		time.Sleep(c.interval)
	}

	c.logger.Info("finished fetching data from SAP API", "total_segments", len(allSegments))

	return allSegments, nil
}
