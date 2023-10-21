package protocol

import "encoding/json"

type KafkaRequest interface {
	GetApiKey() ApiKey
}

type KafkaRequestImpl struct {
	ApiKey ApiKey `json:"api_key"`
}

func (r KafkaRequestImpl) GetApiKey() ApiKey {
	return r.ApiKey
}

func SerializeKafkaRequest(r KafkaRequest) ([]byte, error) {
	j, err := json.Marshal(r)
	if err != nil {
		return j, nil
	}
	return nil, err
}

func DeserializeKafkaRequest(js []byte) (KafkaRequest, error) {
	var r KafkaRequestImpl
	err := json.Unmarshal(js, &r)
	if err != nil {
		return nil, err
	}
	switch r.GetApiKey() {
	case PRODUCE:
		var p ProduceRequest
		json.Unmarshal(js, &p)
		return p, nil
	default:
	}
	return nil, INVALID_REQUEST
}

type ProduceRequest struct {
	KafkaRequestImpl
	Acks             int                 `json:"acks"`
	PartitionRecords map[string][]string `json:"partition_records"`
}
