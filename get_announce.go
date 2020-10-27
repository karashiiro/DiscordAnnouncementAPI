package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func getAnnounce(client *discordgo.Session, ctx *gin.Context) {
	messages, err := client.ChannelMessages(os.Getenv("DANNOUNCEAPI_CHANNELID"), 100, "", "", "")
	if err != nil {
		log.Print(err)
		return
	}

	announcements := make([]announcement, 0)
	for _, m := range messages {
		pluginInternalName := strings.Split(m.Content, "\n")[0]
		if strings.Index(pluginInternalName, " ") == -1 {
			log.Println("Malformed announcement (Message ID: " + m.ID + ")")
			continue
		}

		announce := announcement{
			Author:             m.Member.Nick,
			Message:            m.Content[len(pluginInternalName)-1:],
			PluginInternalName: pluginInternalName,
		}

		if announce.Author == "" {
			announce.Author = m.Author.Username
		}

		announcements = append(announcements, announce)
	}

	payload, err := json.Marshal(announcements)
	if err != nil {
		log.Print("error:", err)
		payload = []byte("[]")
	}

	ctx.JSON(200, gin.H{
		"announce": string(payload),
	})
}
