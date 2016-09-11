package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/BenchR267/goalfred"
)

func main() {
	queryTerms := os.Args[1:]

	resp := goalfred.NewResponse()

	pods := getPods(strings.Join(queryTerms, "%20"))

	for _, pod := range pods {
		resp.AddItem(pod)
	}

	resp.Print()
}

type Pod struct {
	Link      string
	Name      string `json:"id"`
	Summary   string
	Version   string
	Platforms []string
}

func (p Pod) Documentation() string {
	return fmt.Sprintf("http://cocoadocs.org/docsets/%s/%s", p.Name, p.Version)
}

func (p Pod) Item() *goalfred.Item {
	title := fmt.Sprintf("%s (%s)", p.Name, p.Version)
	instruction := fmt.Sprintf("pod '%s', '%s'", p.Name, p.Version)
	i := &goalfred.Item{
		Title:    title,
		Subtitle: p.Summary,
		Arg:      p.Link,
	}
	i.Mod.Cmd = &goalfred.ModContent{
		Arg:      p.Documentation(),
		Subtitle: "Open documentation!",
	}
	i.Mod.Alt = &goalfred.ModContent{
		Arg:      instruction,
		Subtitle: "Copy pod install instructions",
	}

	return i
}

type Response struct {
	Allocations [][]json.RawMessage `json:"allocations"`
}

func getPods(searchTerm string) []Pod {

	url := fmt.Sprintf("https://search.cocoapods.org/api/v1/pods.picky.hash.json?query=%v&ids=20&offset=0&sort=quality", searchTerm)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response Response
	json.Unmarshal(body, &response)

	rawJson := response.Allocations[0][5]

	var pods []Pod
	err = json.Unmarshal(rawJson, &pods)

	return pods
}
