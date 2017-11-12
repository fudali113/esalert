package storage

import (
	"testing"
)

var testErRequest = EsStorage{
	Host:     "localhost",
	Port:     "9200",
	Username: "elastic",
	Password: "changme",
	Index:    "logstash-*",
}

func TestEsRequest_getUrl(t *testing.T) {
	url := testErRequest.getURL()
	if url != "http://localhost:9200/logstash-*/_search" {
		t.Error(url)
	}
}

func TestEsRequest_RunQuery(t *testing.T) {
	res, err := testErRequest.GetData()
	if err != nil {
		t.Error(res, err)
	}
	// if hits.Total < 10 {
	// 	t.Fail()
	// }
	if err != nil {
		t.Error(err)
	}
}
