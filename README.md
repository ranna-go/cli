# Ranna CLI

## Usage

```
Usage cli:
  -auth string
        ranna API auth token
  -endpoint string
        ranna API endpoint (default "http://api.ranna.zekro.de")
  -f string
        code file to read (reads from stdin if not set)
  -i    Execute code inline
  -s string
        spec to use
  -v    Verbose logging
  -version string
        ranna API version (default "v1")
```

## Config

You can store persistent configuration in a JSON, YAML or TOML file either directly in a `config.*` file next to the binary or in a `.ranna` directory in your home path.

```yaml
endpoint:      https://public.ranna.zekro.de
version:       v1
authorization: some-optional-required-token
```

## Binaries

You can download latest binaries for most common architectures from the [**GitHub Action Artifacts**](https://github.com/ranna-go/cli/actions/workflows/artifacts.yml).