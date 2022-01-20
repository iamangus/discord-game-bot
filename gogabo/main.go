// Declare this file to be part of the main package so it can be compiled into
// an executable.
package main

// Import all Go packages required for this file.
import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {
	// Discord Authentication Token
	Session.Token = os.Getenv("DG_TOKEN")
}

func main() {

	// Declare any variables needed later.
	var err error

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}
	
	session.AddHandler(message)

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}

func message(bot *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot { return }
	switch {
	case strings.HasPrefix(message.Content, config.BotPrefix):
		ping := bot.HeartbeatLatency().Truncate(60)
		if message.Content == "&ping" {
			bot.ChannelMessageSend(message.ChannelID,`My latency is **` + ping.String() + `**!`)
		}
		if message.Content == "&author" {
			bot.ChannelMessageSend(message.ChannelID, "My author is Gonz#0001, I'm only a template discord bot made in golang.")
		}
		if message.Content == "&github" {
			embed := embed.NewEmbed().
				SetAuthor(message.Author.Username, message.Author.AvatarURL("1024")).
				SetThumbnail(message.Author.AvatarURL("1024")).
				SetTitle("My repository").
				SetDescription("You can find my repository by clicking [here](https://github.com/gonzyui/Discord-Template).").
				SetColor(0x00ff00).MessageEmbed
			bot.ChannelMessageSendEmbed(message.ChannelID, embed)
		}
		if message.Content == "&botinfo" {
			guilds := len(bot.State.Guilds)
			embed := embed.NewEmbed().
				SetTitle("My informations").
				SetDescription("Some informations about me :)").
				AddField("GO version:", runtime.Version()).
				AddField("DiscordGO version:", discordgo.VERSION).
				AddField("Concurrent tasks:", strconv.Itoa(runtime.NumGoroutine())).
				AddField("Latency:", ping.String()).
				AddField("Total guilds:", strconv.Itoa(guilds)).MessageEmbed
			bot.ChannelMessageSendEmbed(message.ChannelID, embed)
		}
	}
}
