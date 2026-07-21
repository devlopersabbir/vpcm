package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// serverNameCompletions returns all registered server names from the local
// inventory database. Used as a ValidArgsFunction for commands that accept a
// server name or id as a positional argument.
func serverNameCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	repo, _, err := initRepoAndService(context.Background())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	servers, err := repo.List(context.Background())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, s := range servers {
		names = append(names, s.Name+"\t"+s.Username+"@"+s.Host)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// serverHostCompletions returns all registered server host IPs from the local
// inventory database. Used as a flag completion function for --host flags.
func serverHostCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	repo, _, err := initRepoAndService(context.Background())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	servers, err := repo.List(context.Background())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var hosts []string
	for _, s := range servers {
		hosts = append(hosts, s.Host+"\t"+s.Name)
	}
	return hosts, cobra.ShellCompDirectiveNoFileComp
}

// completionCmd generates the shell completion script for bash, zsh, fish, or powershell.
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell auto-completion script",
	Long: `Generate a shell completion script for vpsm.

Server names and host IPs are dynamically completed from your local inventory
database. Add the output to your shell profile to enable tab-completion.

  # Bash (add to ~/.bashrc):
  source <(vpsm completion bash)

  # Zsh (add to ~/.zshrc):
  source <(vpsm completion zsh)

  # Fish:
  vpsm completion fish | source

  # PowerShell:
  vpsm completion powershell | Out-String | Invoke-Expression`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		}
		return nil
	},
}
