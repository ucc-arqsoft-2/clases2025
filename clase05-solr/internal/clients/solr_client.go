package clients

import (
	"clase05-solr/internal/domain"
	"context"
	"encoding/json"
	"fmt"
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
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
		Name:      item.Name,
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
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
		return domain.PaginatedResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	items := make([]domain.Item, len(solrResp.Response.Docs))
	for i, doc := range solrResp.Response.Docs {
		items[i] = domain.Item{
			ID:        doc.ID,
			Name:      doc.Name,
			Price:     doc.Price,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
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
