package config

import "project-app/model"

func DbConfig() model.DbConfig {

	return model.DbConfig{
		Host:     "192.168.53.251",
		Port:     "5432",
		Username: "root",
		Password: "root",
		Database: "db_todolist",
	}
}
