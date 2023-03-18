package main

import "imooc-RabbitMQ/RabbitMQ"

func main() {
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	imoocTwo.RecieveRouting()
}
