package main

import (
	"asig/publisher"
	"asig/pub/input"
	"encoding/json"
	"asig/utils"
)

func main(){
	if err := utils.InitiateLogger("application.log"); err == nil {
		defer utils.CloseLogger()
	}
	p, err := pub.NewPublisher("guest", "guest", "localhost", "5672", "", "CM_OFFER_LOAD") // (username, password, host, port, exchange, queue string)
	if err != nil {
		utils.Log("failure", "999", err.Error(), nil)
		return
	}
	// p.Channel.NotifyConfirm can be used to keep track of deleivery confirmations. Currently not implemented.
	var payloadMap map[string]interface{}
	if err = json.Unmarshal([]byte(input.Json), &payloadMap); err != nil {
		utils.Log("failure", "998", err.Error(), map[string]interface{}{"input": input.Json, "pub_exchange":p.Exchange, "pub_queue": p.Queue})
		return
	}
	if err = p.Publish(false, false, []byte(input.Json)); err != nil {
		utils.Log("failure", "997", err.Error(), map[string]interface{}{"payload": payloadMap, "pub_exchange":p.Exchange, "pub_queue": p.Queue})
		return
	}
	utils.Log("success", "200", "Offers Published!", map[string]interface{}{"payload": payloadMap, "pub_exchange":p.Exchange, "pub_queue": p.Queue})
}


