package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/db"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/models"
	_orgRepo "github.com/imtanmoy/authn/organization/repository"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	"github.com/imtanmoy/logx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createUserCmd)
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create User Command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a user argument")
		}
		if args[0] == "su" {
			return nil
		}
		return fmt.Errorf("invalid user specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		logx.Info(args)
		err := db.InitDB()
		if err != nil {
			logx.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logx.Info("Database Initiated...")
		switch args[0] {
		case "su":
			createSuperUser(db.DB)
			break
		default:
			logx.Errorf("invalid user specified: %s", args[0])
		}
	},
}

func createSuperUser(storage *pg.DB) {
	orgRepo := _orgRepo.NewRepository(storage)
	userRepo := _userRepo.NewRepository(storage)

	authxConfig := authx.AuthxConfig{
		SecretKey:             config.Conf.JWT_SECRET_KEY,
		AccessTokenExpireTime: config.Conf.JWT_ACCESS_TOKEN_EXPIRES,
	}

	au := authx.New(userRepo, &authxConfig)

	ctx := context.Background()

	var org models.Organization
	org.Name = "Example Organization"
	o, err := orgRepo.Save(ctx, &org)
	if err != nil {
		logx.Fatalf("invalid user specified: %s", err.Error())
	}

	if o != nil {
		hashedPassword, err := au.HashPassword("password")
		if err != nil {
			logx.Fatalf("could not create super user, reason %s", err.Error())
		}
		var u models.User
		u.Name = "Super User"
		u.Email = "su@gmail.com"
		u.Password = hashedPassword

		err = userRepo.Save(ctx, &u)
		if err != nil {
			logx.Fatalf("could not create super user, reason %s", err.Error())
		}
		logx.Infof("Super User Created \nEmail: %s\nPassword: %s", u.Email, "password")
	}
}
