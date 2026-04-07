package models

import (
	"time"
)

type AlertAttribute struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AlertAttributesResponse struct {
	AlertAttributes []AlertAttribute `json:"alert_attributes"`
}

type AlertSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AlertSourcesResponse struct {
	AlertSources []AlertSource `json:"alert_sources"`
}

type Alert struct {
	AlertSourceID string `json:"alert_source_id"`
	Attributes    []struct {
		ArrayValue []struct {
			CatalogEntry struct {
				CatalogTypeID string `json:"catalog_type_id"`
				ID            string `json:"id"`
				Name          string `json:"name"`
			} `json:"catalog_entry"`
			Label   string `json:"label"`
			Literal string `json:"literal"`
		} `json:"array_value"`
		Attribute struct {
			Array    bool   `json:"array"`
			Emoji    string `json:"emoji"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Required bool   `json:"required"`
			Type     string `json:"type"`
		} `json:"attribute"`
		Value struct {
			CatalogEntry struct {
				CatalogTypeID string `json:"catalog_type_id"`
				ID            string `json:"id"`
				Name          string `json:"name"`
			} `json:"catalog_entry"`
			Label   string `json:"label"`
			Literal string `json:"literal"`
		} `json:"value"`
	} `json:"attributes"`
	CreatedAt        time.Time  `json:"created_at"`
	DeduplicationKey string     `json:"deduplication_key"`
	Description      string     `json:"description"`
	ID               string     `json:"id"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	SourceURL        string     `json:"source_url"`
	Status           string     `json:"status"`
	Title            string     `json:"title"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type AlertsResponse struct {
	Alerts         []Alert `json:"alerts"`
	PaginationMeta struct {
		After    string `json:"after"`
		PageSize int    `json:"page_size"`
	} `json:"pagination_meta"`
}
