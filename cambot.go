package main

import (
	"fmt"
	"os/exec"
	"time"
	"os"
	"log"
	"strconv"

	_"embed"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}


var (
	//go:embed token
	token string

	//go:embed admin
	admin string
)

func init() {
	if len(token) == 0 {
		log.Fatal("Empty token file")
	}

	if len(admin) == 0 {
		log.Fatal("Empty admin file")
	}
}

func newBot(chatID int64) echotron.Bot {
    return &bot{
            chatID,
            echotron.NewAPI(token),
        }
}

func (b *bot) Update(update *echotron.Update) {
	path := "./archive/" + strconv.FormatInt(b.chatID, 10) + "/"
	b.checkFolder(path)

	log.Println("Message recieved from: " + strconv.FormatInt(b.chatID, 10))
	
	if strconv.Itoa(int(b.chatID)) != admin {
		b.SendMessage("📷", b.chatID, nil)
	} else {
		msg := update.Message.Text
		if msg == "/start" {
		        b.SendMessage("Ready to take a shot!", b.chatID, nil)
		} else if msg == "/shot" {
			b.SendMessage("Taking a shot", b.chatID, nil)
			date := time.Now().Unix()
			strdate := strconv.FormatInt(date, 10)
			name := "pic" + strdate + ".jpg"
			_, err := exec.Command("libcamera-still", "--immediate", "1", "--nopreview", "1", "--output",  path + name).Output()
			
			if err != nil {
				b.SendMessage(fmt.Sprintf("%s", err), b.chatID, nil)
				log.Fatal(err)
			}

			b.SendMessage("Sending image", b.chatID, nil)
			b.SendDocument(echotron.NewInputFilePath(path + name), b.chatID, nil)
		} else {
			b.SendMessage("👀", b.chatID, nil)
		}
	}
}

func (b *bot) checkFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		log.Println("Directory for user " + strconv.FormatInt(b.chatID, 10) + " created")
	}
}

func main() {
	if _, err := os.Stat("archive"); os.IsNotExist(err) {
		os.Mkdir("archive", 0755)
	}
	dsp := echotron.NewDispatcher(token, newBot)
	log.Println("Running CamBot")
	log.Println(dsp.Poll())
}

