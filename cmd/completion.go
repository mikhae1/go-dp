package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// TODO: at cobra version >= 0.1.3 __custom_func should be renamed to __gp_custom_func
const (
	defaultBashCompletionFunc = `
# __gp_parse_envs is a custom completion func to support dynamic env/targets names
__gp_parse_envs() {
	local out word

	word="$cur"
	[ "$cword" -gt "2" ] && word="$prev"

	__gp_debug "__custom_func: for \"$word\": $(gp envs $word)"

	if out=$(gp envs $word 2>/dev/null); then
		COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
	fi
}

# runs last
__custom_func() {
	[ "$cword" -gt "1" ] && __gp_parse_envs
}
`
)

var completionCmd = &cobra.Command{
	Use:   "completion <bash|zsh>",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

	. <(gp completion)

	To configure your bash shell to load completions for each session add to your bashrc

	# ~/.bashrc or ~/.profile
	. <(gp completion bash)

	# ~/.zshrc
	. <(gp zcompletion zsh)
	`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh"},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "zsh" {
			rootCmd.GenZshCompletion(os.Stdout)
			return
		}

		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// enable dynamic env/targets completion
	rootCmd.BashCompletionFunction = defaultBashCompletionFunc
}
