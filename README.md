`otracegen` is a simple http server useful to generate a trace using opentracing
sdk and the zipkin client.

you can run it:

```
go run cmd/main.go
```
Every time you do a `GET :8080/ping` it generate a span.

yes it is complex like a :rocket: .
