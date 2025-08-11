# Lowest Latency Endpoint Plugin

A Traefik middleware plugin that tests multiple endpoints for latency and sets a header with the fastest responding endpoint.

## How it works

The plugin performs HEAD requests to all configured endpoints and determines which one responds fastest. It then sets a header with the URL of the lowest latency endpoint.

## Configuration

Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    lowestlatency:
      moduleName: github.com/cilasbeltrame/lowestlatencyendpoint
      version: v1.0.0
```

Configure the middleware in your dynamic configuration:

```yaml
http:
  middlewares:
    latency-check:
      plugin:
        lowestlatency:
          endpoints:
            - https://api.example.com/health
            - https://api-eu.example.com/health
            - https://api-us.example.com/health
          headerName: X-Fastest-Endpoint

  routers:
    my-app:
      rule: Host(`myapp.localhost`)
      middlewares:
        - latency-check
      service: my-service
```

### Options

- `endpoints`: List of URLs to test for latency
- `headerName`: Name of the header to set with the fastest endpoint (default: "X-Lowest-Latency")

## Development

To test the plugin locally, place it in the `./plugins-local` directory structure:

```
./plugins-local/
    └── src
        └── github.com
            └── cilasbeltrame
                └── lowestlatencyendpoint
                    ├── main.go
                    ├── main_test.go
                    ├── go.mod
                    └── ...
```

Configure Traefik to use the local plugin:

```yaml
experimental:
  localPlugins:
    lowestlatency:
      moduleName: github.com/cilasbeltrame/lowestlatencyendpoint
```
