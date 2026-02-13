package asteriskamitools

import "testing"

func TestRespToMap(T *testing.T) {
	TestedStr := "Data:sip:100@10.0.0.5:5060"
	TestStr := respToMap(TestedStr)
	for k, v := range TestStr {
		T.Log(k)
		T.Log(v)
	}
}
