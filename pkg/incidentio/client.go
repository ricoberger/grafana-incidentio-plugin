package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ricoberger/grafana-incidentio-plugin/pkg/models"
	"github.com/ricoberger/grafana-incidentio-plugin/pkg/roundtripper"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Client interface {
	GetAlertAttributes(ctx context.Context) ([]models.AlertAttribute, error)
	GetAlertAttributeValues(ctx context.Context, attribute string) ([]models.CatalogEntry, error)
	GetAlerts(ctx context.Context, from, to time.Time, filters []models.Filter, limit int) ([]models.Alert, error)
	GetIncidentAttributes(ctx context.Context) ([]models.AlertAttribute, error)
	GetIncidentAttributeValues(ctx context.Context, attribute string) ([]models.CatalogEntry, error)
	GetIncidents(ctx context.Context, from, to time.Time, filters []models.Filter, limit int) ([]models.Incident, error)
}

type client struct {
	logger     log.Logger
	httpClient *http.Client
}

func (c *client) getAlertAttributes(ctx context.Context) ([]models.AlertAttribute, error) {
	return []models.AlertAttribute{
		{ID: "status", Name: "Alert Status"},
		{ID: "alert_source", Name: "Alert Source"},
	}, nil

	// The incident.io API does not support filtering alerts by custom
	// attributes, so we currently return only the default attributes "status"
	// and "alert_source". If the API adds support for filtering by custom
	// attributes in the future, we can uncomment the code below to fetch the
	// custom attributes from the API and return them as well.
	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v2/alert_attributes", nil)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// resp, err := c.httpClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	c.logger.Error("failed to get alert attributes", "status_code", resp.StatusCode, "response_body", string(body))
	// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }
	//
	// var alertAttributesResponse models.AlertAttributesResponse
	// if err := json.NewDecoder(resp.Body).Decode(&alertAttributesResponse); err != nil {
	// 	return nil, fmt.Errorf("failed to decode response: %w", err)
	// }
	//
	// return append(
	// 	[]models.AlertAttribute{
	// 		{ID: "status", Name: "Alert Status"},
	// 		{ID: "alert_source", Name: "Alert Source"},
	// 	},
	// 	alertAttributesResponse.AlertAttributes...,
	// ), nil
}

