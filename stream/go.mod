module github.com/flylib/mq/stream

go 1.18

require (
	github.com/flylib/mq v0.0.0-20231013034215-85ddffca41eb
	github.com/mitchellh/mapstructure v1.5.0
	github.com/redis/go-redis/v9 v9.2.1
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/flylib/goutils/codec/json v0.0.0-20231012070911-2cf6c2bcb71d // indirect
	github.com/flylib/goutils/codec/protobuf v0.0.0-20231012070911-2cf6c2bcb71d // indirect
	github.com/flylib/goutils/logger/log v0.0.0-20231012070911-2cf6c2bcb71d // indirect
	github.com/flylib/interface v0.0.0-20231030075616-76c4e9b38c2a // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/flylib/mq v0.0.0-20231013034215-85ddffca41eb => ../../mq
