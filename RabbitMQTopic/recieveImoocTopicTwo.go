package main

import "imooc-RabbitMQ/RabbitMQ"

func main() {
	imoocOne := RabbitMQ.NewRabbitMQTopic("exImoocTopic", "imooc.*.two")
	imoocOne.RecieveTopic()
}
