package bot

import (
	"github.com/CoyoTan/JACK-AL/bot/responders"
	"github.com/CoyoTan/JACK-AL/structs"
	"github.com/bwmarrin/discordgo"
	"github.com/txgruppi/parseargs-go"
	"os"
)

var (
	jackal *structs.CoreCfg
)

//Init accepts a pointer to the core of Jackal. This core will be used to establish a connection with Discord and act as a registration point for command handlers, as well as a service provider for other linked services.
func Init(core *structs.CoreCfg) {

	jackal = core

	jackal.Logger.Info.Println("Main process has handed off &Core to package: bot")

	dgOpen()

	for {

		resp, _ := parseargs.Parse(GetInput())

		if len(resp) > 0 {
			switch resp[0] {
			case "ping":
				jackal.Logger.Console.Println("Pong")
			case "leave":
				jackal.Discord.Session.Close()
				os.Exit(100)
			}
		}
	}
}

func dgOpen() {

	if !(len(jackal.Discord.Token) > 0) {
		jackal.Logger.Error.Println("The configuration provided does not contain an API Token. Please provide a token to the Jackal Configuration file.", len(jackal.Discord.Token))
	}

	dg, err := discordgo.New("Bot " + jackal.Discord.Token)

	dg.StateEnabled = true
	dg.State.MaxMessageCount = jackal.Discord.MaxMessageCount

	if err != nil {
		jackal.Logger.Error.Println("There was an error when attempting to begin a session with Discord.")
	}

	responders.InitAll(jackal)

	dg.AddHandler(ready)
	dg.AddHandler(createDispatch)

	err = dg.Open()

	if err != nil {
		jackal.Logger.Error.Println(err)
		os.Exit(101)
	}
}
