package encoder

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/lei006/go-assist/servers/datapacket"
)

func ParseSDP(sdpRaw string) *datapacket.SdpData {
	sdpMap := datapacket.MakeSdpData()
	var info *datapacket.SdpItem
	for _, line := range strings.Split(sdpRaw, "\n") {
		line = strings.TrimSpace(line)
		typeval := strings.SplitN(line, "=", 2)
		if len(typeval) == 2 {
			fields := strings.SplitN(typeval[1], " ", 2)
			switch typeval[0] {
			case "m":
				if len(fields) > 0 {
					switch fields[0] {
					case "audio", "video":
						sdpMap.Items[fields[0]] = &datapacket.SdpItem{AVType: fields[0]}
						info = sdpMap.Items[fields[0]]
						mfields := strings.Split(fields[1], " ")
						if len(mfields) >= 3 {
							info.PayloadType, _ = strconv.Atoi(mfields[2])
						}
					}
				}

			case "a":
				if info != nil {
					for _, field := range fields {
						keyval := strings.SplitN(field, ":", 2)
						if len(keyval) >= 2 {
							key := keyval[0]
							val := keyval[1]
							switch key {
							case "control":
								info.Control = val
							case "rtpmap":
								info.Rtpmap, _ = strconv.Atoi(val)
							}
						}
						keyval = strings.Split(field, "/")
						if len(keyval) >= 2 {
							key := keyval[0]
							switch key {
							case "MPEG4-GENERIC":
								info.Codec = "aac"
							case "H264":
								info.Codec = "h264"
							case "H265":
								info.Codec = "h265"
							}
							if i, err := strconv.Atoi(keyval[1]); err == nil {
								info.TimeScale = i
							}
						}
						keyval = strings.Split(field, ";")
						if len(keyval) > 1 {
							for _, field := range keyval {
								keyval := strings.SplitN(field, "=", 2)
								if len(keyval) == 2 {
									key := strings.TrimSpace(keyval[0])
									val := keyval[1]
									switch key {
									case "config":
										info.Config, _ = hex.DecodeString(val)
									case "sizelength":
										info.SizeLength, _ = strconv.Atoi(val)
									case "indexlength":
										info.IndexLength, _ = strconv.Atoi(val)
									case "sprop-parameter-sets":
										fmt.Println("sprop-parameter-sets:", val)
										fields := strings.Split(val, ",")
										for _, field := range fields {
											val, _ := base64.StdEncoding.DecodeString(field)
											fmt.Println("           :", val)
											info.SpropParameterSets = append(info.SpropParameterSets, val)
										}

									}
								}
							}
						}
					}
				}

			}
		}
	}
	return sdpMap
}
