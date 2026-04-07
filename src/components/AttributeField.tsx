import { Combobox, ComboboxOption, InlineField } from '@grafana/ui';
import React from 'react';
import { useAsync } from 'react-use';

import { DataSource } from '../datasource';

interface Props {
  datasource: DataSource;
  type?: string;
  attribute?: string;
  onAttributeChange: (value: string) => void;
}

export function AttributeField({
  datasource,
  type,
  attribute,
  onAttributeChange,
}: Props) {
  const state = useAsync(async (): Promise<ComboboxOption[]> => {
    const result = await datasource.metricFindQuery({
      refId: 'attributes',
      queryType: 'attributes',
      type: type,
    });

    const attributes = result.map((value) => {
      return { value: value.value as string, label: value.text };
    });
    return attributes;
  }, [datasource, type]);

  return (
    <InlineField label="Attribute" labelWidth={15}>
      <Combobox<string>
        data-testid="attribute-combobox"
        width={35}
        value={attribute || ''}
        createCustomValue={true}
        options={state.value || []}
        onChange={(option: ComboboxOption<string>) => {
          onAttributeChange(option.value);
        }}
      />
    </InlineField>
  );
}
