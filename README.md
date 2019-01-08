opencensus-go-exporter-iowriter
===============================

Overview
---------

This exporter implements both the `ExportSpan(sd *trace.SpanData)` and
`ExportView(vd *view.Data)` interfaces.
By default, it will serialize the data in `json` format and will write it
out to `os.Stdout`.
It can be configured to write out to any `io.writer`.

The [example](examples/main.go) takes the example code from the
[opencensus tutorial](https://opencensus.io/exporters/custom-exporter/go/) and
demonstrates how to integrate it with this `iowriter`.


Demo
----

To see the exporter in action, do the following:

```
go get github.com/gtrevg/opencensus-go-exporter-iowriter
cd ${GOPATH}/src/github.com/gtrevg/opencensus-go-exporter-iowriter/examples
go get -u
go build
./examples
```


To-do:
------
* [ ] Bug - Fix logging of `Rows[i].Tags.Key` as the value is currently `{}`
* [ ] Feature - Do not log empty keys in order to reduce log size
* [ ] Feature - Provide option not to log view data description
