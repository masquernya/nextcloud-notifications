package main

import (
	"github.com/masquernya/nextcloud-notifications/cloud"
	logger "github.com/masquernya/nextcloud-notifications/log"
	"github.com/masquernya/nextcloud-notifications/storage"
	"time"
)

var log = logger.New("main")

func main() {

	c := cloud.New()
	// Normal user
	if c.IsLoginRequired() {
		log.Info("Login required")
		data := c.RequestLogin()
		log.Info("Copy and paste this URL into your web browser, and login with your normal NextCloud account:\n", data.LoginUrl)
		result, err := c.PollLogin(data.Poll)
		if err != nil {
			panic(err)
		}
		storage.Get().LoginUsername = result.Username
		storage.Get().LoginPassword = result.Password
		storage.Save()
		log.Info("Login data saved successfully")
	}
	// Admin
	if storage.Get().AdminUsername == "" || storage.Get().AdminPassword == "" {
		log.Info("Admin login required")
		data := c.RequestLogin()
		log.Info("Copy and paste this URL into your web browser, and login with your admin NextCloud account:\n", data.LoginUrl)
		result, err := c.PollLogin(data.Poll)
		if err != nil {
			panic(err)
		}
		storage.Get().AdminUsername = result.Username
		storage.Get().AdminPassword = result.Password
		storage.Save()
		log.Info("Admin saved successfully")
	}
	log.Info("Continuing...")
	for {
		c.SendNotifications()
		log.Info("Sleeping for 1 minute")
		time.Sleep(time.Minute * 1)
	}
}
