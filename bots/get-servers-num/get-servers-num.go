package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/go-homedir"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

const (
	hookString = "get-servers-num"
)

var (
	botId string
)

type Bot struct {
	api *slack.Client
	rtm *slack.RTM
}

func NewBot(token string) *Bot {
	bot := new(Bot)
	bot.api = slack.New(token)
	bot.rtm = bot.api.NewRTM()
	return bot
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".infra-bot")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		fmt.Println("Could not read config file: ", err)
		os.Exit(1)
	}
}

func numberOfInstances() int {
	var instanceNum = 0

	// ローカル確認時用
	// ec2svc := ec2.New(session.New(&aws.Config{
	// 	Region:      aws.String("ap-northeast-1"),
	// 	Credentials: credentials.NewSharedCredentials("", "readyfor"),
	// }))

	sess := session.Must(session.NewSession())
	ec2svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Env"),
				Values: []*string{aws.String("prd")},
			},
			{
				Name:   aws.String("tag:Act_Web"),
				Values: []*string{aws.String("True")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	for _, res := range resp.Reservations {
		instanceNum += len(res.Instances)
	}

	return instanceNum
}

func main() {
	token := viper.GetString("token")

	bot := NewBot(token)
	go bot.rtm.ManageConnection()

	for {
		select {
		case msg := <-bot.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
			case *slack.HelloEvent:
				log.Printf("start up.")
			case *slack.MessageEvent:
				command := strings.Replace(ev.Text, "<@"+botId+">", "", 1)
				fmt.Println(ev.Text)
				fmt.Println(command)
				if strings.Contains(command, hookString) && strings.HasPrefix(ev.Text, "<@"+botId+">") {
					num := numberOfInstances()
					bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("現在の Web サーバー台数は... "+strconv.Itoa(num)+"台です！！！１", ev.Channel))
				}
			}
		}
	}
}
