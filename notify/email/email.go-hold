/* 2018-12-27 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main


import (
	"fmt"
	"net/smtp"
	"strings"

    yml "gopkg.in/yaml.v2"
	promp "github.com/prometheus/client_golang/prometheus"
)

func createEmailTicket(a *AmgrAlert) (string, error) {

	if *args.Debug {
		fmt.Println("DEBUG: create mock api ticket for: ")
		yout, _ := yml.Marshal(*a)
		fmt.Println(string(yout))
	}

	c, err := smtp.Dial(*args.SMTPAddr)
	if err != nil {
		return "", fmt.Errorf("smtp.Dial: %s %s",*args.SMTPAddr,err)
	}
	defer c.Close()
	c.Mail(*args.EmailFrom)
	c.Rcpt(*args.EmailTo)

	wc, err := c.Data()
	if err != nil {
		return "", fmt.Errorf("smtp.Data: %s %s",*args.SMTPAddr,err)
	}
	defer wc.Close()

	node := strings.Split(a.Labels["instance"],":")[0]
	ybody, err := yml.Marshal(*a)
	if err != nil {
		return "", fmt.Errorf("yml.Marshal: %s\n%+v\n",err.Error(),*a)
	}

	_, err = wc.Write([]byte("Subject: "+node+": "+a.Labels["alertname"]))
	if err != nil {
		return "", fmt.Errorf("smtp.Write: %s",err.Error())
	}
	if _, err = wc.Write(ybody); err != nil {
		return "", fmt.Errorf("smtp.Write: %s",err.Error())
	}

	prom.TicketsGend.With(
		promp.Labels{
			"type": "email",
			"dest": *args.EmailTo}).Inc()
	return "", nil
}
