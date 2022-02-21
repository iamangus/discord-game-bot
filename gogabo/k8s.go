// Declare this file to be part of the main package so it can be compiled into
// an executable.
package main

// Import all Go packages required for this file.
import (
	//Shared
	//"container/list"
	"log"
	"time"
	"fmt"
	//"time"

	//K8s
	"context"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"github.com/kr/pretty"

	//Discord
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

// Session is declared in the global space so it can be easily used throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	// Discord Authentication Token
	token := os.Getenv("DG_TOKEN")

    // Verify a Token was provided
	if token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	Session, _ = discordgo.New("Bot " + token)

	// Declare any variables needed later.
    var err error

	//Session.Identify.Intents = discordgo.IntentsGuildMessageReactions
	
	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	log.Printf(`Now running. Press CTRL-C to exit.`)

}

func main() {


	//Session.AddHandler(buildMessages)

	buildMessages(Session)

	Session.AddHandler(reactionRecieved)

	time.Sleep(10 * time.Second)
	
	// Wait for a CTRL-C
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()
	// Exit Normally.
}

func buildMessages(Session *discordgo.Session) {

	guildSlice := make([]string, 0)
	for _, guild := range Session.State.Guilds {
		guildSlice = append(guildSlice, guild.ID)
	}

	//call k8s function to generate deployment info for each guildID in guildList
	payLoad := genServerList(guildSlice)

	for _, guild := range Session.State.Guilds {
		channels, _ := Session.GuildChannels(guild.ID)

		var serverChannel string

		for _, c := range channels {
			if c.Name == "game-servers" {
				serverChannel = c.ID
			}
		}


		for game := range payLoad[guild.ID] {
			//Check to see if a game message has been sent for each game server, if it has, update the struct with said messages ID
			var limit int = 100
			messageHistory, err := Session.ChannelMessages(serverChannel, limit, "", "", "")
			if err != nil {
				log.Printf("error, %s\n", err)
			}
	
			for m := range messageHistory {
				if strings.Contains(messageHistory[m].Author.ID, "933803745647657010") {
					if strings.Contains(messageHistory[m].Embeds[0].Title, payLoad[guild.ID][game].gameName) {
						payLoad[guild.ID][game].messageID = messageHistory[m].ID
					}
					
	
				}
	
			}

			//Checking server status
			var currStatus string
			if payLoad[guild.ID][game].serverStatus == 1 {
				currStatus = "Online"
			} else {
				currStatus = "Offline"
			}

			currentGame := payLoad[guild.ID][game]

			if payLoad[guild.ID][game].messageID == "" {
				sendMessage(currentGame, currStatus, serverChannel, Session)
			} else {
				updateMessage(currentGame, currStatus, serverChannel, Session)
			}

		}

	}

}

type gameServer struct {
	gameName string
	nameSpace string
	serverStatus int32
	deploymentName string
	messageID string
}

func genServerList(guildList []string) map[string][]gameServer {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deployList, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	deployMap := make(map[string][]gameServer)

	for guild := range guildList {

		deploySlice := make([]gameServer, 0)
		
		for i := range deployList.Items {	

			if strings.Contains(deployList.Items[i].Labels["custguild"], fmt.Sprintf("%v", guildList[guild])) {

				currServer := gameServer{deployList.Items[i].Labels["gamename"], deployList.Items[i].Namespace, deployList.Items[i].Status.AvailableReplicas, deployList.Items[i].Name, ""}

				deploySlice = append(deploySlice, currServer)

			}

		}

		deployMap[guildList[guild]] = deploySlice

	}

	return deployMap
}



func sendMessage(gameStruct gameServer, currStatus string, serverChannel string, Session *discordgo.Session) string {

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: "",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Status",
				Value:  currStatus,
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/88b6f848-7f22-44b7-a036-8e65d20a4e3d/dau8yxs-de0458a4-d980-4ad6-845b-556c71353940.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7InBhdGgiOiJcL2ZcLzg4YjZmODQ4LTdmMjItNDRiNy1hMDM2LThlNjVkMjBhNGUzZFwvZGF1OHl4cy1kZTA0NThhNC1kOTgwLTRhZDYtODQ1Yi01NTZjNzEzNTM5NDAucG5nIn1dXSwiYXVkIjpbInVybjpzZXJ2aWNlOmZpbGUuZG93bmxvYWQiXX0.xpMIe6fAJ1wjxuC7987RbqT-jGsXpVI1GoKAgSmP-ps",
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     gameStruct.gameName,
	}


	msg, _ := Session.ChannelMessageSendEmbed(serverChannel, embed)

	Session.MessageReactionsRemoveAll(serverChannel, msg.ID)

	Session.MessageReactionAdd(serverChannel, msg.ID, "🟢")
	Session.MessageReactionAdd(serverChannel, msg.ID, "🔴")

	return msg.ID

}

