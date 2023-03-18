package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url 格式是固定的amqp:// 账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://imoocuser:imoocuser@192.168.163.128:5672/imooc" // 虚拟机的地址

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 连接信息
	Mqurl string
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		conn:      nil,
		channel:   nil,
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 简单模式第一步：1.创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	return rabbitmq
}

// 简单模式第二步：简单模式下生成代码
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	// 保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 控制消息是否持久化
		false,
		// 是否自动删除
		false,
		// 是否就要排他性，只能自己可见，其他用户不能访问
		false,
		// 是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		// 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 3.简单模式第三步：简单模式下生成代码
func (r *RabbitMQ) ConsumeSimple() {
	// 1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	// 保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 控制消息是否持久化
		false,
		// 是否自动删除
		false,
		// 是否就要排他性，只能自己可见，其他用户不能访问
		false,
		// 是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞
		false,
		// 其他参数
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message,To exit press ctrl+C")
	<-forever
}

// 订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed on connect rabbitmq")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		// Yes表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(r.Exchange, "fanout", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exch"+"ange")
	// 2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		// 在Pub/Sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)
	// 消费信息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	fmt.Println("退出请按CTRL+C\n")
	<-forever
}

// 路由模式
// 创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 订阅模式一定是direct
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		// 要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 路由模式接收消息
func (r *RabbitMQ) RecieveRouting() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+"ange")
	// 2.试探性创建队列，这里注意队列名不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		// 需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

// 话题模式
// 创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	// 创建RabbitMQ实例, 和路由模式下的一样
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		// 要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 话题模式接收消息
// 要注意key,规则
// 其中*用于匹配一个单词，#用于匹配多个单词（可以是零个）
// 匹配 imooc .* 表示匹配 imooc.hello，但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+"ange")
	// 2.试探性创建队列，这里注意队列名不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机生成队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		// 需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}