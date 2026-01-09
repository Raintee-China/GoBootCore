package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient 是一个封装了 RabbitMQ 连接和通道的客户端结构体。
type RabbitMQClient struct {
	Conn    *amqp.Connection // 与 RabbitMQ 服务器的连接
	Channel *amqp.Channel    // 用于发送和接收消息的通道
}

// NewRabbitMQClient 创建一个新的 RabbitMQ 客户端实例。
// 参数:
//   - host: RabbitMQ 服务器地址
//   - port: RabbitMQ 服务器端口
//   - username: 登录用户名
//   - password: 登录密码
//
// 返回值:
//   - *RabbitMQClient: 成功时返回 RabbitMQ 客户端实例
//   - error: 失败时返回错误信息
func NewRabbitMQClient(host string, port int, username string, password string) (*RabbitMQClient, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close 关闭 RabbitMQ 的通道和连接。
func (c *RabbitMQClient) Close() {
	c.Channel.Close()
	c.Conn.Close()
}

// Publish 向指定队列发布一条消息。
// 参数:
//   - queueName: 目标队列名称
//   - message: 要发布的消息内容
//
// 返回值:
//   - error: 发布过程中发生的错误
func (c *RabbitMQClient) Publish(queueName, message string) error {
	return c.Channel.Publish(
		"",        // exchange: 使用默认交换机
		queueName, // routing key: 队列名称
		false,     // mandatory: 如果为 true，当消息无法路由到队列时会返回给发布者
		false,     // immediate: 如果为 true，当消息不能立即投递时会返回给发布者（已废弃）
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// DeclareQueue 声明一个队列，如果该队列不存在则创建它。
// 参数:
//   - name: 队列名称
//
// 返回值:
//   - error: 声明或创建队列过程中发生的错误
func (c *RabbitMQClient) DeclareQueue(name string) error {
	_, err := c.Channel.QueueDeclare(
		name,  // name: 队列名称
		false, // durable: 队列是否持久化
		false, // delete when unused: 是否在不使用时自动删除
		false, // exclusive: 是否为排他队列（仅当前连接可用）
		false, // no-wait: 是否等待服务器确认
		nil,   // arguments: 队列参数（可选）
	)
	return err
}

// Consume 从指定队列消费消息。
// 参数:
//   - queueName: 消息来源队列名称
//
// 返回值:
//   - <-chan amqp.Delivery: 消息传递通道
//   - error: 消费过程中的错误
func (c *RabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.Channel.Consume(
		queueName, // queue: 队列名称
		"",        // consumer: 消费者标识符（空表示自动生成）
		true,      // auto-ack: 是否自动确认消息
		false,     // exclusive: 是否为排他消费者
		false,     // no-local: 是否禁止同一连接发送的消息被自己消费
		false,     // no-wait: 是否等待服务器确认
		nil,       // args: 其他参数（可选）
	)
	return msgs, err
}
