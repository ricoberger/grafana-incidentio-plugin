import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { InlineField, SecretInput } from '@grafana/ui';
import React, { ChangeEvent } from 'react';

import { Options, OptionsSecure } from '../types';

interface Props
  extends DataSourcePluginOptionsEditorProps<Options, OptionsSecure> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { secureJsonFields, secureJsonData } = options;

  return (
    <>
      <InlineField
        data-testid="api-key"
        label="API Key"
        labelWidth={15}
        interactive
      >
        <SecretInput
          required
          isConfigured={secureJsonFields.apiKey}
          value={secureJsonData?.apiKey}
          width={35}
          onReset={() => {
            onOptionsChange({
              ...options,
              secureJsonFields: {
                ...options.secureJsonFields,
                apiKey: false,
              },
              secureJsonData: {
                ...options.secureJsonData,
                apiKey: '',
              },
            });
          }}
          onChange={(event: ChangeEvent<HTMLInputElement>) => {
            onOptionsChange({
              ...options,
              secureJsonData: {
                apiKey: event.target.value,
              },
            });
          }}
        />
      </InlineField>
    </>
  );
}
