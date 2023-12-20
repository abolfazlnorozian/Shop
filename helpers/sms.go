package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendSms(verifyCode, phoneNumber string) {
	url := "http://rest.ippanel.com/v1/messages/patterns/send"
	method := "POST"

	payload := fmt.Sprintf(`
  {
	"pattern_code": "wshbtvj694w8uni",
	"originator": "+983000505",
	"recipient": "%s",
	"values": {
	  "verification-code": "%s"
	}
  }
`, phoneNumber, verifyCode)
	// Convert payload to io.Reader
	payloadReader := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "AccessKey wKIm00-HBx_YeLIpJxnChnd61ZqT8Go_89DQ_XL4yxQ=")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
