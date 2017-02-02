package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pashutk/mtproto"
	"github.com/robfig/cron"
)

const (
	maxDelayMinutes = 5
)

func goToForestJob(m *mtproto.MTProto) {
	log.Println("goToForest job: started")
	jobRandDelay()
	log.Println("goToForest job: delay ended, sending command")
	err := m.SendMessageToBot("chatwarsbot", "🌲Лес")
	if err != nil {
		log.Printf("goToForest job: error - %s", err)
	} else {
		log.Println("goToForest job: command sent")
	}
}

func korovanJob(m *mtproto.MTProto) {
	log.Println("korovan job: started")
	jobRandDelay()
	log.Println("korovan job: delay ended, sending command")
	err := m.SendMessageToBot("chatwarsbot", "🐫ГРАБИТЬ КОРОВАНЫ")
	if err != nil {
		log.Printf("korovan job: error - %s", err)
	} else {
		log.Println("korovan job: command sent")
	}
}

func defJob(m *mtproto.MTProto) {
	log.Println("def job: started")
	jobRandDelay()
	log.Println("def job: delay ended, sending command")
	err := m.SendMessageToBot("chatwarsbot", "🛡 Защита")
	time.Sleep(time.Duration(10) * time.Second)
	err = m.SendMessageToBot("chatwarsbot", "🇨🇾")
	if err != nil {
		log.Printf("def job: error - %s", err)
	} else {
		log.Println("def job: command sent")
	}
}

func jobRandDelay() {
	delayMinutes := rand.Intn(maxDelayMinutes)
	delaySeconds := rand.Intn(59)
	sleepDuration := time.Duration(delayMinutes)*time.Minute + time.Duration(delaySeconds)*time.Second
	log.Printf("Setting delay 00:%02d:%02d", delayMinutes, delaySeconds)
	time.Sleep(sleepDuration)
}

func registerCronJobs(m *mtproto.MTProto) {
	c := cron.New()

	c.AddFunc("0 15 0,9-23 * * *", func() { goToForestJob(m) })
	c.AddFunc("0 25 1-8/2 * * *", func() { korovanJob(m) })
	c.AddFunc("0 45 3-23/4 * * *", func() { defJob(m) })

	c.Start()
}

func runBot(m *mtproto.MTProto) error {
	log.Printf("ChatWars game bot activated")

	log.Printf("Register cron jobs")
	registerCronJobs(m)

	log.Printf("Jobs registered, infinite loop start")
	<-make(chan bool)

	return nil
}

func usage() {
	fmt.Print("ChatWars game bot.\n\nUsage:\n\n")
	fmt.Print("    ./cwgb <command> [arguments]\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("    auth <phone_number>                auth connection by code\n")
	fmt.Print("    msgToBot <botname> <msgtext>       send message to telegram bot\n")
	fmt.Print("    bot                                run game bot\n")
	fmt.Println()
}

func main() {
	var err error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	commands := map[string]int{"auth": 1, "msgToBot": 2, "bot": 0}
	valid := false
	for k, v := range commands {
		if os.Args[1] == k {
			if len(os.Args) < v+2 {
				usage()
				os.Exit(1)
			}
			valid = true
			break
		}
	}

	if !valid {
		usage()
		os.Exit(1)
	}

	m, err := mtproto.NewMTProto(os.Getenv("HOME") + "/.telegram_go")
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}

	err = m.Connect()
	if err != nil {
		fmt.Printf("Connect failed: %s\n", err)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "auth":
		err = m.Auth(os.Args[2])

	case "msgToBot":
		err = m.SendMessageToBot(os.Args[2], os.Args[3])

	case "bot":
		err = runBot(m)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
