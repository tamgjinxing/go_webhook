package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func ParseWebhookDataAndCreated9735Event(bodyjson string) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(bodyjson), &result)
	if err != nil {
		log.Printf("获取到的jsonbody格式不对:%v\n", err)
		return
	}

	typess := result["type"]
	if typess == "incoming" {
		payment_request := result["payment_request"]
		preimage := result["preimage"]
		payer_pubkey := result["payer_pubkey"]
		metadata := result["metadata"].(map[string]interface{})

		zapRequestRaw := metadata["zap_request_raw"].(string)

		// 创建一个空的 map 来接收解析后的 JSON 数据
		var result2 map[string]interface{}

		// 将 JSON 字符串解析到 map 中
		err := json.Unmarshal([]byte(zapRequestRaw), &result2)
		if err != nil {
			log.Println(err)
			return
		}

		// 提取字段
		tags := result2["tags"].([]interface{}) // tags 是一个数组
		var needSendToRelays []interface{}
		var receiverPublicKey []interface{}

		// 打印 tags 数组中的元素
		for _, tag := range tags {
			tagPair := tag.([]interface{}) // 每个 tag 是一个数组

			if tagPair[0].(string) == "relays" {
				needSendToRelays = tagPair[1:]
			}

			if tagPair[0].(string) == "p" {
				receiverPublicKey = tagPair[1:]
			}
		}

		tags1 := make([]Tag, 0, 4) // 初始长度为0，容量为5

		description := []string{"description", zapRequestRaw}

		p := []string{"p", receiverPublicKey[0].(string)}
		P := []string{"P", payer_pubkey.(string)}
		bolt11 := []string{"bolt11", payment_request.(string)}
		preimage1 := []string{"preimage", preimage.(string)}

		tags1 = append(tags1, p)
		tags1 = append(tags1, P)
		tags1 = append(tags1, bolt11)
		tags1 = append(tags1, preimage1)
		tags1 = append(tags1, description)

		eventString, err := Gen9735Event("", tags1)
		fmt.Println(eventString)
		if err != nil {
			return
		}

		if len(needSendToRelays) > 0 {
			// StartRelayConnections(needSendToRelays, eventString)
		}
	}
}
