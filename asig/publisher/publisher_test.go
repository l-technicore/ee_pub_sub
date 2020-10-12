package pub

import (
	"github.com/streadway/amqp"
	"testing"
	"asig/subscriber"
	"asig/pub/input"
	"strings"
	"time"
	"asig/utils"
)

func TestNewPublisher(t *testing.T) {
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
	p, err := NewPublisher("guest", "guest", "localhost", "5672","", "CM_OFFER_LOAD_TEST")
	if err != nil {
		t.Error(err.Error())
	} else if p == nil {
		t.Error("Publisher Structure is null.")
	}
	if p.Queue != "CM_OFFER_LOAD_TEST" {
		t.Error("Queue in Subscriber struct does not match with input.")
	}
	if p.Exchange != "" {
		t.Error("Exchange in Subscriber struct does not match with input.")
	}
	utils.Log("success", "297", "Successfuly Creater New RabbitMQ Subscriber!", nil)
}

func TestPublish(t *testing.T) {
	if p, err := NewPublisher("guest", "guest", "localhost", "5672", "", "CM_OFFER_LOAD_TEST"); err != nil {
		t.Error(err.Error())
	} else if err = p.Publish(false, false, []byte(input.Json)); err != nil {
		t.Error(err.Error())
	}
	utils.Log("success", "296", "Successfuly Published messages using RabbitMQ Publisher!", nil)
	s, err := sub.NewSubscriber("guest", "guest", "localhost", "5672", "CM_OFFER_LOAD_TEST", "cm_offer_loader_1", false, false, false, false, nil)
	if err != nil {
		t.Error(err.Error())
	} 
	// Try Registering consumer on rabbitMQ until Successful
	if err = s.Register(); err != nil {
		t.Error(err.Error())
	} else {
		// defer s.DeRegister()
	}
	var msg amqp.Delivery
	var msg_recieved string
	func(){
		for {
			select{
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