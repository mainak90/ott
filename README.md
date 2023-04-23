# OTT: Opentelemetry Tracing Utilities
Golang provider library that provides tracers for auto instrumentation of tracing providers based on services using http handlers/routers.

## Providers
Providers are distributed tracing backends that provides and interface that is compatible with [OpenTelemetry specification](https://opentelemetry.io/docs/reference/specification/). 
Currently this tracer support the following providers
* `Zipkin`
* `Jaeger`
* `noop`
* `stdout`

## Usage
Using the library is simple and straightforward.
1. Import the library on your router module/package.
```go
import tl "github.com/mainak90/ott"
```
2. Initialize and empty `Config` struct.
```go
cfg := tl.Config{AppName: "Test-App",Provider: "stdout"}
```

3. Use the config struct to initialize the tracer utility.
```go
trace, err := tl.InitTracing(cfg, "Test-App")
if err != nil {
	fmt.Errorf(err)
}
```

4. Get the router options
```go
opts := tl.GetMuxOptions(trace)
```

5. Get the http transport config
```go
tr := tl.NewTransport(trace)
```

6. Update your mux router with the updated config. In application
```go
r := mux.NewRouter()
client := &http.Client{Transport: tr}
r.Use(otelmux.Middleware("Test-App", opts...), tl.EchoFirstTraceNodeInfo(trace.Propagator))
r.HandleFunc("/", requestHandlerFunc)
```

7. While running the app, the library can be tested by using checking the traceID and spanID of each request.
```go
traceID, spanID, _ := tl.ExtractTraceInfo(r.Context())
```

## License
Uses the MIT license. Please check out [LICENSE.md](./LICENSE.md) for more details.