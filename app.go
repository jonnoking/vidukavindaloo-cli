package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli"

	"github.com/joho/godotenv"
	fpl "github.com/jonnoking/vidukavindaloo-fpl"
	"github.com/jonnoking/vidukavindaloo-fpl/config"
)

var FPL *fpl.FPL

func main() {
	app := cli.NewApp()
	app.Name = "VV Fantasy Premier League CLI"
	app.Version = "0.0.0.1"

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "player",
			Aliases:     []string{"plyr"},
			Category:    "Players",
			Usage:       "get a player",
			UsageText:   "",
			Description: "",
			ArgsUsage:   "[]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "The player's full name",
				},
			},
			Action: playerSearch,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func playerSearch(c *cli.Context) error {

	player := c.String("name")
	if player == "" {
		log.Fatal("No player name supplied")
	}

	p, err := FPL.Bootstrap.Players.GetPlayerByFullName(player)
	if err != nil {
		log.Fatalf("Could not find the player %s\n", player)
	}

	nc := float32(p.NowCost)
	ncd := nc / 10

	fmt.Println()
	fmt.Printf("%s (%s)\n", p.GetFullName(), p.GetTeam(FPL.Bootstrap.Teams).Name)
	fmt.Printf("* Position: %s\n", p.GetPlayerType(FPL.Bootstrap.PlayerTypes).SingularName)
	fmt.Printf("* Total Points: %d\n", p.TotalPoints)
	fmt.Printf("* Now Cost: £%.1fm\n", ncd)
	if p.News != "" {
		fmt.Printf("* News: %s\n", p.News)
	}
	fmt.Println()

	return nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	fplConfig := config.New(getEnv("FPL_USER", ""), getEnv("FPL_PASSWORD", ""), 8, "", "", "")
	FPL = fpl.New(fplConfig)
	FPL.LoadBoostrapLive()
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
