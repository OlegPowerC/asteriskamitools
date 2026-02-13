package asteriskamitools

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type extData struct {
	Extension string
	Contacts  string
	Name      string
}

// ExtensionData contains complete information about a registered PJSIP endpoint.
type ExtensionData struct {
	Extension string
	Contacts  string
	IP        string
	Name      string
}

const AMI_AUTH_ACCEPTED = "Authentication accepted"
const AMI_LOGIN_CMD = "login"
const AMI_LOGOFF_CMD = "logoff"
const AMI_PJSIPSHOWENDPOINTS_CMD = "PJSIPShowEndpoints"
const AMI_PJSIPSHOWENDPOINT_CMD = "PJSIPShowEndpoint"
const AMI_PJSIPSHOWENDPOINTS_EVENT_VAL = "EndpointList"
const AMI_PJSIPSHOWENDPOINTS_OBJECTTYPE_VAL = "endpoint"

func respToMap(stresp string) map[string]string {
	rmap := make(map[string]string)
	ass := strings.Split(stresp, "\r\n")
	for _, assentry := range ass {
		assformap := strings.SplitN(assentry, ":", 2)
		if len(assformap) == 2 {
			rmap[assformap[0]] = strings.TrimSpace(assformap[1])
		}
	}
	return rmap
}

func endpointslistDataGet(inmap map[string]string, ActionID string) (string, string, error) {
	Ext := ""
	Contacts := ""
	if val, ok := inmap["ActionID"]; ok {
		if val == ActionID {
			if val, ok := inmap["Response"]; ok {
				if val == "Error" {
					return "", "", fmt.Errorf("Error: %s", inmap["Message"])
				}
			}
			if val, ok := inmap["Event"]; ok {
				if val == AMI_PJSIPSHOWENDPOINTS_EVENT_VAL {
					if val, ok := inmap["ObjectType"]; ok {
						if val == AMI_PJSIPSHOWENDPOINTS_OBJECTTYPE_VAL {
							if val, ok := inmap["ObjectName"]; ok {
								if val != "" {
									Ext = val
								}
							}
						}
					}
					if val, ok := inmap["Contacts"]; ok {
						if val != "" {
							Contacts = val
						}
					}
				}
			}
		}
	}
	return Ext, Contacts, nil
}

func endpointDetailGet(inmap map[string]string, ActionID string) (string, string, error) {
	Ext := ""
	CallerId := ""
	if val, ok := inmap["ActionID"]; ok {
		if val == ActionID {
			if val, ok := inmap["Response"]; ok {
				if val == "Error" {
					return "", "", fmt.Errorf("Error: %s", inmap["Message"])
				}
			}
			if val, ok := inmap["Event"]; ok {
				if val == "EndpointDetail" {
					if val, ok := inmap["ObjectType"]; ok {
						if val == AMI_PJSIPSHOWENDPOINTS_OBJECTTYPE_VAL {
							if val, ok := inmap["Callerid"]; ok {
								if val != "" {
									CallerId = val
								}
							}
						}
					}
					if val, ok := inmap["ObjectName"]; ok {
						if val != "" {
							Ext = val
						}
					}
				}
			}
		}
	}
	return Ext, CallerId, nil
}

func getIPExtNameMap(ExSlDt []extData) map[string]ExtensionData {
	RetMap := make(map[string]ExtensionData)
	for _, ExdVal := range ExSlDt {
		if ExdVal.Contacts != "" {
			lf1 := strings.Split(ExdVal.Contacts, "@")
			if len(lf1) == 2 {
				lf2 := strings.Split(lf1[1], ";")
				if len(lf2) > 0 {
					lf3 := strings.Split(lf2[0], ":")
					if len(lf3) > 0 {
						RetMap[lf3[0]] = ExtensionData{Extension: ExdVal.Extension, Contacts: ExdVal.Contacts, IP: lf3[0], Name: ExdVal.Name}
					}
				}
			}
		}
	}
	return RetMap
}

