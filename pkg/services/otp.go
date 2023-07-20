package services

import (
    "os"
    "net/http"    
    "fmt"
    "io/ioutil"
)

func SendOTP(phone,msg,tid string) error{
    client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://sapteleservices.com/SMS_API/sendsms.php", nil)
	if err != nil {
        return err
	}

	q := req.URL.Query()
	q.Add("username", os.Getenv("OTP_USERNAME"))
	q.Add("password", os.Getenv("OTP_PASSWORD"))
	q.Add("mobile", phone)
	q.Add("sendername", os.Getenv("OTP_SENDER"))
	q.Add("message", msg)
	q.Add("routetype", "1")
	q.Add("tid", os.Getenv(tid))

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))

    return nil 
}
