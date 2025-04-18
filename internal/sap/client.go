package sap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"go-test/internal/models"
	"go-test/pkg/config"
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	authHeader  string
	userAgent   string
	batchSize   int
	interval    time.Duration
	logger      *slog.Logger
	useTestData bool
}

type Response struct {
	Items []models.Segmentation `json:"items"`
}

func NewClient(cfg *config.Config, logger *slog.Logger) *Client {
	logger.Info("config:", "cfg", cfg)
	auth := base64.StdEncoding.EncodeToString([]byte(cfg.Connection.AuthLoginPwd))

	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Connection.Timeout,
		},
		baseURL:     cfg.Connection.URI,
		authHeader:  fmt.Sprintf("Basic %s", auth),
		userAgent:   cfg.Connection.UserAgent,
		batchSize:   cfg.Import.BatchSize,
		interval:    cfg.Connection.Interval,
		logger:      logger,
		useTestData: cfg.Import.UseTestData,
	}
}

// generateTestData создает тестовые данные для разработки и тестирования
func (c *Client) generateTestData() []*models.Segmentation {
	c.logger.Info("generating test data for development")

	// Инициализируем генератор случайных чисел
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Генерируем случайные данные для тестирования
	testData := make([]*models.Segmentation, 30)
	prefixes := []string{"SAP-", "SEG-", "ADR-"}
	segments := []string{"Premium", "Standard", "VIP", "Corporate", "SMB"}

	for i := 0; i < 30; i++ {
		prefix := prefixes[r.Intn(len(prefixes))]
		segmentIdx := r.Intn(len(segments))

		testData[i] = &models.Segmentation{
			AddressSapID: fmt.Sprintf("%s%03d", prefix, i+1),
			AdrSegment:   segments[segmentIdx],
			SegmentID:    int64(1000 + i),
		}
	}

	return testData
}

func (c *Client) FetchSegmentation() ([]*models.Segmentation, error) {
	if c.useTestData {
		c.logger.Info("using test data as configured by USE_TEST_DATA=true")
		return c.generateTestData(), nil
	}

	var allSegments []*models.Segmentation
	offset := 0

	c.logger.Info("testing connection to SAP API", "url", c.baseURL)
	testURL := fmt.Sprintf("%s?p_limit=1&p_offset=0", c.baseURL)
	testReq, err := http.NewRequest(http.MethodGet, testURL, nil)
	if err != nil {
		c.logger.Error("error creating test request", "error", err.Error())
		return nil, fmt.Errorf("error creating test request: %w", err)
	}

	testReq.Header.Set("Authorization", c.authHeader)
	testReq.Header.Set("User-Agent", c.userAgent)
	testResp, testErr := c.httpClient.Do(testReq)

	if testErr != nil {
		c.logger.Error("error connecting to SAP API", "error", testErr.Error())
		return nil, fmt.Errorf("error connecting to SAP API: %w", testErr)
	}

	if testResp != nil && testResp.Body != nil {
		defer testResp.Body.Close()
	}

	if testResp != nil && testResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(testResp.Body)
		c.logger.Error("SAP API returned error status",
			"status", testResp.StatusCode,
			"body", string(bodyBytes),
		)

		// Если получена ошибка 401 Unauthorized, возвращаем тестовые данные
		if testResp.StatusCode == http.StatusUnauthorized {
			c.logger.Warn("SAP API is not available, using test data",
				"error", nil,
				"status", testResp.StatusCode,
			)
			return c.generateTestData(), nil
		}

		return nil, fmt.Errorf("error response from SAP API: status=%d, body=%s",
			testResp.StatusCode, string(bodyBytes))
	}

	for {
		c.logger.Info("fetching data from SAP API",
			"url", c.baseURL,
			"offset", offset,
			"limit", c.batchSize,
		)

		url := fmt.Sprintf("%s?p_limit=%d&p_offset=%d", c.baseURL, c.batchSize, offset)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", c.authHeader)
		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)

			// Если получена ошибка 401 Unauthorized на этом этапе, также используем тестовые данные
			if resp.StatusCode == http.StatusUnauthorized {
				c.logger.Warn("SAP API is not available during fetching, using test data",
					"error", nil,
					"status", resp.StatusCode,
				)
				return c.generateTestData(), nil
			}

			return nil, fmt.Errorf("error response from SAP API: status=%d, body=%s",
				resp.StatusCode, string(bodyBytes))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if len(bodyBytes) == 0 || string(bodyBytes) == "[]" || string(bodyBytes) == "{}" {
			break
		}

		var response Response
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if len(response.Items) == 0 {
			break
		}

		segments := make([]*models.Segmentation, 0)
		for _, item := range response.Items {
			segments = append(segments, &models.Segmentation{
				AddressSapID: item.AddressSapID,
				AdrSegment:   item.AdrSegment,
				SegmentID:    item.SegmentID,
			})
		}

		allSegments = append(allSegments, segments...)

		offset += c.batchSize

		time.Sleep(c.interval)
	}

	c.logger.Info("finished fetching data from SAP API", "total_segments", len(allSegments))

	return allSegments, nil
}
