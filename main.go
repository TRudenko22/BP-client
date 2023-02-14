package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

const (
	SERVER = "http://anton:8089"
)

// Entry Object
type Entry struct {
	Burpees int `json:"burpees"`
	Minutes int `json:"mins"`
}

func (e *Entry) New(amount, minutes string) {
	var err error
	e.Burpees, err = strconv.Atoi(amount)
	e.Minutes, err = strconv.Atoi(minutes)

	if err != nil {
		panic(err)
	}
}

func (e Entry) Debug() string {
	return fmt.Sprintf("Entry -> Burpees: %v | Minutes: %v", e.Burpees, e.Minutes)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// Commands
func Test(c *cli.Context) {
	client := http.Client{}
	request, err := http.NewRequest("GET", SERVER+"/workout/", nil)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	text, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(text))
}

func Create(c *cli.Context) {
	var entry Entry

	client := &http.Client{}

	entry.New(c.Args().Get(0), c.Args().Get(1))

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}

  fmt.Println(string(jsonEntry))

	request, err := http.NewRequest("POST", SERVER+"/workout/", bytes.NewBuffer(jsonEntry))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Created")
	}

	defer request.Body.Close()

	print(response.Status)

}

func main() {
	app := &cli.App{
		Commands: []cli.Command{
			{
				Name:    "create",
				Usage:   "Creates and entry [burpees, amount]",
				Aliases: []string{"c"},
				Action:  Create,
			},
			{
				Name:   "test",
				Usage:  "Tests server connectivity",
				Action: Test,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}
