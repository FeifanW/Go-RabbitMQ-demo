package main

import "imooc-RabbitMQ/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple") // 要和生产端的queueName一致
	rabbitmq.ConsumeSimple()
}
