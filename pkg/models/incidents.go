package models

import (
	"time"
)

type IncidentStatusesResponse struct {
	IncidentStatuses []CatalogEntry `json:"incident_statuses"`
}

type IncidentTypesResponse struct {
	IncidentTypes []CatalogEntry `json:"incident_types"`
}

type IncidentSeveritiesResponse struct {
	IncidentSeverities []CatalogEntry `json:"severities"`
}

type CustomField struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CatalogTypeID string `json:"catalog_type_id"`
}

type IncidentCustomFieldsResponse struct {
	IncidentCustomFields []CustomField `json:"custom_fields"`
}

type CustomFieldOption struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type IncidentCustomFieldOptionsResponse struct {
	IncidentCustomFieldsptions []CustomFieldOption `json:"custom_field_options"`
}

type Incident struct {
	CallURL   string    `json:"call_url"`
	CreatedAt time.Time `json:"created_at"`
	Creator   struct {
		Alert struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"alert"`
		APIKey struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"api_key"`
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Role        string `json:"role"`
			SlackUserID string `json:"slack_user_id"`
		} `json:"user"`
		Workflow struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"workflow"`
	} `json:"creator"`
	CustomFieldEntries []struct {
		CustomField struct {
			Description string `json:"description"`
			FieldType   string `json:"field_type"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Options     []struct {
				CustomFieldID string `json:"custom_field_id"`
				ID            string `json:"id"`
				SortKey       int64  `json:"sort_key"`
				Value         string `json:"value"`
			} `json:"options"`
		} `json:"custom_field"`
		Values []struct {
			ValueCatalogEntry struct {
				Aliases    []string `json:"aliases"`
				ExternalID string   `json:"external_id"`
				ID         string   `json:"id"`
				Name       string   `json:"name"`
			} `json:"value_catalog_entry"`
			ValueLink    string `json:"value_link"`
			ValueNumeric string `json:"value_numeric"`
			ValueOption  struct {
				CustomFieldID string `json:"custom_field_id"`
				ID            string `json:"id"`
				SortKey       int64  `json:"sort_key"`
				Value         string `json:"value"`
			} `json:"value_option"`
			ValueText string `json:"value_text"`
		} `json:"values"`
	} `json:"custom_field_entries"`
	DurationMetrics []struct {
		DurationMetric struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"duration_metric"`
		ValueSeconds int64 `json:"value_seconds"`
	} `json:"duration_metrics"`
	ExternalIssueReference struct {
		IssueName      string `json:"issue_name"`
		IssuePermalink string `json:"issue_permalink"`
		Provider       string `json:"provider"`
	} `json:"external_issue_reference"`
	HasDebrief              bool   `json:"has_debrief"`
	ID                      string `json:"id"`
	IncidentRoleAssignments []struct {
		Assignee struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Role        string `json:"role"`
			SlackUserID string `json:"slack_user_id"`
		} `json:"assignee"`
		Role struct {
			CreatedAt    time.Time `json:"created_at"`
			Description  string    `json:"description"`
			ID           string    `json:"id"`
			Instructions string    `json:"instructions"`
			Name         string    `json:"name"`
			Required     bool      `json:"required"`
			RoleType     string    `json:"role_type"`
			Shortform    string    `json:"shortform"`
			UpdatedAt    time.Time `json:"updated_at"`
		} `json:"role"`
	} `json:"incident_role_assignments"`
	IncidentStatus struct {
		Category    string    `json:"category"`
		CreatedAt   time.Time `json:"created_at"`
		Description string    `json:"description"`
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Rank        int64     `json:"rank"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"incident_status"`
	IncidentTimestampValues []struct {
		IncidentTimestamp struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Rank int64  `json:"rank"`
		} `json:"incident_timestamp"`
		Value struct {
			Value time.Time `json:"value"`
		} `json:"value"`
	} `json:"incident_timestamp_values"`
	IncidentType struct {
		CreateInTriage       string    `json:"create_in_triage"`
		CreatedAt            time.Time `json:"created_at"`
		Description          string    `json:"description"`
		ID                   string    `json:"id"`
		IsDefault            bool      `json:"is_default"`
		Name                 string    `json:"name"`
		PrivateIncidentsOnly bool      `json:"private_incidents_only"`
		UpdatedAt            time.Time `json:"updated_at"`
	} `json:"incident_type"`
	Mode                  string   `json:"mode"`
	Name                  string   `json:"name"`
	Permalink             string   `json:"permalink"`
	PostmortemDocumentIds []string `json:"postmortem_document_ids"`
	PostmortemDocumentURL string   `json:"postmortem_document_url"`
	Reference             string   `json:"reference"`
	Severity              struct {
		CreatedAt   time.Time `json:"created_at"`
		Description string    `json:"description"`
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Rank        int64     `json:"rank"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"severity"`
	SlackChannelID          string    `json:"slack_channel_id"`
	SlackChannelName        string    `json:"slack_channel_name"`
	SlackTeamID             string    `json:"slack_team_id"`
	Summary                 string    `json:"summary"`
	UpdatedAt               time.Time `json:"updated_at"`
	Visibility              string    `json:"visibility"`
	WorkloadMinutesLate     float64   `json:"workload_minutes_late"`
	WorkloadMinutesSleeping float64   `json:"workload_minutes_sleeping"`
	WorkloadMinutesTotal    float64   `json:"workload_minutes_total"`
	WorkloadMinutesWorking  float64   `json:"workload_minutes_working"`
}

type IncidentsResponse struct {
	Incidents      []Incident `json:"incidents"`
	PaginationMeta struct {
		After            string `json:"after"`
		PageSize         int64  `json:"page_size"`
		TotalRecordCount int64  `json:"total_record_count"`
	} `json:"pagination_meta"`
}
