// Pretty much this entire file is a copy of Ashish Kumar's Mufetch config.go file
package cmd

import (
	"fmt"
	"os"

	"github.com/yubaix/mofetch/config"

	"github.com/spf13/cobra"
)

// authCmd represents the authentication command for Spotify API
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with OMDB API",
	Long: `Set up your OMDB API credentials.

You need to:
1. Go to https://developer.spotify.com/dashboard
2. Create a new app
3. Copy your Client ID and Client Secret`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Spotify API Authentication Setup")
		fmt.Println()
		fmt.Println("To get your OMDB API credentials:")
		fmt.Println("1. Go to: https://www.omdbapi.com/apikey.aspx")
		fmt.Println("2. Select 'FREE! (1,000 daily limit)'")
		fmt.Println("3. Enter in your email and name, put anything in for 'Use'")
		fmt.Println("4. Click Submit")
		fmt.Println("5. Check your email for the API key")
		fmt.Println()

		var apiKey string

		fmt.Print("Enter your OMDB API Key: ")
		fmt.Scanln(&apiKey)

		// Validate credentials
		if apiKey == "" {
			fmt.Println("Error: API Key is required.")
			os.Exit(1)
		}
		if len(apiKey) < 8 {
			fmt.Println("Warning: API Key seems too short. Please verify they are correct.")
		}

		// Set credentials in config
		if err := config.SetCredentials(apiKey); err != nil {
			fmt.Printf("Failed to save API Key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Credentials saved successfully!")
		fmt.Println("You can now use 'mufetch search <query>' to search for music.")
	},
}

// init adds the auth command to the root command
func init() {
	rootCmd.AddCommand(authCmd)
}
