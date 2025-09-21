module Streamify/Cliente

go 1.24.5

require (
	Streamify/ServidorCanciones v0.0.0-00010101000000-000000000000
	Streamify/ServidorStreaming v0.0.0-00010101000000-000000000000
	github.com/faiface/beep v1.1.0
	google.golang.org/grpc v1.75.1
)

require (
	github.com/hajimehoshi/go-mp3 v0.3.0 // indirect
	github.com/hajimehoshi/oto v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp v0.0.0-20190306152737-a1d7652674e8 // indirect
	golang.org/x/image v0.0.0-20190227222117-0694c2d4d067 // indirect
	golang.org/x/mobile v0.0.0-20190415191353-3e0bab5405d6 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

replace Streamify/ServidorCanciones => ../ServidorCanciones

replace Streamify/ServidorStreaming => ../ServidorStreaming
