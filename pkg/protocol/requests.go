package protocol

import "encoding/json"

type KafkaRequest interface {
	GetApiKey() ApiKey
	Serialize() ([]byte, error)
}

type ProduceRequest struct {
	ApiKey    ApiKey `json:"api_key"`
	Partition int    `json:"partition"`
}

func (r ProduceRequest) GetApiKey() ApiKey {
	return r.ApiKey
}

func (r ProduceRequest) Serialize() ([]byte, error) {
	j, err := json.Marshal(r)
	if err != nil {
		return j, nil
	}
	return nil, err
}

func DeserializeProduceRequest(js []byte) (*ProduceRequest, error) {
	var p ProduceRequest
	err := json.Unmarshal(js, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
