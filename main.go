package main

import pez "github.com/pivotalservices/pezinventory/service"

func main() {
	s := pez.NewServer()
	s.Run()
}
