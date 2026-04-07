package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ricoberger/grafana-incidentio-plugin/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/tracing"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/concurrent"
	"go.opentelemetry.io/otel/codes"
)

func (d *Datasource) handleAttributesQueries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAttributesQueries")
	defer span.End()

	return concurrent.QueryData(ctx, req, d.handleAttributes, 10)
}

func (d *Datasource) handleAttributes(ctx context.Context, query concurrent.Query) backend.DataResponse {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAttributes")
	defer span.End()

	var qm models.QueryModeAttributeValues
	err := json.Unmarshal(query.DataQuery.JSON, &qm)
	if err != nil {
		d.logger.Error("Failed to unmarshal query model", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var attributes []models.AlertAttribute

	switch qm.Type {
	case models.QueryTypeAlerts:
		tmpAttributes, err := d.incidentioClient.GetAlertAttributes(ctx)
		if err != nil {
			d.logger.Error("Failed to get alert attributes", "error", err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return backend.ErrorResponseWithErrorSource(err)
		}
		attributes = tmpAttributes
	case models.QueryTypeIncidents:
		tmpAttributes, err := d.incidentioClient.GetIncidentAttributes(ctx)
		if err != nil {
			d.logger.Error("Failed to get incident attributes", "error", err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return backend.ErrorResponseWithErrorSource(err)
		}
		attributes = tmpAttributes
	default:
		err = fmt.Errorf("unsupported query type: %s", qm.Type)
		d.logger.Error("Unsupported query type", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var ids []string
	var names []string

	for _, attribute := range attributes {
		ids = append(ids, attribute.ID)
		names = append(names, attribute.Name)
	}

	frame := data.NewFrame(
		"Attributes",
		data.NewField("ids", nil, ids),
		data.NewField("names", nil, names),
	)

	frame.SetMeta(&data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
		Type:                   data.FrameTypeTable,
	})

	var response backend.DataResponse
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *Datasource) handleAttributeValuesQueries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAttributeValuesQueries")
	defer span.End()

	return concurrent.QueryData(ctx, req, d.handleAttributeValues, 10)
}

func (d *Datasource) handleAttributeValues(ctx context.Context, query concurrent.Query) backend.DataResponse {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAttributeValues")
	defer span.End()

	var qm models.QueryModeAttributeValues
	err := json.Unmarshal(query.DataQuery.JSON, &qm)
	if err != nil {
		d.logger.Error("Failed to unmarshal query model", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var values []models.CatalogEntry

	switch qm.Type {
	case models.QueryTypeAlerts:
		tmpValues, err := d.incidentioClient.GetAlertAttributeValues(ctx, qm.Attribute)
		if err != nil {
			d.logger.Error("Failed to get alert attribute values", "error", err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return backend.ErrorResponseWithErrorSource(err)
		}
		values = tmpValues
	case models.QueryTypeIncidents:
		tmpValues, err := d.incidentioClient.GetIncidentAttributeValues(ctx, qm.Attribute)
		if err != nil {
			d.logger.Error("Failed to get incident attribute values", "error", err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return backend.ErrorResponseWithErrorSource(err)
		}
		values = tmpValues
	default:
		err = fmt.Errorf("unsupported query type: %s", qm.Type)
		d.logger.Error("Unsupported query type", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var ids []string
	var names []string

	for _, value := range values {
		ids = append(ids, value.ID)
		names = append(names, value.Name)
	}

	frame := data.NewFrame(
		"Attributes",
		data.NewField("ids", nil, ids),
		data.NewField("names", nil, names),
	)

	frame.SetMeta(&data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
		Type:                   data.FrameTypeTable,
	})

	var response backend.DataResponse
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *Datasource) handleAlertsQueries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAlertsQueries")
	defer span.End()

	return concurrent.QueryData(ctx, req, d.handleAlerts, 10)
}

func (d *Datasource) handleAlerts(ctx context.Context, query concurrent.Query) backend.DataResponse {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleAlerts")
	defer span.End()

	var qm models.QueryModelAlerts
	err := json.Unmarshal(query.DataQuery.JSON, &qm)
	if err != nil {
		d.logger.Error("Failed to unmarshal query model", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	values, err := d.incidentioClient.GetAlerts(ctx, query.DataQuery.TimeRange.From, query.DataQuery.TimeRange.To, qm.Filters, qm.Limit)
	if err != nil {
		d.logger.Error("Failed to get alerts", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var createdAt []time.Time
	var updatedAt []*time.Time
	var resolvedAt []*time.Time
	var id []string
	var title []string
	var status []string
	var description []string

	attributes := make(map[string][]string)

	for _, value := range values {
		createdAt = append(createdAt, value.CreatedAt)
		updatedAt = append(updatedAt, value.UpdatedAt)
		resolvedAt = append(resolvedAt, value.ResolvedAt)
		id = append(id, value.ID)
		title = append(title, value.Title)
		status = append(status, value.Status)
		description = append(description, value.Description)

		for _, attribute := range value.Attributes {
			attributeName := attribute.Attribute.Name
			attributeValue := attribute.Value.Label
			attributes[attributeName] = append(attributes[attributeName], attributeValue)
		}
	}

	frame := data.NewFrame(
		"Alerts",
		data.NewField("Created At", nil, createdAt),
		data.NewField("Updated At", nil, updatedAt),
		data.NewField("Resolved At", nil, resolvedAt),
		data.NewField("ID", nil, id),
		data.NewField("Title", nil, title),
		data.NewField("Status", nil, status),
		data.NewField("Description", nil, description),
	)

	for attributeName, attributeValues := range attributes {
		frame.Fields = append(frame.Fields, data.NewField(attributeName, nil, attributeValues))
	}

	frame.SetMeta(&data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
		Type:                   data.FrameTypeTable,
	})

	var response backend.DataResponse
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *Datasource) handleIncidentsQueries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleIncidentsQueries")
	defer span.End()

	return concurrent.QueryData(ctx, req, d.handleIncidents, 10)
}

func (d *Datasource) handleIncidents(ctx context.Context, query concurrent.Query) backend.DataResponse {
	ctx, span := tracing.DefaultTracer().Start(ctx, "handleIncidents")
	defer span.End()

	var qm models.QueryModelAlerts
	err := json.Unmarshal(query.DataQuery.JSON, &qm)
	if err != nil {
		d.logger.Error("Failed to unmarshal query model", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	values, err := d.incidentioClient.GetIncidents(ctx, query.DataQuery.TimeRange.From, query.DataQuery.TimeRange.To, qm.Filters, qm.Limit)
	if err != nil {
		d.logger.Error("Failed to get incidents", "error", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return backend.ErrorResponseWithErrorSource(err)
	}

	var createdAt []time.Time
	var id []string
	var reference []string
	var name []string
	var severity []string
	var status []string
	var incidentType []string
	var summary []string
	var permalink []string
	var workloadMinutesTotal []float64
	var workloadMinutesWorking []float64
	var workloadMinutesLate []float64
	var workloadMinutesSleeping []float64

	durations := make(map[string][]int64)
	customFields := make(map[string][]string)

	for _, value := range values {
		createdAt = append(createdAt, value.CreatedAt)
		id = append(id, value.ID)
		reference = append(reference, value.Reference)
		name = append(name, value.Name)
		severity = append(severity, value.Severity.Name)
		status = append(status, value.IncidentStatus.Name)
		incidentType = append(incidentType, value.IncidentType.Name)
		summary = append(summary, value.Summary)
		permalink = append(permalink, value.Permalink)
		workloadMinutesTotal = append(workloadMinutesTotal, value.WorkloadMinutesTotal)
		workloadMinutesWorking = append(workloadMinutesWorking, value.WorkloadMinutesWorking)
		workloadMinutesLate = append(workloadMinutesLate, value.WorkloadMinutesLate)
		workloadMinutesSleeping = append(workloadMinutesSleeping, value.WorkloadMinutesSleeping)

		for _, duration := range value.DurationMetrics {
			durationName := duration.DurationMetric.Name
			durationValue := duration.ValueSeconds
			durations[durationName] = append(durations[durationName], durationValue)
		}

		for _, customField := range value.CustomFieldEntries {
			var values []string
			for _, value := range customField.Values {
				if value.ValueText != "" {
					values = append(values, value.ValueText)
				} else if value.ValueCatalogEntry.Name != "" {
					values = append(values, value.ValueCatalogEntry.Name)
				} else if value.ValueOption.Value != "" {
					values = append(values, value.ValueOption.Value)
				} else {
					values = append(values, "")
				}
			}

			customFieldName := customField.CustomField.Name
			customFields[customFieldName] = append(customFields[customFieldName], strings.Join(values, ","))
		}
	}

	frame := data.NewFrame(
		"Incidents",
		data.NewField("Created At", nil, createdAt),
		data.NewField("ID", nil, id),
		data.NewField("Reference", nil, reference),
		data.NewField("Name", nil, name),
		data.NewField("Severity", nil, severity),
		data.NewField("Status", nil, status),
		data.NewField("Type", nil, incidentType),
		data.NewField("Summary", nil, summary),
		data.NewField("Permalink", nil, permalink),
		data.NewField("Workload Minutes Total", nil, workloadMinutesTotal),
		data.NewField("Workload Minutes Working", nil, workloadMinutesWorking),
		data.NewField("Workload Minutes Late", nil, workloadMinutesLate),
		data.NewField("Workload Minutes Sleeping", nil, workloadMinutesSleeping),
	)

	for durationName, durationValues := range durations {
		frame.Fields = append(frame.Fields, data.NewField(durationName, nil, durationValues))
	}
	for customFieldName, customFieldValues := range customFields {
		frame.Fields = append(frame.Fields, data.NewField(customFieldName, nil, customFieldValues))
	}

	frame.SetMeta(&data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
		Type:                   data.FrameTypeTable,
	})

	var response backend.DataResponse
	response.Frames = append(response.Frames, frame)

	return response
}
