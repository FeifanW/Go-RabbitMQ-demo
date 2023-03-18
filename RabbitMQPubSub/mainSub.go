package main

import "imooc-RabbitMQ/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" + "newProduct")
	rabbitmq.RecieveSub()
}
