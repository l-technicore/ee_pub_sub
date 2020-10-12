package pub

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"asig/utils"
)

type Publisher struct {
	connection *amqp.Connection
	channel *amqp.Channel
	Exchange string
	Queue string
}

func NewPublisher(username, password, host, port, exchange, queue string) (p *Publisher, err error) {
	p = &Publisher{
		Exchange: exchange,
		Queue: queue,
	}
	// Dail amqp Connection
	if p.connection, err = amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port); err != nil {
		return
	}
	// Create amqp Channel
	if p.channel, err = p.connection.Channel(); err != nil {
		return
	}
	return
}

type Payload struct {
	mandatory bool
	immediate bool
	body []byte
}

func (p *Publisher) Publish(mandatory, immediate bool, body []byte) error {
	return p.channel.Publish(
				p.Exchange, 	// exchange
				p.Queue,    	// queue/routing key
				mandatory,    	// mandatory
				immediate,    	// immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        body,
				},
			)
}

func (p *Publisher) Worker(ch chan *Payload) {
	var pl *Payload
	var plMap map[string]interface{}
	var err error
	for pl = range ch {
		if err = json.Unmarshal(pl.body, &plMap); err != nil {
			utils.Log("failure", "998", err.Error(), map[string]interface{}{"input": string(pl.body), "pub_exchange":p.Exchange, "pub_queue": p.Queue})
			continue
		}
		if err = p.Publish(pl.mandatory, pl.immediate, pl.body); err != nil {
			utils.Log("failure", "997", err.Error(), map[string]interface{}{"payload": plMap, "pub_exchange":p.Exchange, "pub_queue": p.Queue})
			continue
		}
		utils.Log("success", "200", err.Error(), map[string]interface{}{"payload": plMap, "pub_exchange":p.Exchange, "pub_queue": p.Queue})
	}
}