package main

import (
	"encoding/json"
	"log"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/env"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/starter"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/utils"
)

var (
	zipperAddr = env.GetString("YOMO_SOURCE_MQTT_ZIPPER_ADDR", "localhost:9999")
	brokerAddr = env.GetString("YOMO_SOURCE_MQTT_BROKER_ADDR", "0.0.0.0:1883")
	source     yomo.Source
)

type ThermometerData struct {
	Temperature float32 `json:"tem"`
	Humidity    float32 `json:"hum"`
}

func main() {
	// connect to YoMo-Zipper.
	source = yomo.NewSource("yomo-source", zipperAddr)
	err := source.Connect()
	if err != nil {
		log.Printf("[source] ❌ Connect to YoMo-Zipper %s failure with err: %v", zipperAddr, err)
		return
	}

	defer source.Close()

	// start a new MQTT Broker.
	starter.NewBrokerSimply(brokerAddr, "thermometer").
		Run(handler)
}

func handler(topic string, payload []byte) {
	log.Printf("receive: topic=%v, payload=%v\n", topic, string(payload))

	// get data from MQTT
	var data ThermometerData
	err := json.Unmarshal(payload, &data)
	if err != nil {
		log.Printf("Unmarshal payload error:%v", err)
	}

	sendingBuf, _ := json.Marshal(data)

	// send data to YoMo-Zipper.
	err = source.Write(0x34, sendingBuf)
	if err != nil {
		log.Printf("stream.Write error: %v, sendingBuf=%#x\n", err, sendingBuf)
	}

	log.Printf("write: sendingBuf=%v\n", utils.FormatBytes(sendingBuf))
}
