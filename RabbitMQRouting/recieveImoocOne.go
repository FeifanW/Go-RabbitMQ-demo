package main

import "imooc-RabbitMQ/RabbitMQ"

func main() {
	imoocOne := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	imoocOne.RecieveRouting()
}
