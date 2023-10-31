module github.com/flylib/mq/rabbitmq

go 1.18

require (
	github.com/flylib/mq v0.0.0-20231013034215-85ddffca41eb
	github.com/rabbitmq/amqp091-go v1.9.0
)

require (
	github.com/flylib/goutils/codec/json v0.0.0-20231012070911-2cf6c2bcb71d // indirect
	github.com/flylib/interface v0.0.0-20231030075616-76c4e9b38c2a // indirect
	github.com/flylib/pkg/log/builtinlog v0.0.0-20231031020815-c8a2ce9c3f8a // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/flylib/mq v0.0.0-20231013034215-85ddffca41eb => ../../mq
