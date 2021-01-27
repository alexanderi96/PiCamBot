package main

import (
	"github.com/NicoNex/echotron"
	"os/exec"
	"time"
	"os"
	"log"
	"strconv"
)

type bot struct {
	chatId int64
	echotron.Api
}

func newBot(api echotron.Api, chatId int64) echotron.Bot {
    return &bot{
            chatId,
            api,
        }
}

func (b *bot) Update(update *echotron.Update) {
	path := "./archive/" + strconv.FormatInt(b.chatId, 10) + "/"
	b.checkFolder(path)

	log.Println("Message recieved from: " + strconv.FormatInt(b.chatId, 10))
	if update.Message.Text == "/start" {
	        b.SendMessage("Ready to take a shot!", b.chatId)
	} else if update.Message.Text == "/shot" {
		b.SendMessage("Taking a shot", b.chatId)
		date := time.Now().Unix()
		strdate := strconv.FormatInt(date, 10)
		name := "pic" + strdate + ".jpg"
		_, err := exec.Command("raspistill", "-o",  path + name).Output()
		if err == nil {
			b.SendDocument(path + name, name, b.chatId)
		} else {
			log.Fatal(err)
			b.SendMessage("error", b.chatId)
		}
	}
}

func (b *bot) checkFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		
		os.Mkdir(path, 0755)
		log.Println("Directory for user " + strconv.FormatInt(b.chatId, 10) + " created")
	}
}

func main() {
	if _, err := os.Stat("archive"); os.IsNotExist(err) {
		os.Mkdir("archive", 0755)
	}
	dsp := echotron.NewDispatcher("your-telegram-token", newBot)
	log.Println("Running CamBot")
	dsp.Run()
}

