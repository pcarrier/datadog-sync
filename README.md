# datadog-sync: automate Datadog monitor updates

The Datadog user interface is great, but standard files can be a better way to maintain your monitors.

Maybe you'd like a backup. Maybe you'd like to use Git and pull requests. Or maybe your infrastructure has
many similar environments and you'd like to provision monitors programmatically.

Whatever the scenario, we hope `datadog-sync` can help.

## Installation

We do not currently provide binaries or packages.

Please [setup a Go environment](https://golang.org/doc/install), then run:

    $ go install github.com/meteor/datadog-sync

This tool will then be available as `$GOPATH/bin/datadog-sync` and can be copied around (but not shared
between Windows, OSX and Linux, nor across CPU architectures).

## Format

We follow the same format as the Datadog APIs.
Please see ["Create a monitor" in the Datadog API documentation](http://docs.datadoghq.com/api/#monitor-create).

## Usage

`datadog-sync` requires both an API key and an application key.

You can find them in [your Datadog account settings](https://app.datadoghq.com/account/settings#api),
and pass them either through the `-api-key`/`-app-key` command line arguments or the
`DATADOG_API_KEY`/`DATADOG_APP_KEY` environment variables.

### General flags

- `-only [regex]` restricts the set of monitors being processed. `datadog-sync` only cares about monitors
whose title match `regex` (partial match, so `foobarbaz` is a match for `-only bar`).
For example, prefix alert names with "Prod: " or "Staging: " then maintain them through separate files
using `-only '^Prod:'` and `-only '^Staging:'`.

- `-http-debug` is useful to log the exact requests and responses, for example when API calls fail.
  Please help us help you by providing the resulting output whenever reporting an issue.
  Your keys should be hidden; please double-check.

- `-format [json|yaml]` allows you to choose between JSON and YAML for the output of `-mode pull`
  and the input of `-mode push`.
  YAML is the default as it is much more friendly to humans for so many long strings.
  We only print JSON in its raw form; please look at [jq](https://github.com/stedolan/jq)
  for pretty-printing, colors and much more.

### `-mode pull` (default)

This mode dumps your current Datadog metrics to standard output.

- `-ids` will dump the monitor IDs. `-mode push` uses IDs to update monitors in place instead of deleting and
  recreating a monitor whenever it changes; however it will fail if you try to push a metric with an ID that no
  longer exists on Datadog.

### `-mode push`

This mode reads metrics from standard input and syncs them on Datadog (creating metrics in the input that aren't
on Datadog, deleting Datadog metrics that aren't in the input, updating metrics if the input provides an ID).

- `-dry-run` outputs the various changes without performing them. We highly recommend trying it before making
any changes.

- `-verbose` will log details about each change, instead of only showing the metric ID and name.

## Similar projects

- [Interferon](https://github.com/airbnb/interferon) from Airbnb
- [Barkdog](https://github.com/winebarrel/barkdog) from Genki Sugawara
