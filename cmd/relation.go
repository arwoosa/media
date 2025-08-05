/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/arwoosa/vulpes/relation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// relationCmd represents the relation command
var relationCmd = &cobra.Command{
	Use:   "relation",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		relation.Initialize(
			relation.WithWriteAddr(viper.GetString("relation.write_uri")),
			relation.WithReadAddr(viper.GetString("relation.read_uri")))

		relation.AddUserResourceRole(context.Background(), "d6e0a454-4484-41a3-a0d0-652dc5c4aad4", "Image", "b04cc187-11c8-4eb7-19b7-5dcd0064c400", relation.RoleOwner)

		fmt.Println(relation.Check(context.Background(), "Image", "b04cc187-11c8-4eb7-19b7-5dcd0064c400", "viewer", "User", "d6e0a454-4484-41a3-a0d0-652dc5c4aad4"))
		fmt.Println(relation.CheckBySubjectId(context.Background(), "Image", "abc", "viewer", "user:123"))

		fmt.Println(relation.DeleteObjectId(context.Background(), "Image", "b04cc187-11c8-4eb7-19b7-5dcd0064c400"))

		fmt.Println("relation")
	},
}

func init() {
	rootCmd.AddCommand(relationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// relationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// relationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
