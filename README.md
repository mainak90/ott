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

7. `EchoFirstTraceInfo` injects trace info into the response path as long as the request is part of the parent span(first one in the request chain).

8. While running the app, the library can be tested by using checking the traceID and spanID of each request.
```go
traceID, spanID, _ := tl.ExtractTraceInfo(r.Context())
```

9. If the trace_ids and span_ids is not printed out from the helper function like below.
```go
log.Println(fmt.Sprintf("Trace ID for this request in %s is: %s and Span Id is: %s", "Test-App", traceID, spanID))
```

10. Throws generic output like
```azure
2023/04/23 21:48:49 Trace ID for this request in ServiceA is: 9ede35c46c2dc8297e5d3f15283c59c8 and Span Id is: 33837d83b34cfabb
```

11. This can also output the full span context in the response(if its `stdout` provider). A typical span context would look like this.
```azure
{
        "Name": "HTTP POST",
        "SpanContext": {
                "TraceID": "9ede35c46c2dc8297e5d3f15283c59c8",
                "SpanID": "d28674755081cfd9",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "Parent": {
                "TraceID": "9ede35c46c2dc8297e5d3f15283c59c8",
                "SpanID": "33837d83b34cfabb",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "SpanKind": 3,
        "StartTime": "2023-04-23T21:48:49.825656+02:00",
        "EndTime": "2023-04-23T21:48:49.827174434+02:00",
        "Attributes": [
                {
                        "Key": "http.method",
                        "Value": {
                                "Type": "STRING",
                                "Value": "POST"
                        }
                },
                {
                        "Key": "http.flavor",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.1"
                        }
                },
                {
                        "Key": "http.url",
                        "Value": {
                                "Type": "STRING",
                                "Value": "http://localhost:8001"
                        }
                },
                {
                        "Key": "net.peer.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "localhost"
                        }
                },
                {
                        "Key": "net.peer.port",
                        "Value": {
                                "Type": "INT64",
                                "Value": 8001
                        }
                },
                {
                        "Key": "http.status_code",
                        "Value": {
                                "Type": "INT64",
                                "Value": 200
                        }
                }
        ],
        "Events": null,
        "Links": null,
        "Status": {
                "Code": "Unset",
                "Description": ""
        },
        "DroppedAttributes": 0,
        "DroppedEvents": 0,
        "DroppedLinks": 0,
        "ChildSpanCount": 0,
        "Resource": [
                {
                        "Key": "service.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "unknown_service:___go_build_github_com_mainak90_ott_test_serviceA"
                        }
                },
                {
                        "Key": "telemetry.sdk.language",
                        "Value": {
                                "Type": "STRING",
                                "Value": "go"
                        }
                },
                {
                        "Key": "telemetry.sdk.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "opentelemetry"
                        }
                },
                {
                        "Key": "telemetry.sdk.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.14.0"
                        }
                }
        ],
        "InstrumentationLibrary": {
                "Name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
                "Version": "semver:0.40.0",
                "SchemaURL": ""
        }
}

```

## License
Uses the MIT license. Please check out [LICENSE.md](./LICENSE.md) for more details. 