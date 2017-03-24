// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"text/tabwriter"
        "github.com/danehans/nxgo/nx"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var (
        RootCmd = &cobra.Command{
     	    Use:   "cli",
	    Short: "A command line utility for managing Cisco Nexus Switches",
	    Long: `A command line utility for managing Cisco Nexus Switches
                  to get help about a resource or command run 
                  "nxsctl help resource"`,
        }

        globalFlags = struct {
                hosts []string
                username string
                password string
                debug bool
        }{}

        nxsFlags = struct{
            intf        string
            iftype      string
            trunkVlan   string
            nativeVlan  string
            vnSeg       string
            vlanid      string
        }{}
)

func init() {
        RootCmd.PersistentFlags().StringSliceVar(&globalFlags.hosts, "hosts", []string{""}, "Nexus API Endpoints")
        // gRPC TLS Server Verification
        RootCmd.PersistentFlags().StringVar(&globalFlags.username, "user", "admin", "Nexus Username")
        // gRPC TLS Client Authentication
        RootCmd.PersistentFlags().StringVar(&globalFlags.password, "pass", "Cisco12345", "Nexus Password")
        RootCmd.PersistentFlags().BoolVar(&globalFlags.debug, "debug", false, "Enable debug") 
        cobra.EnablePrefixMatching = true
}


// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// mustClientFromCmd returns a Nexus client or exits.
func mustClientFromCmd(cmd *cobra.Command) *nx.Client {
        hosts := hostsFromCmd(cmd)
        user := userFromCmd(cmd)
        pass := passFromCmd(cmd)
        debug := debugFromCmd(cmd)

        // Sets ACI client configuration options.
        opts := nx.ClientOptions{
                Hosts: hosts,
                User:  user,
                Pass:  pass,
                Debug: debug,
        }

        // Creates an instance of the ACI Client.
        client, err := nx.New(opts)
        if err != nil {
                exitWithError(ExitBadConnection, err)
        }
        return client
}

// hostsFromCmd returns the hosts argument.
func hostsFromCmd(cmd *cobra.Command) []string {
        hosts, err := cmd.Flags().GetStringSlice("hosts")
        if err != nil {
                exitWithError(ExitBadArgs, err)
        }
        return hosts
}

func userFromCmd(cmd *cobra.Command) string {
        user, err := cmd.Flags().GetString("user")
        if err != nil {
                exitWithError(ExitBadArgs, err)
        }
        return user
}

// passFromCmd returns the password argument.
func passFromCmd(cmd *cobra.Command) string {
        password, err := cmd.Flags().GetString("pass")
        if err != nil {
                exitWithError(ExitBadArgs, err)
        }
        return password
}

// passFromCmd returns the password argument.
func debugFromCmd(cmd *cobra.Command) bool {
        _, err := cmd.Flags().GetBool("debug")
        if err != nil {
                return false
        }
        return true
}

func startTransaction(cmd *cobra.Command) (*tabwriter.Writer, *nx.Client){

        tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
        defer tw.Flush()

        // ACI Client
        client := mustClientFromCmd(cmd)

        //Log into Nexus NXAPI
        err := client.Login()
        if err != nil {
                exitWithError(ExitError, err)
        }

        return tw, client
}
