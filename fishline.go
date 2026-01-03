package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	version := "v1.02"
	publicIP := GetPublicIP()
	port := ConfigValue.Port

	http.HandleFunc("/", PipelineHandler)

	title := figure.NewFigure("Fishline", "smslant", true)
	title.Print()

	fmt.Printf("From Hyvr • %s", version)
	fmt.Print("\n")

	box := box.New(box.Config{Px: 2, Py: 0, Type: "Round", Color: "Yellow"})
	content := fmt.Sprintf("Public:  http://%s:%s\nPrivate: http://127.0.0.1:%s", publicIP, port, port)
	box.Print("Pipeline URLS", content)

	WriteLog("")
	WriteLog("Fishline service is started")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
