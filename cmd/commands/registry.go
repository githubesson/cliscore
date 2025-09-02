package commands

// Command interface defines the structure for all CLI commands
type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

// GetCommands returns all available commands
func GetCommands() []Command {
	return []Command{
		&SearchCommand{},
		&CountCommand{},
		&SetupCommand{},
		&ConfigCommand{},
		&MachineInfoCommand{},
		&DownloadCommand{},
		&CreditsCommand{},
		&SpinnerCommand{},
	}
}