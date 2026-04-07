# Grafana incident.io Plugin

The Grafana incident.io Plugin allows you to explore your incident.io alerts and
incidents via the Grafana.

<div align="center">
  <table>
    <tr>
      <td><img src="https://raw.githubusercontent.com/ricoberger/grafana-incidentio-plugin/refs/heads/main/src/img/screenshots/alerts.png" /></td>
      <td><img src="https://raw.githubusercontent.com/ricoberger/grafana-incidentio-plugin/refs/heads/main/src/img/screenshots/incidents.png" /></td>
    </tr>
  </table>
</div>

## Installation

1. Before you can install the plugin, you have to add
   `ricoberger-incidentio-datasource` to the
   [`allow_loading_unsigned_plugins`](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#allow_loading_unsigned_plugins)
   configuration option or to the `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS`
   environment variable.
2. The plugin can then be installed by adding
   `ricoberger-incidentio-datasource@<VERSION>@https://github.com/ricoberger/grafana-incidentio-plugin/releases/download/v<VERSION>/ricoberger-incidentio-datasource-<VERSION>.zip`
   to the
   [`preinstall_sync`](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#preinstall_sync)
   configuration option or the `GF_PLUGINS_PREINSTALL_SYNC` environment
   variable.

### Configuration File

```ini
[plugins]
allow_loading_unsigned_plugins = ricoberger-incidentio-datasource
preinstall_sync = ricoberger-incidentio-datasource@0.1.0@https://github.com/ricoberger/grafana-incidentio-plugin/releases/download/v0.1.0/ricoberger-incidentio-datasource-0.1.0.zip
```

### Environment Variables

```bash
export GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=ricoberger-incidentio-datasource
export GF_PLUGINS_PREINSTALL_SYNC=ricoberger-incidentio-datasource@0.1.0@https://github.com/ricoberger/grafana-incidentio-plugin/releases/download/v0.1.0/ricoberger-incidentio-datasource-0.1.0.zip
```

## Contributing

If you want to contribute to the project, please read through the
[contribution guideline](https://github.com/ricoberger/grafana-incidentio-plugin/blob/main/CONTRIBUTING.md).
Please also follow our
[code of conduct](https://github.com/ricoberger/grafana-incidentio-plugin/blob/main/CODE_OF_CONDUCT.md)
in all your interactions with the project.
