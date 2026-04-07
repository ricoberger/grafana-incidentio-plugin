import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export const DEFAULT_QUERIES: Record<QueryType, Partial<Query>> = {
  attributes: {
    type: '',
  },
  attributevalues: {
    type: '',
    attribute: '',
  },
  alerts: {
    filters: [],
    limit: 50,
  },
  incidents: {
    filters: [],
    limit: 50,
  },
};

export const DEFAULT_QUERY: Partial<Query> = {
  queryType: 'alerts',
  ...DEFAULT_QUERIES['alerts'],
};

export type QueryType =
  | 'attributes'
  | 'attributevalues'
  | 'alerts'
  | 'incidents';

export interface Query
  extends DataQuery,
  QueryModelAttributes,
  QueryModelAttributeValues,
  QueryModelAlerts,
  QueryModelIncidents {
  queryType: QueryType;
}

interface QueryModelAttributes {
  type?: string;
}

interface QueryModelAttributeValues {
  type?: string;
  attribute?: string;
}

interface QueryModelAlerts {
  filters?: Filter[];
  limit?: number;
}

interface QueryModelIncidents { }

export interface Filter {
  attribute: string;
  operator: string;
  value: string;
}

export interface Options extends DataSourceJsonData { }

export interface OptionsSecure {
  apiKey?: string;
}
