package asteriskamitools

import (
	"flag"
	"os"
	"testing"
)

var (
	pjsip            = flag.String("pjsip", "", "Asterisk with PJSIP ip address")
	pjsipuser        = flag.String("pjsipu", "", "Asterisk with PJSIP ami username")
	prjippassword    = flag.String("pjsipp", "", "Asterisk with PJSIP ami password")
	sipchain         = flag.String("chansip", "", "Asterisk with chan_sip ip address")
	sipchainuser     = flag.String("chansipu", "", "Asterisk with chan_sip ami username")
	sipchainpassword = flag.String("chansipp", "", "Asterisk with chan_sip ami password")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestRespToMap(T *testing.T) {
	TestedStr := "Contact:sip:100@10.0.0.5:5060"
	TestStr := respToMap(TestedStr)
	for k, v := range TestStr {
		T.Log(k)
		T.Log(v)
		if k != "Contact" {
			T.Fatal("error in RespToMap, incorrect key")
		}
		if v != "sip:100@10.0.0.5:5060" {
			T.Fatal("error in RespToMap, incorrect value")
		}
	}
}

func TestGetPJSIPEndpointsIPtoDataMap(t *testing.T) {
	Dt, er := GetPJSIPEndpointsIPtoDataMap(*pjsip, 5038, *pjsipuser, *prjippassword)
	if er != nil {
		t.Error(er)
	} else {
		for k, v := range Dt {
			t.Log("key:", k)
			t.Log("IP:", v.IP)
			t.Log("Contact:", v.Contacts)
			t.Log("Ext:", v.Extension)
			t.Log("Name", v.Name)
		}
	}
}

func TestGetSIPEndpointsIPtoDataMap(t *testing.T) {
	Dt, er := GetSIPEndpointsIPtoDataMap(*sipchain, 5038, *sipchainuser, *sipchainpassword)
	if er != nil {
		t.Error(er)
	} else {
		for k, v := range Dt {
			t.Log("key:", k)
			t.Log("IP:", v.IP)
			t.Log("Ext:", v.Extension)
			t.Log("Name", v.Name)
		}
	}
}
