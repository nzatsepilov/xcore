package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	errInvalidArgs = errors.New("invalid args")
)

var rootCmd = &cobra.Command{
	Use: "xcore",
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func StartCLI() {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "$ ",
		Templates: templates,
	}
	for {
		str, err := prompt.Run()
		if err != nil {
			log.Println(err)
		}

		rootCmd.SetArgs(strings.Split(str, " "))
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
