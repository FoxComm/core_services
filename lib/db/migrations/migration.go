package migrations

import (
	"fmt"
	"os"

	feature_manager_migration "github.com/FoxComm/core_services/feature_manager/migration"
	user_migration "github.com/FoxComm/core_services/user/migration"
	"github.com/FoxComm/core_services/user/service"
	"github.com/FoxComm/libs/endpoints"
)

func RunAllMigrations() {
	feature_manager_migration.Run()

	userMigration := user_migration.Migration{}
	if err := userMigration.InitializeForFeature(endpoints.UserAPI.Name); err != nil {
		fmt.Println("Failed to initialize migration for the users feature")
		fmt.Printf("Error was: %s\n", err.Error())
		os.Exit(1)
	}

	userMigration.Run()

	u := &user.User{}
	if err := u.InitializeForFeature(endpoints.UserAPI.Name); err != nil {
		fmt.Println("Failed to initialize migration for the users feature")
		fmt.Printf("Error was: %s\n", err.Error())
		os.Exit(1)
	}

	u.Where(user.User{Email: "admin@wearebeautykind.com"}).Assign("Role", "admin").FirstOrCreate(u)
	u.UpdatePassword("123qwe123")
}
