module example

go 1.21.1

require (
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/zerologr v1.2.3
	github.com/koolay/centrifuge-sse-client/sseclient v0.0.0
	github.com/rs/zerolog v1.30.0
)

require (
	github.com/centrifugal/gocent/v3 v3.2.0 // indirect
	github.com/centrifugal/protocol v0.10.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	github.com/segmentio/asm v1.1.4 // indirect
	github.com/segmentio/encoding v0.3.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.0.0-20191116160921-f9c825593386 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
)

replace github.com/koolay/centrifuge-sse-client/sseclient => ../
