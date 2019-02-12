# Trying out WebAssembly in Go (Cologne Go Meetup)

## Prerequisites

* Go 1.11.x, with the GOROOT env variable set to that installation
* Make

## How to use it

```
make all
make run
```

Then go to http://localhost:8080/index.html

## Goals achieved

* Compiling to WASM and running it in the browser (see https://github.com/golang/go/wiki/WebAssembly#getting-started)
* Fetching a Golang HTML template from the server, executing it on the client and injecting it into the DOM
* GET to a REST end point, and using the result to update the DOM, while using transport structs shared with the server

## Insights

* main() has to end with a select {}, so it does not exit, thereby killing all event handlers.
* the generated WASM file is quite large, around 10 MB in our example. Gzipped 2.2MB. Not mobile ready.
* Whenever something goes wrong, the page reloads. If you are lucky, you get a Go stack trace in the console.
* Go rpc does not work, as it wants to open a TCP connection.
* Using net/http in a JS event handler does not seem to work.

## Next steps

* Investigate https://github.com/dennwc/dom
* Interesting example there: gRPC over web sockets: https://github.com/dennwc/dom/tree/master/examples/grpc-over-ws
