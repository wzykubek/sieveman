package main

import (
	"fmt"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"go.wzykubek.xyz/sieveman/cmd"
)

func main() {
	manPage, err := mcobra.NewManPage(1, cmd.Root())
	if err != nil {
		panic(err)
	}

	sections := []struct {
		Name    string
		Content string
	}{
		{Name: "AUTHOR", Content: "Wiktor Zykubek <dev at wzykubek dot xyz>"},
		{Name: "HOMEPAGE", Content: "https://github.com/wzykubek/sieveman"},
		{Name: "LICENSE", Content: "ISC"},
	}

	for _, s := range sections {
		manPage = manPage.WithSection(s.Name, s.Content)
	}

	fmt.Println(manPage.Build(roff.NewDocument()))
}
