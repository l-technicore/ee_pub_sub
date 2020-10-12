package sub

import (
	"github.com/streadway/amqp"
	"testing"
	"asig/publisher"
	"asig/pub/input"
	"strings"
	"time"
	"asig/utils"
)

func TestNewSubscriber(t *testing.T) {
	// Dail amqp Connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Error("Unable to connect RabbitMQ. Test Aborted. Error: "+err.Error())
	} else {
		defer conn.Close()
	}
	// Create amqp Channel
	var ch *amqp.Channel
	ch, err = conn.Channel()
	if err != nil {
		t.Error("Unable to get channel to RabbitMQ. Test Aborted. Error: "+err.Error())
	} else {
		defer ch.Close()
	}
	for {
		if _, err := ch.QueueDeclare(
			"CM_OFFER_LOAD_TEST",			// Queue
			true,						// Durable
			false,						// AutoDelete
			false,						// Exclusive
			false,						// NoWait
			nil,						// Arguments
		); err != nil {
			if !strings.Contains(err.Error(), "channel/connection is not open") {
				time.Sleep(10*time.Millisecond)
			} else {
				t.Error(err.Error())
			}
		} else {
			break
		}
	}
	s, err := NewSubscriber("guest", "guest", "localhost", "5672", "CM_OFFER_LOAD_TEST", "cm_offer_loader", false, false, false, false, nil)
	if err != nil {
		t.Error(err.Error())
	} else if s == nil {
		t.Error("Subscriber Structure is null.")
	}
	if s.Queue != "CM_OFFER_LOAD_TEST" {
		t.Error("Queue in Subscriber struct does not match with input.")
	}
	if s.Consumer != "cm_offer_loader" {
		t.Error("Consumer in Subscriber struct does not match with input.")
	}
	if s.auto_ack != false {
		t.Error("auto_ack in Subscriber struct does not match with input.")
	}
	if s.exclusive != false {
		t.Error("exclusive in Subscriber struct does not match with input.")
	}
	if s.no_local != false {
		t.Error("no_local in Subscriber struct does not match with input.")
	}
	if s.no_wait != false {
		t.Error("no_wait in Subscriber struct does not match with input.")
	}
	if s.args != nil {
		t.Error("args in Subscriber struct does not match with input.")
	}
	utils.Log("success", "297", "Successfuly Creater New RabbitMQ Subscriber!", nil)
}

func TestRegister(t *testing.T) {
	s := &Subscriber{
		Queue: "CM_OFFER_LOAD_TEST",
		Consumer: "cm_offer_loader",
		auto_ack: false,
		exclusive: false,
		no_local: false,
		no_wait: false,
		args: nil,
	}
	var err error
	// Dail amqp Connection
	if s.connection, err = amqp.Dial("amqp://guest:guest@localhost:5672"); err != nil {
		t.Error("Unable to connect RabbitMQ. Test Aborted. Error: "+err.Error())
	} else {
		defer s.connection.Close()
	}
	// Create amqp Channel
	if s.channel, err = s.connection.Channel(); err != nil {
		t.Error("Unable to get channel to RabbitMQ. Test Aborted. Error: "+err.Error())
	} else {
		defer s.channel.Close()
	}
	// Try Registering consumer on rabbitMQ until Successful
	if err = s.Register(); err != nil {
		t.Error("Registration failed. Error: "+err.Error())
	} else {
		defer s.DeRegister()
	}
}

func TestWorker(t *testing.T) {
	if p, err := pub.NewPublisher("guest", "guest", "localhost", "5672", "", "CM_OFFER_LOAD_TEST"); err != nil {
		t.Error(err.Error())
	} else if err = p.Publish(false, false, []byte(input.Json)); err != nil {
		t.Error(err.Error())
	}
	utils.Log("success", "296", "Successfuly Published messages using RabbitMQ Publisher!", nil)
	s, err := NewSubscriber("guest", "guest", "localhost", "5672", "CM_OFFER_LOAD_TEST", "cm_offer_loader_1", false, false, false, false, nil)
	if err != nil {
		t.Error(err.Error())
	} 
	// Try Registering consumer on rabbitMQ until Successful
	if err = s.Register(); err != nil {
		t.Error(err.Error())
	} else {
		// defer s.DeRegister()
	}
	notify := s.connection.NotifyClose(make(chan *amqp.Error))
	var msg amqp.Delivery
	var msg_recieved string
	func() {
		for {
			select{
			case err = <-notify:
				t.Error(err.Error())
			case msg = <-s.Msgs:
				msg_recieved = string(msg.Body)
				msg.Ack(false)
				return
			}
		}
	}()

	if msg_recieved != input.Json {
		t.Error("Message recieved different from message sent.")
	}
	utils.Log("success", "296", "Successfuly Fetched message Published using RMQ Publisher through RMQ Subscriber!", nil)
}