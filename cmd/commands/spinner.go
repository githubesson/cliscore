package commands

import (
	"fmt"
)

type SpinnerCommand struct{}

func (c *SpinnerCommand) Name() string {
	return "spinner"
}

func (c *SpinnerCommand) Description() string {
	return "Show available spinner styles"
}

func (c *SpinnerCommand) Execute(args []string) error {
	fmt.Println("Available spinner styles:")
	fmt.Println()
	fmt.Println("🎨 Visual Styles:")
	fmt.Println("  default    - Braille dots (⠋⠙⠹⠸⠼)")
	fmt.Println("  dots       - Heavy dots (⣾⣽⣻⢿⡿⣟⣯⣷)")
	fmt.Println("  arrows     - Rotating arrows (←↖↑↗→↘↓↙)")
	fmt.Println("  bounce     - Bouncing dots (⠁⠂⠄⡀⢀⠠⠐⠈)")
	fmt.Println("  pulse      - Pulsing bar (▁▂▃▄▅▆▇█)")
	fmt.Println("  braille    - Braille pattern (⠋⠙⠚⠒⠂⠒⠲⠴)")
	fmt.Println()
	fmt.Println("😀 Emoji Styles:")
	fmt.Println("  emoji      - Lightning and refresh (⚡🔄⏳)")
	fmt.Println("  planet     - Rotating Earth (🌍🌎🌏)")
	fmt.Println("  clock      - Clock faces (🕐🕑🕒...🕛)")
	fmt.Println()
	fmt.Println("📝 Text Styles:")
	fmt.Println("  simple     - Progressive dots (....)")
	fmt.Println("  text       - Loading bar ([=   ])")
	fmt.Println("  matrix     - Matrix characters (ｱｲｳｴｵ)")
	fmt.Println()
	fmt.Println("🔧 Configuration:")
	fmt.Println("  none       - No spinner")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cliscore setup  - Configure spinner style")
	fmt.Println("  or set CLISCORE_SPINNER_STYLE environment variable")
	
	return nil
}