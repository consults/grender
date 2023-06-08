package model

type Task struct {
	Url     string
	Xpath   string
	TimeOut int
	Cookies []Cook
}
type Cook struct {
	Name   string
	Value  string
	Domain string
}
