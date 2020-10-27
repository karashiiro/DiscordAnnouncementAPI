package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func main() {
	client, err := discordgo.New("Bot " + os.Getenv("DANNOUNCEAPI_BOT_TOKEN"))
	if err != nil {
		log.Print(err)
	}

	if err = client.Open(); err != nil {
		log.Print(err)
	}

	done := make(chan error, 2) // servers signal when they are done using this channel
	stop := make(chan struct{}) // servers are commanded to stop with this channel

	go func() {
		done <- startAPI(client, stop)
	}()
	go func() {
		done <- startBot(client, stop)
	}()

	stopped := false
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil { // blocks until one server closes
			log.Println("error:", err)
		}
		if !stopped {
			stopped = true
			close(stop) // force both servers to stop
		}
	}
}

func startAPI(client *discordgo.Session, stop <-chan struct{}) error {
	router := gin.Default()

	router.GET("/announce", func(ctx *gin.Context) {
		getAnnounce(client, ctx)
	})

	server := http.Server{
		Addr:    ":9023",
		Handler: router,
	}

	go func() {
		<-stop
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func startBot(client *discordgo.Session, stop <-chan struct{}) error {
	user, err := client.User("@me")
	if err != nil {
		return err
	}
	log.Println("Logged in as", user.Username+"#"+user.Discriminator)

	<-stop

	client.Close()

	return nil
}
