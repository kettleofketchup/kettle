package terminal

import (
	"fmt"

	"github.com/kettleofketchup/kettle/src/cmd/helpers"
	"github.com/spf13/cobra"
)

// StarshipCmd represents the starship command
var StarshipCmd = &cobra.Command{
	Use:   "starship",
	Short: "Starship cross-shell prompt commands",
	Long:  `Install and configure Starship, a fast, customizable cross-shell prompt.`,
}

var starshipInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Starship cross-shell prompt",
	Long:  `Downloads and installs Starship using the official installation script from starship.rs.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := installStarship()
		if err != nil {
			helpers.PrintError("Failed to install Starship", err)
		} else {
			helpers.PrintSuccess("Starship installed successfully!")
			helpers.PrintInfo("Starship initialization has been added to your kettle shell profile")
		}
	},
}

func installStarship() error {
	// Check if starship is already installed
	if helpers.CommandExists("starship") {
		if !helpers.PromptYesNo("Starship is already installed. Do you want to reinstall it?") {
			return nil
		}
	}

	helpers.PrintInfo("Downloading Starship installation script...")

	// Use the helpers package to download and run the install script
	url := "https://starship.rs/install.sh"
	filename := "starship_install.sh"
	if err := helpers.DownloadFile(filename, url); err != nil {

		return fmt.Errorf("failed to download Starship install script: %w", err)
	}

	cmd := `sh -c "./starship_install.sh -b ~/.local/bin -f" `

	if err := helpers.RunCmd(cmd); err != nil {
		return fmt.Errorf("failed to run Starship install script: %w", err)
	}

	// Add starship init to shell profile
	addStarshipToShellProfile()

	return nil
}

func addStarshipToShellProfile() {
	shell := helpers.GetCurrentShell()
	var initLine string

	switch shell {
	case "bash":
		initLine = `eval "$(starship init bash)"`
	case "zsh":
		initLine = `eval "$(starship init zsh)"`
	case "fish":
		initLine = `eval "$(starship init fish)"`
	default:
		helpers.PrintInfo(fmt.Sprintf("Manual setup required for shell: %s", shell))
		helpers.PrintInfo("Add the appropriate starship init command to your shell profile")
		return
	}

	// Add to kettle shell profile instead of main profile for better organization
	added := helpers.AddLineToShellProfile(initLine)
	if added {
		helpers.PrintSuccess(fmt.Sprintf("Added Starship initialization to %s profile", shell))
		helpers.PrintInfo("Restart your shell or source your profile to activate Starship")
	} else {
		helpers.PrintInfo(fmt.Sprintf("Starship initialization already present in %s profile", shell))
	}
}

func init() {
	StarshipCmd.AddCommand(starshipInstallCmd)
}
