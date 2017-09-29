package main

type app struct {
	dummy struct {
		enable bool
		target string
	}
	email struct {
		host, username, password string
		port                     int
	}
}
