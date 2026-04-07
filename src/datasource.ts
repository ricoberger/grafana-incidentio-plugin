import {
  CoreApp,
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  LegacyMetricFindQueryOptions,
  MetricFindValue,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { lastValueFrom, Observable } from 'rxjs';

import { DEFAULT_QUERY, Filter, Options, Query } from './types';
import { VariableSupport } from './variablesupport';

export class DataSource extends DataSourceWithBackend<Query, Options> {
  constructor(instanceSettings: DataSourceInstanceSettings<Options>) {
    super(instanceSettings);
    this.variables = new VariableSupport(this);
  }

  getDefaultQuery(_: CoreApp): Partial<Query> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: Query, scopedVars: ScopedVars) {
    const filters: Filter[] = [];
    if (query.filters) {
      for (const filter of query.filters) {
        const value = getTemplateSrv().replace(filter.value, scopedVars);
        filters.push({ ...filter, value: value });
      }
    }

    return {
      ...query,
      queryType: query.queryType || DEFAULT_QUERY.queryType,
      filters: filters,
    };
  }

  query(request: DataQueryRequest<Query>): Observable<DataQueryResponse> {
    return super.query(request);
  }

  async metricFindQuery(
    query: Query,
    options?: LegacyMetricFindQueryOptions,
  ): Promise<MetricFindValue[]> {
    const q = this.query({
      targets: [
        {
          ...query,
          refId: query.refId
            ? `metricsFindQuery-${query.refId}`
            : 'metricFindQuery',
        },
      ],
      range: options?.range,
    } as DataQueryRequest<Query>);

    const response = await lastValueFrom(q as Observable<DataQueryResponse>);

    if (
      response &&
      (!response.data.length || !response.data[0].fields.length)
    ) {
      return [];
    }

    return response
      ? (response.data[0] as DataFrame).fields[0].values.map((_, index) => {
        const name = (response.data[0] as DataFrame).fields[1].values[
          index
        ].toString();

        return {
          text: name,
          value: _.toString(),
        };
      })
      : [];
  }

  filterQuery(query: Query): boolean {
    if (!query.queryType) {
      return false;
    }

    return true;
  }
}
