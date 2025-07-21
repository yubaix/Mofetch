// Whole idea for this file was from Ashish Kumar's Mufetch
package cmd

import (
	"fmt"
	"os"

	"github.com/yubaix/mofetch/config"
	"github.com/yubaix/mofetch/pkg/display"
	"github.com/yubaix/mofetch/pkg/omdb"

	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	omdbClient *omdb.Api
	imageSize  int
	Verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mofetch",
	Short: "Mofetch is a Go-based neofetch-like tool for movies",
	Long:  `Mofetch is a command-line tool written in Go that fetches and displays movie information using the OMDB API.`,
}

// searchCmd is the search command
var searchCmd = &cobra.Command{
	Use:   "search [movie title]",
	Short: "Search for a movie (note that dashes (-) must be used instead of spaces)",
	Long:  `Search for a movie by title and display its information.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		if !config.HasCredentials() {
			fmt.Println("No API key found. Please set your OMDB API key using 'mofetch auth'.")
			os.Exit(1)
		}

		// Load configuration
		var err error
		cfg, err = config.GetConfig()
		if err != nil {
			fmt.Printf("Error reading configuration: %v\n", err)
			os.Exit(1)
		}

		// Initialize OMDB client with API key
		omdbClient = omdb.NewOMDBClient(cfg.OmdbApiKey)

		// Perform Search
		if query != "" {
			searchomdb(query, Verbose)
		} else {
			fmt.Println("Please provide a movie title to search for.")
			os.Exit(1)
		}

		fmt.Print("\033[?25l")
		defer fmt.Print("\033[?25h")

		fmt.Printf("\n")
	},
}

func searchomdb(query string, verbose bool) {
	// Perform the search using the OMDB client
	if result, err := omdbClient.Search(query); err == nil {
		fmt.Printf("\n")
		if verbose { // Determines whether to display film information in verbose mode or not
			display.DisplayFilmVerbose(*result, *omdbClient, imageSize)
		} else {
			display.DisplayFilm(*result, *omdbClient, imageSize)
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately
// Thanks, Ashish :)
func Execute() {
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Thanks, Ashish :)
func init() {
	// Disable the default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	searchCmd.Flags().IntVarP(&imageSize, "size", "s", 20, "Image size (20-50)")
	searchCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")

	// Add the search command to the root command
	rootCmd.AddCommand(searchCmd)
}
