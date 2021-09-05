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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "Simple CLI",
}

var vertUpCmd = &cobra.Command{
	Use:   "resize",
	Short: "Vertical Upscale",
	Long:  `Do a vertical upscaling by increasing memory and cpu`,
	Run: func(cmd *cobra.Command, args []string) {
		VerticalScaleUp(cores, memory)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status",
	Long:  `Prints the current status of the system`,
	Run: func(cmd *cobra.Command, args []string) {
		Status()
	},
}

func Status() {
	fmt.Println("Hello Status")
}

func VerticalScaleUp(mem, cpu int) {
	fmt.Printf("Scaling up cpu to '%d' and memory to '%d'\n", cores, memory)
}