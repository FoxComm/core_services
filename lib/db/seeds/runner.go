package seeds

import (
	"log"
	"os"

	"gopkg.in/mgo.v2/bson"

	"github.com/FoxComm/core_services/lib/db/migrations"
)

func (sr *SeedRunner) AttachToCommandLine() {
	// If you didn't run main.go with an optional 'seeds' param, let's just skip.
	if len(os.Args) < 2 || os.Args[1] != "seeds" && os.Args[1] != "migrations" {
		return
	}

	log.Printf("The SeedRunner has been attached to the commandLine %+v", os.Args)

	if len(os.Args) == 3 {
		if os.Args[1] == "seeds" {
			//You added the seeds command, let's run them.
			switch os.Args[2] {
			case "all":
				sr.GenerateAllRules()
			case "reset":
				sr.DeleteAllRules()
			case "stores":
				sr.GenerateStores()
			}
		}
		if os.Args[1] == "migrations" {
			migrations.RunAllMigrations()
		}
	}
}

func (sr *SeedRunner) GenerateAllRules() {
	log.Printf("Generating All Rules...")
	sr.GenerateStores()
}
