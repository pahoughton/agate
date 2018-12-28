/* 2018-12-27 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main


import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

    yml "gopkg.in/yaml.v2"
	promp "github.com/prometheus/client_golang/prometheus"
)

func createEmailTicket(a *AmgrAlert){

	if *args.Debug {
		fmt.Println("DEBUG: create mock api ticket for: ")
		yout, _ := yml.Marshal(*a)
		fmt.Println(string(yout))
	}

	c, err := smtp.Dial(*args.SMTPAddr)
	if err != nil {
		fmt.Println("FATAL-smtp.Dial: ",*args.SMTPAddr," ",err)
		os.Exit(2)
	}
	defer c.Close()
	c.Mail(*args.EmailFrom)
	c.Rcpt(*args.EmailTo)

	wc, err := c.Data()
	if err != nil {
		fmt.Println("FATAL-smtp.Data: ",*args.SMTPAddr," ",err)
		os.Exit(2)
	}
	defer wc.Close()

	node := strings.Split(a.Labels["instance"],":")[0]
	ybody, err := yml.Marshal(*a)
	if err != nil {
		fmt.Printf("FATAL-yml.Marshal: %s\n%+v\n",err.Error(),*a)
		os.Exit(2)
	}

	_, err = wc.Write([]byte("Subject: " +
		node + ": " + a.Labels["alertname"]))
	if err != nil {
		fmt.Println("FATAL-smtp.Write: ",err.Error())
		os.Exit(2)
	}
	if _, err = wc.Write(ybody); err != nil {
		fmt.Println("FATAL-smtp.Write: ",err.Error())
		os.Exit(2)
	}

	prom.TicketsGend.With(
		promp.Labels{
			"type": "email",
			"dest": *args.EmailTo}).Inc()
}