// GetPJSIPEndpointsIPtoDataMap retrieves PJSIP endpoint data from Asterisk via AMI protocol.
// Returns a map indexed by IP address containing extension numbers, contact URIs, and caller IDs.
//
// Parameters:
//   - AMIAddr: Asterisk server IP address or hostname
//   - AMIPort: AMI port (typically 5038)
//   - AMIUsername: AMI username
//   - AMIPassword: AMI password
//
// Example:
//
//	endpoints, err := GetPJSIPEndpointsIPtoDataMap("10.0.0.1", 5038, "admin", "secret")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for ip, data := range endpoints {
//	    fmt.Printf("IP: %s, Ext: %s, Name: %s\n", ip, data.Extension, data.Name)
//	}
func GetPJSIPEndpointsIPtoDataMap(AMIAddr string, AMIPort int, AMIUsername string, AMIPassword string) (map[string]ExtensionData, error) {
	ActionID := "23456063340"
	AuthStr := fmt.Sprintf("Action: %s\r\nUsername: %s\r\nSecret: %s\r\nEvents: off\r\nActionID: %s\r\n\r\n", AMI_LOGIN_CMD, AMIUsername, AMIPassword, ActionID)
	EndpointsStr := fmt.Sprintf("Action: %s\r\nActionID: %s\r\n\r\n", AMI_PJSIPSHOWENDPOINTS_CMD, ActionID)
	LogooffStr := fmt.Sprintf("Action: %s\r\n\r\n", AMI_LOGOFF_CMD)

	AMIurl := fmt.Sprintf("%s:%d", AMIAddr, AMIPort)
	conn, err := net.DialTimeout("tcp", AMIurl, 30*time.Second)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(180 * time.Second))
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		_, werr := conn.Write([]byte(AuthStr))
		if werr != nil {
			return nil, werr
		}
		break
	}
	if scanerr := scanner.Err(); scanerr != nil {
		return nil, scanerr
	}

	scanner2 := bufio.NewScanner(conn)
	scanner2.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), "\r\n\r\n"); i >= 0 {
			return i + len("\r\n\r\n"), data[0 : i+len("\r\n\r\n")], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return
	})

	scanner2.Scan()
	if scan2err := scanner2.Err(); scan2err != nil {
		return nil, scan2err
	}

	AuthRespT := scanner2.Text()

	authrespm := respToMap(AuthRespT)

	exts := make([]extData, 0)
	if authrespm["Message"] == AMI_AUTH_ACCEPTED {
		_, werr := conn.Write([]byte(EndpointsStr))
		if werr != nil {
			return nil, werr
		}
		for scanner2.Scan() {
			respt := scanner2.Text()
			sipregm := respToMap(respt)
			if val, ok := sipregm["EventList"]; ok {
				if val == "Complete" {
					break
				}
			}
			Ests, ExtCnt, Exerr := endpointslistDataGet(sipregm, ActionID)
			if Exerr == nil && Ests != "" {
				exts = append(exts, extData{Extension: Ests, Contacts: ExtCnt, Name: ""})
			} else {
				if Exerr != nil {
					return nil, Exerr
				}
			}
		}
		if scan2err := scanner2.Err(); scan2err != nil {
			return nil, scan2err
		}
	}

	for exindex, extnum := range exts {
		if extnum.Contacts == "" {
			continue
		}

		ShowEndpointDetail := fmt.Sprintf("Action: %s\r\nEndpoint: %s\r\nActionID: %s\r\n\r\n", AMI_PJSIPSHOWENDPOINT_CMD, extnum.Extension, ActionID)
		_, werr := conn.Write([]byte(ShowEndpointDetail))
		if werr != nil {
			return nil, werr
		}
		for scanner2.Scan() {
			respt := scanner2.Text()
			sipregm := respToMap(respt)
			if val, ok := sipregm["EventList"]; ok {
				if val == "Complete" {
					break
				}
			}
			Ext, CallerId, CallErr := endpointDetailGet(sipregm, ActionID)
			if CallerId != "" && Ext == extnum.Extension && CallErr == nil {
				exts[exindex].Name = CallerId
			}
		}
		if scan2err := scanner2.Err(); scan2err != nil {
			return nil, scan2err
		}
	}

	IpToExtAndNameMap := getIPExtNameMap(exts)

	_, werr := conn.Write([]byte(LogooffStr))
	if werr != nil {
		return nil, werr
	}
	scanner2.Scan()

	return IpToExtAndNameMap, nil
}
