package main

import (
	"flag"
	"log"
	"os"

	"github.com/gotoailab/trendhub/web"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to config file")
	keywordPath := flag.String("keywords", "config/frequency_words.txt", "Path to keywords file")
	webMode := flag.Bool("web", false, "Run in web mode")
	webAddr := flag.String("addr", ":8080", "Web server address")
	flag.Parse()

	runner := web.NewTaskRunner(*configPath, *keywordPath)

	if *webMode {
		server := web.NewServer(runner)
		log.Fatal(server.Run(*webAddr))
	} else {
		runner.ExtraWriter = os.Stdout
		if _, err := runner.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
