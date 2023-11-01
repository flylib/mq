module test

go 1.18

require (
	github.com/flylib/goutils/codec/json v0.0.0-20231026110424-19dfbb98ff56
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d
	github.com/flylib/mq/rabbitmq v0.0.0-20231031025750-b7bdc43a231e
	github.com/flylib/pkg/log/builtinlog v0.0.0-20231031025337-eee45d016863
)

require (
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/flylib/mq/rabbitmq v0.0.0-20231031025750-b7bdc43a231e => ../rabbitmq
