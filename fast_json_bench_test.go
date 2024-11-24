package isunippets

import (
	"encoding/json"
	"testing"
)

func BenchmarkFastJSONMarshal(b *testing.B) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	data := JSONResponse{
		Message: "hello",
		Amount:  100,
	}

	for i := 0; i < b.N; i++ {
		_, err := FastJSONMarshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

// 比較用
func BenchmarkJsonMarshal(b *testing.B) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	data := JSONResponse{
		Message: "hello",
		Amount:  100,
	}

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

// 比較用
func BenchmarkFastJSONUnmarshal(b *testing.B) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	data := []byte(`{"message":"hello","amount":100}`)

	for i := 0; i < b.N; i++ {
		var response JSONResponse
		err := FastJSONUnmarshal(data, &response)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	type JSONResponse struct {
		Message string `json:"message"`
		Amount  int64  `json:"amount"`
	}

	data := []byte(`{"message":"hello","amount":100}`)

	for i := 0; i < b.N; i++ {
		var response JSONResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			b.Error(err)
		}
	}
}
