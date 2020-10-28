package main

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func getAnnounce(client *discordgo.Session, ctx *gin.Context) {
	channelID := os.Getenv("DANNOUNCEAPI_CHANNELID")

	messages, err := client.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		log.Print(err)
		return
	}

	announcements := make([]announcement, 0)
	for _, m := range messages {
		lines := strings.Split(m.Content, "\n")
		if len(lines) <= 1 {
			log.Println("malformed announcement, skipping... (Message ID: " + m.ID + ")")
			continue
		}

		pluginInternalName := lines[0]
		link := lines[len(lines)-1]
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			link = ""
		}

		channel, _ := client.State.Channel(channelID)
		member, err := client.GuildMember(channel.GuildID, m.Author.ID)
		if err != nil {
			log.Print("error:", err)
			continue
		}

		announce := announcement{
			Author: &author{
				Nickname:  member.Nick,
				Username:  m.Author.Username + "#" + m.Author.Discriminator,
				AvatarURL: m.Author.AvatarURL(""),
			},
			Message:            strings.TrimSpace(m.Content[len(pluginInternalName) : len(m.Content)-len(link)]),
			PluginInternalName: pluginInternalName,
			Link:               link,
		}

		announcements = append(announcements, announce)
	}

	ctx.JSON(200, gin.H{
		"announce": announcements,
	})
}
