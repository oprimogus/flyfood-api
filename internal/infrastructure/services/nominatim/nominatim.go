package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oprimogus/flyfood-api/internal/config"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	conf   *config.Nominatim
	client *http.Client
)

func init() {
	conf = config.GetInstance().Nominatim
	client = &http.Client{Timeout: 15 * time.Second}
}

type Query struct {
	Q          string `json:"q"`
	Amenity    string `json:"amenity"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	County     string `json:"county"`
	Country    string `json:"country"`
	PostalCode string `json:"postalcode"`
	Format     string `json:"format"`
	Limit      int    `json:"limit"`
}

type Location struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OSMType     string   `json:"osm_type"`
	OSMID       int      `json:"osm_id"`
	Latitude    string   `json:"lat"`
	Longitude   string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	AddressType string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
}

func (q Query) toQueryParams() string {
	values := url.Values{}

	values.Add("format", "json")
	values.Add("limit", strconv.Itoa(max(1, q.Limit)))

	if q.Q != "" {
		values.Add("q", q.Q)
	} else {
		if q.Amenity != "" {
			values.Add("amenity", q.Amenity)
		}
		if q.Street != "" {
			values.Add("street", q.Street)
		}
		if q.City != "" {
			values.Add("city", q.City)
		}
		if q.State != "" {
			values.Add("state", q.State)
		}
		if q.County != "" {
			values.Add("county", q.County)
		}
		if q.Country != "" {
			values.Add("country", q.Country)
		}
		if q.PostalCode != "" {
			values.Add("postalcode", q.PostalCode)
		}
	}

	return values.Encode()
}

func Search(ctx context.Context, parameters Query) ([]Location, error) {
	queryParams := parameters.toQueryParams()
	fullURL := conf.BaseURL + "/search?" + queryParams

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request failed", "url", fullURL, "error", err)
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("error closing response body", "error", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		slog.Warn("Nominatim API returned error", "status", resp.StatusCode, "response", string(body))
		return nil, fmt.Errorf("nominatim API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result []Location
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return result, nil
}
