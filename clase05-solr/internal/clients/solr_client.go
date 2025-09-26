package clients

import (
	"clase05-solr/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)


type SolrClient struct {
	baseURL string
	core    string
	client  *http.Client
}

type SolrDocument struct {
	ID        string    `json:"id"`
	Name      []string  `json:"name"`
	Price     []float64 `json:"price"`
	CreatedAt []string  `json:"created_at"`
	UpdatedAt []string  `json:"updated_at"`
}

type SolrResponse struct {
	Response struct {
		NumFound int            `json:"numFound"`
		Start    int            `json:"start"`
		Docs     []SolrDocument `json:"docs"`
	} `json:"response"`
}

type SolrUpdateResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
}

const (
	defaultCount = 10
)

func NewSolrClient(host, port, core string) *SolrClient {
	baseURL := fmt.Sprintf("http://%s:%s/solr/%s", host, port, core)
	return &SolrClient{
		baseURL: baseURL,
		core:    core,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SolrClient) Index(ctx context.Context, item domain.Item) error {
	doc := SolrDocument{
		ID:        item.ID,
		Name:      []string{item.Name},
		Price:     []float64{item.Price},
		CreatedAt: []string{item.CreatedAt.Format(time.RFC3339)},
		UpdatedAt: []string{item.UpdatedAt.Format(time.RFC3339)},
	}

	data, err := json.Marshal([]SolrDocument{doc})
	if err != nil {
		return fmt.Errorf("error marshalling document: %w", err)
	}

	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		return fmt.Errorf("solr update failed with status %d", updateResp.ResponseHeader.Status)
	}

	return nil
}

func (s *SolrClient) Search(ctx context.Context, query string, page int, count int) (domain.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if count <= 0 {
		count = defaultCount
	}

	// calcular offset
	start := (page - 1) * count

	params := url.Values{}
	params.Set("q", query)
	params.Set("wt", "json")
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("rows", fmt.Sprintf("%d", count))

	url := fmt.Sprintf("%s/select?%s", s.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.PaginatedResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.PaginatedResponse{}, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.PaginatedResponse{}, fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		// Log response body for debugging as json
		slog.Info("Response body for debugging", "body", json.NewDecoder(resp.Body))
		return domain.PaginatedResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	items := make([]domain.Item, len(solrResp.Response.Docs))
	for i, doc := range solrResp.Response.Docs {
		var name string
		if len(doc.Name) > 0 {
			name = doc.Name[0]
		}

		var price float64
		if len(doc.Price) > 0 {
			price = doc.Price[0]
		}

		var createdAt, updatedAt time.Time
		if len(doc.CreatedAt) > 0 {
			if t, err := time.Parse(time.RFC3339, doc.CreatedAt[0]); err == nil {
				createdAt = t
			}
		}
		if len(doc.UpdatedAt) > 0 {
			if t, err := time.Parse(time.RFC3339, doc.UpdatedAt[0]); err == nil {
				updatedAt = t
			}
		}

		items[i] = domain.Item{
			ID:        doc.ID,
			Name:      name,
			Price:     price,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return domain.PaginatedResponse{
		Page:    page,
		Count:   len(items),
		Total:   solrResp.Response.NumFound, // total de coincidencias
		Results: items,
	}, nil
}

func (s *SolrClient) Delete(ctx context.Context, id string) error {
	data := fmt.Sprintf(`{"delete":{"id":"%s"}}`, id)
	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		return fmt.Errorf("solr delete failed with status %d", updateResp.ResponseHeader.Status)
	}

	return nil
}
