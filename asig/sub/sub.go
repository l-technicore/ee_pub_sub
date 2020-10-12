package main

import (
	"asig/subscriber"
	"asig/utils"
	"asig/model"
)

func main(){
	if err := utils.InitiateLogger("application.log"); err != nil {
		return
	} else {
		defer utils.CloseLogger()
	}
	if err := model.InitializeDB("offers.db"); err != nil {
		return
	}
	// (username, password, host, port, queue, consumer string, auto_ack, exclusive, no_local, no_wait bool, args amqp.Table)
	s, err := sub.NewSubscriber("guest", "guest", "localhost", "5672", "CM_OFFER_LOAD", "cm_offer_loader", false, false, false, false, nil)
	if err != nil {
		utils.Log("failure", "999", err.Error(), nil)
		return
	} 
	// Try Registering consumer on rabbitMQ until Successful
	if err = s.Register(); err != nil {
		utils.Log("failure", "999", err.Error(), nil)
		return
	}else {
		defer s.DeRegister()
	}
	// Start one Worker
	s.Worker()
}


