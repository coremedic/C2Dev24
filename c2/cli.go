package c2

import (
	"fmt"
	"github.com/chzyer/readline"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var helpMenu string = `
CSC C2 v0.0.1 (2024-03-01)
-------------------------------------------
Main menu commands:
Command                  Description
-------------------------------------------
help                     Show this menu
clear                    Clear the console
agents                   List active agents
exit                     Exit C2
-------------------------------------------

`

// StartCLI starts the pseudo console loop
func StartCLI() {
	// Init readline autoCompleter
	autoCompleter := readline.NewPrefixCompleter(
		readline.PcItem("clear"),
		readline.PcItem("help"),
		readline.PcItem("agents"),
		readline.PcItem("exit"),
	)

	// Init readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "C2 > ",       // Readline prompt to be inserted on every line
		AutoComplete:    autoCompleter, // Add our auto completer object
		InterruptPrompt: "^C",          // "^C" = "Ctrl + C"
		EOFPrompt:       "exit",        // Exit CLI with "exit command"
	})
	if err != nil {
		log.Fatal(err)
	}

	// Defer readline close until function exit
	defer rl.Close()

	// Create history file to log commands
	historyFile, err := os.OpenFile("c2/data/.history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	// Main CLI loop
	for {
		// Read command from CLI
		command, err := rl.Readline()
		if err != nil {
			log.Fatal(err)
		}

		// Write command to history file
		if _, err := historyFile.WriteString(command + "\n"); err != nil {
			fmt.Println(err)
		}

		// Trim spaces from command
		command = strings.TrimSpace(command)

		// Command logic switch
		switch command {
		case "exit": // Exit command logic
			{
				fmt.Println("[C2] Shutting down C2, goodbye...")
				os.Exit(0)
			}
		case "help": // Help command logic
			{
				fmt.Print(helpMenu)
			}
		case "clear": // Clear command logic
			{
				if err := clearConsole(); err != nil {
					fmt.Println(err)
				}
			}
		case "agents": // List agents logic
			{
				for _, agent := range AgentMap.Agents {
					fmt.Printf("ID: %s IP: %s, Last Call: %.0f\n", agent.Id, agent.Ip, time.Since(agent.LastCall).Seconds())
				}
			}
		}
	}
}

// clearConsole clears the console
// OS-agnostic
func clearConsole() error {
	var cmd *exec.Cmd
	// Check which OS we are on
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	// Redirect stdout
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
