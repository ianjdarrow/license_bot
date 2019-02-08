package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	args := parseArgs()
	c := &client{org: args.org}
	c.setAuthToken()

	switch args.action {
	case "licenses":
		repos := c.getAllLicenses()
		good := 0.0
		for _, repo := range repos {
			if licensesAreGood(repo.ObservedLicenses) {
				good++
				fmt.Printf("%s %s\n", repo.FullName, color.GreenString(strings.Join(repo.ObservedLicenses, ", ")))
				continue
			}
			if repo.ObservedLicenses[0] == "NONE" {
				fmt.Printf("%s %s\n", repo.FullName, color.RedString("NONE"))
				continue
			}

			fmt.Printf("%s %s (%s)\n", repo.FullName, color.YellowString(strings.Join(repo.ObservedLicenses, ", ")), repo.License.Name)
		}
		coverage := good / float64(len(repos)) * 100
		var result string
		if coverage < 60 {
			result = color.RedString("%.1f", coverage)
		}
		if coverage >= 60 && coverage < 95 {
			result = color.YellowString("%.1f", coverage)
		}
		if coverage >= 95 {
			result = color.GreenString("%.1f", coverage)
		}
		fmt.Printf("%s license coverage: %s%%\n", c.org, result)
	case "contributors":
		contributors := c.getAllContributors()
		for _, con := range contributors {
			fmt.Printf("• %s (%d)\n", con.Login, con.Contributions)
		}
		fmt.Printf("%d total contributors", len(contributors))
	}
}
