import { ComboboxOption, InlineField, Input, MultiCombobox } from '@grafana/ui';
import React, { ChangeEvent } from 'react';
import { useAsync } from 'react-use';

import { DataSource } from '../datasource';

interface Props {
  datasource: DataSource;
  type?: string;
  attribute?: string;
  attributeValue?: string;
  onAttributeValueChange: (operator: string, value: string) => void;
}

export function AttributeValueField({
  datasource,
  type,
  attribute,
  attributeValue,
  onAttributeValueChange,
}: Props) {
  const state = useAsync(async (): Promise<ComboboxOption[]> => {
    const result = await datasource.metricFindQuery({
      refId: 'attributevalues',
      queryType: 'attributevalues',
      type: type,
      attribute: attribute || '',
    });

    const attributes = result.map((value) => {
      return { value: value.value as string, label: value.text };
    });
    return attributes;
  }, [datasource, type, attribute]);

  if (
    state.loading ||
    state.error ||
    !state.value ||
    state.value.length === 0
  ) {
    return (
      <InlineField label="Value" labelWidth={15}>
        <Input
          width={15}
          value={attributeValue || ''}
          onChange={(event: ChangeEvent<HTMLInputElement>) => {
            onAttributeValueChange('is', event.target.value);
          }}
        />
      </InlineField>
    );
  }

  return (
    <InlineField label="Value" labelWidth={15}>
      <MultiCombobox
        data-testid="attributevalue-combobox"
        width="auto"
        minWidth={35}
        maxWidth={35}
        isClearable={true}
        createCustomValue={true}
        value={
          !attributeValue || attributeValue.length === 0
            ? []
            : attributeValue.split(',')
        }
        options={state.value || []}
        onChange={(option: Array<ComboboxOption<string>>) => {
          onAttributeValueChange(
            'one_of',
            Array.from(option.values())
              .map((value) => value.value)
              .join(','),
          );
        }}
      />
    </InlineField>
  );
}
