package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"billing/internal/apperrors"
)

type StockClient struct {
	baseURL string
	http    *http.Client
}

func NewStockClient(baseURL string) *StockClient {
	return &StockClient{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

type ProductResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int64     `json:"quantity"`
	Price    int64     `json:"price"`
}

func (c *StockClient) GetProduct(ctx context.Context, productID uuid.UUID) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/products/%s", c.baseURL, productID)
	return c.doGet(ctx, url)
}

func (c *StockClient) ReduceStock(ctx context.Context, productID uuid.UUID, quantity int64) error {
	url := fmt.Sprintf("%s/products/%s/reduce-stock", c.baseURL, productID)
	body := map[string]int64{"quantity": quantity}
	return c.doPatch(ctx, url, body)
}

func (c *StockClient) RestoreStock(ctx context.Context, productID uuid.UUID, quantity int64) error {
	url := fmt.Sprintf("%s/products/%s/restore-stock", c.baseURL, productID)
	body := map[string]int64{"quantity": quantity}
	return c.doPatch(ctx, url, body)
}

func (c *StockClient) doGet(ctx context.Context, url string) (*ProductResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, apperrors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, apperrors.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, apperrors.ErrServiceUnavailable
	}

	var product ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (c *StockClient) doPatch(ctx context.Context, url string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return apperrors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperrors.ErrServiceUnavailable
	}
	return nil
}
