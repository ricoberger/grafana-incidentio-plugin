package plugin

import (
	"context"

	"github.com/ricoberger/grafana-incidentio-plugin/pkg/incidentio"
	"github.com/ricoberger/grafana-incidentio-plugin/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/tracing"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin
// in runtime. In this example datasource instance implements
// backend.QueryDataHandler, backend.CheckHealthHandler interfaces. Plugin
// should not implement all these interfaces - only those which are required for
// a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func NewDatasource(_ context.Context, pCtx backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	logger := backend.Logger.With("datasource", pCtx.Name).With("datasourceId", pCtx.ID).With("datasourceUid", pCtx.UID)
	logger.Debug("Creating new datasource instance")

	settings, err := models.LoadPluginSettings(pCtx)
	if err != nil {
		logger.Error("Failed to load plugin settings", "error", err.Error())
		return nil, err
	}

	incidentioClient, err := incidentio.NewClient(logger, settings.Secrets.ApiKey)
	if err != nil {
		logger.Error("Failed to create incident.io client", "error", err.Error())
		return nil, err
	}

	ds := &Datasource{
		incidentioClient: incidentioClient,
		logger:           logger,
	}

	queryTypeMux := datasource.NewQueryTypeMux()
	queryTypeMux.HandleFunc(models.QueryTypeAttributes, ds.handleAttributesQueries)
	queryTypeMux.HandleFunc(models.QueryTypeAttributeValues, ds.handleAttributeValuesQueries)
	queryTypeMux.HandleFunc(models.QueryTypeAlerts, ds.handleAlertsQueries)
	queryTypeMux.HandleFunc(models.QueryTypeIncidents, ds.handleIncidentsQueries)
	ds.queryHandler = queryTypeMux

	return ds, nil
}

// Datasource is an example datasource which can respond to data queries,
// reports its health and has streaming skills.
type Datasource struct {
	queryHandler     backend.QueryDataHandler
	incidentioClient incidentio.Client
	logger           log.Logger
}

// QueryData handles multiple queries and returns multiple responses. The
// queries are matched by their QueryType field against a handler function. See
// the NewDatasource function where the query type multiplexer is created and
// handlers are registered.
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "QueryData")
	defer span.End()

	return d.queryHandler.QueryData(ctx, req)
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a
// new instance created. As soon as datasource settings change detected by SDK
// old datasource instance will be disposed and a new one will be created using
// NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// CheckHealth handles health checks sent from Grafana to the plugin. The main
// use case for these health checks is the test button on the datasource
// configuration page which allows users to verify that a datasource is working
// as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	_, err := d.incidentioClient.GetAlertAttributes(ctx)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Health check failed: " + err.Error()
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
