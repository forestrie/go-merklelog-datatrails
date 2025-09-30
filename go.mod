module github.com/forestrie/go-merklelog-datatrails

go 1.24

toolchain go1.24.4

require (
	github.com/datatrails/go-datatrails-common-api-gen v0.8.0
	github.com/datatrails/go-datatrails-merklelog/massifs v0.0.0
	github.com/datatrails/go-datatrails-merklelog/mmr v0.4.0
	github.com/datatrails/go-datatrails-serialization/eventsv1 v0.0.0-20240909144155-35c44c77e223
	github.com/datatrails/go-datatrails-simplehash v0.0.0-20241001114049-90fd7a82596f
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/datatrails/go-datatrails-common v0.30.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0 // indirect
	github.com/ldclabs/cose/go v0.0.0-20221214142927-d22c1cfc2154 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/veraison/go-cose v1.1.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/zeebo/bencode v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/datatrails/go-datatrails-merklelog/massifs => ../go-datatrails-merklelog/massifs

replace github.com/datatrails/go-datatrails-merklelog/mmr => ../go-datatrails-merklelog/mmr

replace github.com/datatrails/go-datatrails-common => ../go-datatrails-common

replace github.com/datatrails/go-datatrails-simplehash => ../go-datatrails-simplehash

replace github.com/datatrails/go-datatrails-serialization/eventsv1 => ../go-datatrails-serialization/eventsv1
