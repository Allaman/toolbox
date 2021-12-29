package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cores  int
	memory int
)

func main() {
	vertUpCmd.Flags().IntVarP(&cores, "cores", "c", 1, "number of cores")
	vertUpCmd.Flags().IntVarP(&memory, "memory", "m", 2, "amount of memory")
	vertUpCmd.MarkFlagRequired("cores")
	vertUpCmd.MarkFlagRequired("memory")
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(vertUpCmd)
	rootCmd.AddCommand(CompletionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  "To load completions",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "simple-cobra",
	Short: "Server status and modify",
}

var vertUpCmd = &cobra.Command{
	Use:   "resize",
	Short: "Vertical Upscale",
	Long:  `Do a vertical upscaling by increasing memory and cpu`,
	Run: func(cmd *cobra.Command, args []string) {
		verticalScaleUp(cores, memory)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status",
	Long:  `Prints the current status of the system`,
	Run: func(cmd *cobra.Command, args []string) {
		status()
	},
}

func status() {
	fmt.Println("Hello Status")
}

func verticalScaleUp(mem, cpu int) {
	fmt.Printf("Scaling up cpu to '%d' and memory to '%d'\n", cores, memory)
}