func (c *client) getAlertSources(ctx context.Context) ([]models.AlertSource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v2/alert_sources", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("failed to get alert attributes", "status_code", resp.StatusCode, "response_body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var alertSourcesResponse models.AlertSourcesResponse
	if err := json.NewDecoder(resp.Body).Decode(&alertSourcesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return alertSourcesResponse.AlertSources, nil
}

func (c *client) getCatalogTypes(ctx context.Context) ([]models.CatalogType, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v3/catalog_types", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("failed to get catalog types", "status_code", resp.StatusCode, "response_body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var catalogTypesResponse models.CatalogTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&catalogTypesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return catalogTypesResponse.CatalogTypes, nil
}

func (c *client) getCatalogEntries(ctx context.Context, catalogTypeId string) ([]models.CatalogEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.incident.io/v3/catalog_entries?page_size=250&catalog_type_id=%s", catalogTypeId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("failed to get catalog entries", "status_code", resp.StatusCode, "response_body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var catalogEntriesResponse models.CatalogEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&catalogEntriesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return catalogEntriesResponse.CatalogEntries, nil
}

func (c *client) getCustomFields(ctx context.Context) ([]models.CustomField, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v2/custom_fields", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("failed to get custom fields", "status_code", resp.StatusCode, "response_body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var customFieldsResponse models.IncidentCustomFieldsResponse
	if err := json.NewDecoder(resp.Body).Decode(&customFieldsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return customFieldsResponse.IncidentCustomFields, nil
}

func (c *client) getCustomFieldOptions(ctx context.Context, customFieldID string) ([]models.CatalogEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.incident.io/v1/custom_field_options?page_size=250&custom_field_id=%s", customFieldID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("failed to get custom field options", "status_code", resp.StatusCode, "response_body", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var customFieldOptionsResponse models.IncidentCustomFieldOptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&customFieldOptionsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var options []models.CatalogEntry
	for _, option := range customFieldOptionsResponse.IncidentCustomFieldsptions {
		options = append(options, models.CatalogEntry{
			ID:   option.ID,
			Name: option.Value,
		})
	}

	return options, nil
}

func (c *client) GetAlertAttributes(ctx context.Context) ([]models.AlertAttribute, error) {
	return c.getAlertAttributes(ctx)
}

func (c *client) GetAlertAttributeValues(ctx context.Context, attribute string) ([]models.CatalogEntry, error) {
	if attribute == "status" {
		return []models.CatalogEntry{
			{ID: "firing", Name: "Firing"},
			{ID: "resolved", Name: "Resolved"},
		}, nil
	}

	if attribute == "alert_source" {
		alertSources, err := c.getAlertSources(ctx)
		if err != nil {
			return nil, err
		}

		var catalogEntries []models.CatalogEntry
		for _, alertSource := range alertSources {
			catalogEntries = append(catalogEntries, models.CatalogEntry(alertSource))
		}
		return catalogEntries, nil
	}

	alertAttributes, err := c.getAlertAttributes(ctx)
	if err != nil {
		return nil, err
	}

	alertAttribute := getAlertAttribute(attribute, alertAttributes)
	if alertAttribute == nil {
		return nil, fmt.Errorf("alert attribute not found: %s", attribute)
	}

	catalogTypes, err := c.getCatalogTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog types: %w", err)
	}

	// If the catalog type is not found, we do not return an error, because it
	// might happen that the attribute does not have a catalog type. In this
	// case we show an text input field in the frontend instead of a dropdown
	// with options.
	catalogType := getCatalogType(alertAttribute, catalogTypes)
	if catalogType == nil {
		return nil, nil
	}

	return c.getCatalogEntries(ctx, catalogType.ID)
}

func (c *client) GetAlerts(ctx context.Context, from, to time.Time, filters []models.Filter, limit int) ([]models.Alert, error) {
	var after string
	var alerts []models.Alert

	for {
		params := url.Values{}
		params.Add("created_at[date_range]", fmt.Sprintf("%s~%s", from.Format("2006-01-02"), to.Format("2006-01-02")))

		if limit == 0 {
			params.Add("page_size", "50")
		} else if limit <= 50 {
			params.Add("page_size", fmt.Sprintf("%d", limit))
		} else {
			params.Add("page_size", "50")
			params.Add("after", after)
		}

		for _, filter := range filters {
			for value := range strings.SplitSeq(filter.Value, ",") {
				switch filter.Attribute {
				case "status":
					params.Add("status[one_of]", value)
				case "alert_source":
					params.Add("alert_source[one_of]", value)
				default:
					params.Add(fmt.Sprintf("attributes[%s][%s]", filter.Attribute, filter.Operator), value)
				}
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.incident.io/v2/alerts?%s", params.Encode()), nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.logger.Error("failed to get alerts", "status_code", resp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var alertsResponse models.AlertsResponse
		if err := json.NewDecoder(resp.Body).Decode(&alertsResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		alerts = append(alerts, alertsResponse.Alerts...)

		if len(alerts) >= limit || len(alertsResponse.Alerts) < 50 || alertsResponse.PaginationMeta.After == "" {
			break
		}
	}

	return alerts, nil
}

func (c *client) GetIncidentAttributes(ctx context.Context) ([]models.AlertAttribute, error) {
	customFields, err := c.getCustomFields(ctx)
	if err != nil {
		return nil, err
	}

	var attributes []models.AlertAttribute
	for _, customField := range customFields {
		attributes = append(attributes, models.AlertAttribute{
			ID:   customField.ID,
			Name: customField.Name,
		})
	}

	return append([]models.AlertAttribute{
		{ID: "status", Name: "Status"},
		{ID: "status_category", Name: "Status Category"},
		{ID: "severity", Name: "Severity"},
		{ID: "incident_type", Name: "Incident Type"},
	}, attributes...), nil
}

func (c *client) GetIncidentAttributeValues(ctx context.Context, attribute string) ([]models.CatalogEntry, error) {
	if attribute == "status" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v1/incident_statuses", nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.logger.Error("failed to get incident statuses", "status_code", resp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var statusResponse models.IncidentStatusesResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return statusResponse.IncidentStatuses, nil
	}

	if attribute == "status_category" {
		return []models.CatalogEntry{
			{ID: "triage", Name: "Triage"},
			{ID: "declined", Name: "Declined"},
			{ID: "merged", Name: "Merged"},
			{ID: "canceled", Name: "Canceled"},
			{ID: "live", Name: "Live"},
			{ID: "learning", Name: "Learning"},
			{ID: "closed", Name: "Closed"},
		}, nil
	}

	if attribute == "severity" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v1/severities", nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.logger.Error("failed to get incident statuses", "status_code", resp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var statusResponse models.IncidentSeveritiesResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return statusResponse.IncidentSeverities, nil
	}

	if attribute == "incident_type" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.incident.io/v1/incident_types", nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.logger.Error("failed to get incident statuses", "status_code", resp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var typesResponse models.IncidentTypesResponse
		if err := json.NewDecoder(resp.Body).Decode(&typesResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return typesResponse.IncidentTypes, nil
	}

	customFields, err := c.getCustomFields(ctx)
	if err != nil {
		return nil, err
	}

	customField := getCustomField(attribute, customFields)
	if customField == nil {
		return nil, fmt.Errorf("custom field not found: %s", attribute)
	}

	if customField.CatalogTypeID != "" {
		return c.getCatalogEntries(ctx, customField.CatalogTypeID)
	}

	return c.getCustomFieldOptions(ctx, attribute)
}

func (c *client) GetIncidents(ctx context.Context, from, to time.Time, filters []models.Filter, limit int) ([]models.Incident, error) {
	var after string
	var incidents []models.Incident

	for {
		params := url.Values{}
		params.Add("created_at[date_range]", fmt.Sprintf("%s~%s", from.Format("2006-01-02"), to.Format("2006-01-02")))

		if limit == 0 {
			params.Add("page_size", "50")
		} else if limit <= 50 {
			params.Add("page_size", fmt.Sprintf("%d", limit))
		} else {
			params.Add("page_size", "50")
			params.Add("after", after)
		}

		for _, filter := range filters {
			for value := range strings.SplitSeq(filter.Value, ",") {
				if filter.Attribute == "status" || filter.Attribute == "status_category" || filter.Attribute == "severity" || filter.Attribute == "incident_type" {
					params.Add(fmt.Sprintf("%s[%s]", filter.Attribute, filter.Operator), value)
				} else {
					params.Add(fmt.Sprintf("custom_field[%s][%s]", filter.Attribute, filter.Operator), value)
				}
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.incident.io/v2/incidents?%s", params.Encode()), nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.logger.Error("failed to get alerts", "status_code", resp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var incidentsResponse models.IncidentsResponse
		if err := json.NewDecoder(resp.Body).Decode(&incidentsResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		incidents = append(incidents, incidentsResponse.Incidents...)

		if len(incidents) >= limit || len(incidentsResponse.Incidents) < 50 || incidentsResponse.PaginationMeta.After == "" {
			break
		}
	}

	return incidents, nil
}

func NewClient(logger log.Logger, apiKey string) (Client, error) {
	return &client{
		logger: logger,
		httpClient: &http.Client{
			Transport: roundtripper.TokenAuthTransporter{
				Transport: roundtripper.DefaultRoundTripper,
				Token:     apiKey,
			},
		},
	}, nil
}
