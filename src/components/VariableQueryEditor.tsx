import { QueryEditorProps } from '@grafana/data';
import {
  Combobox,
  ComboboxOption,
  InlineField,
  InlineFieldRow,
} from '@grafana/ui';
import React from 'react';

import { DataSource } from '../datasource';
import { DEFAULT_QUERIES, Options, Query, QueryType } from '../types';
import { AttributeField } from './AttributeField';

interface Props extends QueryEditorProps<DataSource, any, Options, Query> { }

export function VariableQueryEditor({
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
            options={[{ label: 'Attribute Values', value: 'attributevalues' }]}
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
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField label="Type" labelWidth={15}>
          <Combobox<string>
            width={35}
            value={query.type}
            options={[
              { label: 'Alerts', value: 'alerts' },
              { label: 'Incidents', value: 'incidents' },
            ]}
            onChange={(option: ComboboxOption<string>) => {
              onChange({
                ...query,
                type: option.value,
              });
            }}
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <AttributeField
          datasource={datasource}
          type={query.type}
          attribute={query.attribute}
          onAttributeChange={(newAttribute: string) => {
            onChange({ ...query, attribute: newAttribute });
            onRunQuery();
          }}
        />
      </InlineFieldRow>
    </>
  );
}