func updateMessage(gameStruct gameServer, currStatus string, serverChannel string, Session *discordgo.Session) {

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: "",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Status",
				Value:  currStatus,
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/88b6f848-7f22-44b7-a036-8e65d20a4e3d/dau8yxs-de0458a4-d980-4ad6-845b-556c71353940.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7InBhdGgiOiJcL2ZcLzg4YjZmODQ4LTdmMjItNDRiNy1hMDM2LThlNjVkMjBhNGUzZFwvZGF1OHl4cy1kZTA0NThhNC1kOTgwLTRhZDYtODQ1Yi01NTZjNzEzNTM5NDAucG5nIn1dXSwiYXVkIjpbInVybjpzZXJ2aWNlOmZpbGUuZG93bmxvYWQiXX0.xpMIe6fAJ1wjxuC7987RbqT-jGsXpVI1GoKAgSmP-ps",
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     gameStruct.gameName,
	}

	Session.ChannelMessageEditEmbed(serverChannel, gameStruct.messageID, embed)

	Session.MessageReactionAdd(serverChannel, gameStruct.messageID, "🟢")
	Session.MessageReactionAdd(serverChannel, gameStruct.messageID, "🔴")

	reactionsGreen, _ := Session.MessageReactions(serverChannel, gameStruct.messageID, "🟢", 99, "", "")
	for i := range reactionsGreen {
		if strings.Contains(reactionsGreen[i].ID, "933803745647657010") {
			continue
		} else {
			Session.MessageReactionRemove(serverChannel, gameStruct.messageID, "🟢", reactionsGreen[i].ID)
		}

	}

	reactionsRed, _ := Session.MessageReactions(serverChannel, gameStruct.messageID, "🔴", 99, "", "")
	for i := range reactionsRed {
		if strings.Contains(reactionsRed[i].ID, "933803745647657010") {
			continue
		} else {
			Session.MessageReactionRemove(serverChannel, gameStruct.messageID, "🔴", reactionsRed[i].ID)
		}

	}

	//Session.MessageReactionsRemoveAll(serverChannel, gameStruct.messageID)



}

func reactionRecieved(Session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// we need to pass in the channel name

	//log.Printf("We got a reaction!")
	//fmt.Printf("%# v", pretty.Formatter(r))

	msg, _ := Session.ChannelMessage(r.ChannelID, r.MessageID)

	fmt.Printf("%# v", pretty.Formatter(msg.Embeds[0].Title))

	actionReq := r.Emoji.Name

	if actionReq == "🔴" {
		actionReq = "stop"
	} else if actionReq == "🟢" {
		actionReq = "start"
	} else {
		log.Printf("Incorrect emoji! Quitting!")
		return;
	}

	//fmt.Printf("%# v", pretty.Formatter(actionReq))


	//Generating list of guilds
	guildSlice := make([]string, 0)
	for _, guild := range Session.State.Guilds {
		guildSlice = append(guildSlice, guild.ID)
	}

	//call k8s function to generate deployment info for each guildID in guildList
	payLoad := genServerList(guildSlice)

	for game := range payLoad[r.GuildID] {
		if payLoad[r.GuildID][game].gameName == msg.Embeds[0].Title {
			scaleDeployment(payLoad[r.GuildID][game].deploymentName, payLoad[r.GuildID][game].nameSpace, actionReq)
		}
	}

	time.Sleep(10 * time.Second)

	buildMessages(Session)

}

func scaleDeployment(deploymentName string, nameSpace string, actionReq string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	s, err := clientset.AppsV1().Deployments(nameSpace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
    if err != nil {
        log.Fatal(err)
    }

	sc := *s

	if actionReq == "stop" {
		sc.Spec.Replicas = 0
	} else if actionReq == "start" {
		sc.Spec.Replicas = 1
	} else {
		log.Printf("Something went wrong when dtermining deployment scale request")
		return;
	}



    us, err := clientset.AppsV1().Deployments(nameSpace).UpdateScale(context.TODO(), deploymentName, &sc, metav1.UpdateOptions{})
    if err != nil {
        log.Fatal(err)
    }

	log.Println(*us)

}
