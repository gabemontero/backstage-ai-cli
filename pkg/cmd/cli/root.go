package cli

import (
     "github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/kserve"
     "github.com/gabemontero/backstage-ai-cli/pkg/config"
     "github.com/spf13/cobra"
)

// NewCmd create a new root command, linking together all sub-commands organized by groups.
func NewCmd() *cobra.Command {
     cfg := &config.Config{}
     bkstgAI := &cobra.Command{
          Use:  "bkstg-ai",
          Long: "Backstage AI is a command line tool that facilitates management of AI related Entities in the Backstage Catalog.",
          Run: func(cmd *cobra.Command, args []string) {
               cmd.Help()
          },
     }

     bkstgAI.PersistentFlags().StringVar(&(cfg.Kubeconfig), "kubeconfig", cfg.Kubeconfig,
          "Path to the kubeconfig file to use for CLI requests.")
     bkstgAI.AddCommand(kserve.NewCmd(cfg))

     return bkstgAI
}
