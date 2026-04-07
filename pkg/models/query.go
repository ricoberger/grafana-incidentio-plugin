package models

type QueryType string

const (
	QueryTypeAttributes      = "attributes"
	QueryTypeAttributeValues = "attributevalues"
	QueryTypeAlerts          = "alerts"
	QueryTypeIncidents       = "incidents"
)

type QueryModeAttribute struct {
	Type string `json:"type"`
}

type QueryModeAttributeValues struct {
	Type      string `json:"type"`
	Attribute string `json:"attribute"`
}

type QueryModelAlerts struct {
	Filters []Filter `json:"filters"`
	Limit   int      `json:"limit"`
}

type Filter struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type QueryModelIncidents struct{}
