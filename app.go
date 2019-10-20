package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli"

	"github.com/joho/godotenv"
	fpl "github.com/jonnoking/vidukavindaloo-fpl"
	"github.com/jonnoking/vidukavindaloo-fpl/config"
	"github.com/jonnoking/vidukavindaloo-fpl/models"
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
		cli.Command{
			Name:        "team",
			Aliases:     []string{"tm"},
			Category:    "Teams",
			Usage:       "get a team",
			UsageText:   "",
			Description: "",
			ArgsUsage:   "[]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "short-name",
					Value: "",
					Usage: "The team's short name",
				},
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "The team's name",
				},
			},
			Action: teamShortSearch,
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

	fmt.Println()
	printPlayer(&p)
	fmt.Println()

	return nil
}

func teamShortSearch(c *cli.Context) error {

	team := c.String("short-name")
	if team == "" {
		log.Fatal("No short name supplied")
	}

	t, err := FPL.Bootstrap.Teams.GetTeamByShortName(team)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	fmt.Println()
	fmt.Printf("%s (%s)\n", t.Name, t.ShortName)
	fmt.Printf("* Strength: %d\n", t.Strength)
	fmt.Println()

	g, d, m, f := t.GetSquad(FPL.Bootstrap.Players)

	printPosition(g, "Goalkeepers", false)
	fmt.Println()
	printPosition(d, "Defenders", false)
	fmt.Println()
	printPosition(m, "Midfielders", false)
	fmt.Println()
	printPosition(f, "Forwards", false)
	fmt.Println()

	return nil
}

func printPosition(pos map[int]models.Player, position string, details bool) {
	fmt.Printf("### %s ###\n", position)
	for _, ff := range pos {
		if details {
			printPlayerDetails(&ff)
		} else {
			printPlayer(&ff)
		}
	}

}

func printPlayer(p *models.Player) {
	nc := float32(p.NowCost)
	ncd := nc / 10
	fmt.Printf("%s [£%.1fm] (%s)\n", p.GetFullName(), ncd, p.GetTeam(FPL.Bootstrap.Teams).Name)
}

func printPlayerDetails(p *models.Player) {
	nc := float32(p.NowCost)
	ncd := nc / 10

	fmt.Printf("%s (%s)\n", p.GetFullName(), p.GetTeam(FPL.Bootstrap.Teams).Name)
	fmt.Printf("* Position: %s\n", p.GetPlayerType(FPL.Bootstrap.PlayerTypes).SingularName)
	fmt.Printf("* Total Points: %d\n", p.TotalPoints)
	fmt.Printf("* Now Cost: £%.1fm\n", ncd)
	if p.News != "" {
		fmt.Printf("* News: %s\n", p.News)
	}
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
