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

	buildMessages(Session)

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

		//find channel with name "game-servers"
		channels, _ := Session.GuildChannels(guild.ID)

		var serverChannel string

		for _, c := range channels {
			if c.Name == "game-servers" {
				serverChannel = c.ID
				//fmt.Printf("%# v", pretty.Formatter(serverChannel))
			}
		}

		for game := range payLoad[guild.ID] {
			fmt.Printf("%# v", pretty.Formatter(payLoad[guild.ID][game].serverStatus))

			var currStatus string

			if payLoad[guild.ID][game].serverStatus == 1 {
				currStatus = "Online"
			} else {
				currStatus = "Offline"
			}

			msg, err := Session.ChannelMessageSend(serverChannel, "\n " + "GAME: " + payLoad[guild.ID][game].gameName + "\nSTATUS " + currStatus)
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("%# v", pretty.Formatter(msg))

		}

	}

}

type gameServer struct {
	gameName string
	nameSpace string
	serverStatus int32
	deploymentName string
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

				currServer := gameServer{deployList.Items[i].Labels["gamename"], deployList.Items[i].Namespace, deployList.Items[i].Status.AvailableReplicas, deployList.Items[i].Name}

				deploySlice = append(deploySlice, currServer)

			}

		}

		deployMap[guildList[guild]] = deploySlice
	}

	return deployMap
}
