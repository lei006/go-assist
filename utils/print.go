package utils

import "fmt"

func PrintRTP(data []byte) {

	if len(data) > 13 {
		FUindicator := data[12]
		FUheader := data[13]
		headerStart := (FUindicator & 0xE0) | (FUheader & 0x1F)
		fmt.Printf("-RtpData: %5d  %2d  %02d    ", len(data), data[12]&0x1f, headerStart&0x1f)
	}

	for i := 0; i < len(data); i++ {
		fmt.Printf(" %02x", data[i])
		if i == 12 {
			fmt.Print(" -")
		}
		if i == 13 {
			fmt.Print(" -")
		}
		if i == 11 {
			fmt.Print(" -")
		}

		if i > 30 {
			break
		}
	}

	fmt.Println("")

}

func PrintH264(data []byte) {

	if len(data) > 12 {
		fmt.Printf("-h264: %5d ", len(data))
	}

	for i := 0; i < len(data); i++ {
		fmt.Printf(" %02x", data[i])
		if i > 20 {
			break
		}
	}
	fmt.Println("")

}

func PrintBin(data []byte, limit int) {

	lineStr := ""

	if len(data) > 2 {
		lineStr += fmt.Sprintf("-data: %5d ", len(data))
	}

	for i := 0; i < len(data); i++ {

		if i%4 == 0 {
			lineStr += " "
		}
		if i%8 == 0 {
			lineStr += "  "
		}

		lineStr += fmt.Sprintf(" %02x", data[i])

		if i > limit {
			break
		}
	}

	lineStr += "  ||  "

	for i := len(data) - 6; i < len(data); i++ {
		if i >= 0 {
			lineStr += fmt.Sprintf(" %02x", data[i])
		}
	}

	fmt.Println(lineStr)

}

func PrintBinStr(data []byte, limit int, str string) {

	lineStr := ""

	if len(data) > 2 {
		lineStr += fmt.Sprintf(str+" -data: %5d ", len(data))
	}

	for i := 0; i < len(data); i++ {

		if i%4 == 0 {
			lineStr += " "
		}
		if i%8 == 0 {
			lineStr += "  "
		}

		lineStr += fmt.Sprintf(" %02x", data[i])

		if i > limit {
			break
		}
	}

	lineStr += "  ||  "

	for i := len(data) - 6; i < len(data); i++ {

		if i >= 0 {
			lineStr += fmt.Sprintf(" %02x", data[i])
		}
	}

	fmt.Println(lineStr)

}

func PrintfBin(data []byte, str string) string {
	lineStr := ""

	if len(data) > 2 {
		lineStr += fmt.Sprintf(str+"-data: %5d ", len(data))
	}

	for i := 0; i < len(data); i++ {

		if i%4 == 0 {
			//lineStr += " "
		}
		if i%16 == 0 {
			lineStr += "  "
		}

		lineStr += fmt.Sprintf(" %02x", data[i])
	}

	return lineStr + "\r\n"
}
