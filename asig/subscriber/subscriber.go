package sub

import (
	"github.com/streadway/amqp"
	// "encoding/json"
	"asig/utils"
	"asig/model"
	"strings"
	"time"
)

type Subscriber struct {
	connection *amqp.Connection
	channel *amqp.Channel
	Queue string
	Consumer string
	auto_ack bool
	exclusive bool
	no_local bool
	no_wait bool
	args amqp.Table
	Msgs <-chan amqp.Delivery
}

func NewSubscriber(username, password, host, port, queue, consumer string, auto_ack, exclusive, no_local, no_wait bool, args amqp.Table) (s *Subscriber, err error) {
	s = &Subscriber{
		Queue: queue,
		Consumer: consumer,
		auto_ack: auto_ack,
		exclusive: exclusive,
		no_local: no_local,
		no_wait: no_wait,
		args: args,
	}

	// Dail amqp Connection
	if s.connection, err = amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port); err != nil {
		return
	}
	// Create amqp Channel
	if s.channel, err = s.connection.Channel(); err != nil {
		return
	}
	// s.Deregister functions can be implemented. Currently not implemented.
	return
}

func (s *Subscriber) DeRegister() {
	s.channel.Cancel(
			s.Consumer,		// Consumer
			false,			// No Wait
		)
	s.channel.Close()
	s.connection.Close()
}

func (s *Subscriber) Register() (err error) {
	for {
		s.Msgs, err = s.channel.Consume(
			s.Queue,		/* Queue name */
			s.Consumer,		/* Consumer */
			s.auto_ack,		/* Auto-Ack Enable */
			s.exclusive,	/* Exclusive */
			s.no_local,		/* No-Local */
			s.no_wait,		/* No-wait */
			s.args,			/* Args */
		)
		if err != nil {
			if strings.Contains(err.Error(), "NOT_FOUND - no queue") {
				utils.Log("failure", "995", "Unable to subscribe queue! Error:"+err.Error(), map[string]interface{}{"sub_queue" : s.Queue, "consumer" : s.Consumer})
				return
			}
			utils.Log("failure", "995", "Unable to subscribe queue, retrying! Error:"+err.Error(), map[string]interface{}{"sub_queue" : s.Queue, "consumer" : s.Consumer})
			time.Sleep(100*time.Millisecond)
		} else {
			utils.Log("success", "299", "Successfuly Registered RabbitMQ Subscriber!", map[string]interface{}{"sub_queue" : s.Queue, "consumer" : s.Consumer})
			return
		}
	}
}

func (s *Subscriber) Worker() {
	notify := s.connection.NotifyClose(make(chan *amqp.Error))
	// var payloads map[string]interface{}
	// var payload interface{}
	var msg amqp.Delivery
	var err error
	for {
		select{
		case err = <-notify:
			utils.Log("failure", "994", err.Error(), map[string]interface{}{"sub_queue" : s.Queue, "consumer" : s.Consumer})
			return
		case msg = <-s.Msgs:
			// if err = json.Unmarshal(msg.Body, &payloads); err != nil {
			// 	utils.Log("failure", "998", err.Error(), map[string]interface{}{"input": string(msg.Body), "sub_queue" : s.Queue, "consumer" : s.Consumer})
			// 	continue
			// }
			// for _, payload = range payloads["offers"].([]interface{}) {
			// 	model.Load(payload.(map[string]interface{}))
			// }
			model.Load(msg.Body)
		}
	}
}