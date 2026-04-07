import { QueryEditorProps } from '@grafana/data';
import {
  Combobox,
  ComboboxOption,
  IconButton,
  InlineField,
  InlineFieldRow,
  Input,
} from '@grafana/ui';
import React, { ChangeEvent } from 'react';

import { DataSource } from '../datasource';
import { DEFAULT_QUERIES, Options, Query, QueryType } from '../types';
import { AttributeField } from './AttributeField';
import { AttributeValueField } from './AttributeValueField';

type Props = QueryEditorProps<DataSource, Query, Options>;

export function QueryEditor({
  datasource,
  query,
  onChange,
  onRunQuery,
}: Props) {
  return (
    <>
      <InlineFieldRow>
        <InlineField label="Query Type" labelWidth={15}>
          <Combobox<QueryType>
            width={35}
            value={query.queryType}
            options={[
              { label: 'Alerts', value: 'alerts' },
              { label: 'Incidents', value: 'incidents' },
            ]}
            onChange={(option: ComboboxOption<QueryType>) => {
              onChange({
                ...query,
                ...DEFAULT_QUERIES[option.value],
                queryType: option.value,
              });
              onRunQuery();
            }}
          />
        </InlineField>

        <IconButton
          name="plus"
          aria-label="Add Filter"
          onClick={() => {
            const newFilters = [
              ...(query.filters || []),
              { attribute: '', operator: '', value: '' },
            ];
            onChange({ ...query, filters: newFilters });
          }}
        />
      </InlineFieldRow>

      {(query.filters || []).map((filter, index) => (
        <InlineFieldRow key={index}>
          <AttributeField
            datasource={datasource}
            type={query.queryType}
            attribute={filter.attribute}
            onAttributeChange={(newAttribute: string) => {
              const newFilters = [...(query.filters || [])];
              newFilters[index] = {
                attribute: newAttribute,
                operator: '',
                value: '',
              };
              onChange({ ...query, filters: newFilters });
            }}
          />
          {filter.attribute && (
            <AttributeValueField
              datasource={datasource}
              type={query.queryType}
              attribute={filter.attribute}
              attributeValue={filter.value}
              onAttributeValueChange={(
                newOperator: string,
                newAttributeValue: string,
              ) => {
                const newFilters = [...(query.filters || [])];
                newFilters[index] = {
                  ...newFilters[index],
                  operator: newOperator,
                  value: newAttributeValue,
                };
                onChange({ ...query, filters: newFilters });
              }}
            />
          )}
          <IconButton
            name="minus"
            aria-label="Remove Filter"
            onClick={() => {
              const newFilters = [...(query.filters || [])];
              newFilters.splice(index, 1);
              onChange({ ...query, filters: newFilters });
            }}
          />
        </InlineFieldRow>
      ))}

      <InlineFieldRow>
        <InlineField label="Limit" labelWidth={15}>
          <Input
            width={35}
            value={query.limit || 50}
            onChange={(event: ChangeEvent<HTMLInputElement>) => {
              onChange({
                ...query,
                limit: parseInt(event.target.value, 10),
              });
            }}
          />
        </InlineField>
      </InlineFieldRow>
    </>
  );
}
