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

	"github.com/spf13/cobra"
)

// trunkVlanCmd represents the trunkVlan command
var trunkVlanCmd = &cobra.Command{
	Use:   "trunkVlan",
	Short: "Manage Nexus trunk vlan allowed and native vlan modifications",
	Long: `root command to show, replace, add, remove trunk allowed 
vlans and native vlan beneath an interface.`,
}

func init() {
	RootCmd.AddCommand(trunkVlanCmd)

        // show trunk vlan data
        trunkVlanCmd.AddCommand(trunkVlanShowCmd)
        trunkVlanShowCmd.Flags().StringVar(
            &nxsFlags.intf, "interface", "",
            "interface to get trunk data . ex: ethernet:x/y, port-channel:z, ethernet, port-channel")
        trunkVlanShowCmd.MarkFlagRequired("interface")

        // add trunk vlan data
        trunkVlanCmd.AddCommand(trunkVlanAddCmd)
        trunkVlanAddCmd.Flags().StringVar(
            &nxsFlags.intf, "interface", "",
            "interface to apply trunk. ex: ethernet:x/y, port-channel:z")
        trunkVlanAddCmd.Flags().StringVar(
            &nxsFlags.trunkVlan, "trunkVlan", "",
            "Add trunk vlan 'None', '200', '200-203,205'")
        trunkVlanAddCmd.Flags().StringVar(
            &nxsFlags.nativeVlan, "nativeVlan", "",
            "Add native vlan 'None', '200', '200-203,205'")
        trunkVlanAddCmd.MarkFlagRequired("interface")
        trunkVlanAddCmd.MarkFlagRequired("trunkVlan")

        // replace trunk vlan data
        trunkVlanCmd.AddCommand(trunkVlanReplaceCmd)
        trunkVlanReplaceCmd.Flags().StringVar(
            &nxsFlags.intf, "interface", "",
            "interface to apply trunk. ex: ethernet:x/y, port-channel:z")
        trunkVlanReplaceCmd.Flags().StringVar(
            &nxsFlags.trunkVlan, "trunkVlan", "",
            "Replace trunk vlan 'None', '200', '200-203,205'")
        trunkVlanReplaceCmd.Flags().StringVar(
            &nxsFlags.nativeVlan, "nativeVlan", "",
            "Replace native vlan 'None', '200', '200-203,205'")
        trunkVlanReplaceCmd.MarkFlagRequired("interface")
        trunkVlanReplaceCmd.MarkFlagRequired("trunkVlan")

        // remove trunk vlan data
        trunkVlanCmd.AddCommand(trunkVlanRemoveCmd)
        trunkVlanRemoveCmd.Flags().StringVar(
            &nxsFlags.intf, "interface", "",
            "interface to remove trunk. ex: ethernet:x/y, port-channel:z")
        trunkVlanRemoveCmd.Flags().StringVar(
            &nxsFlags.trunkVlan, "trunkVlan", "",
            "Remove trunk vlan 'None', '200', '200-203,205'")
        trunkVlanRemoveCmd.Flags().StringVar(
            &nxsFlags.nativeVlan, "nativeVlan", "",
            "Remove native vlan 'None', '200', '200-203,205'")
        trunkVlanRemoveCmd.MarkFlagRequired("interface")
        trunkVlanRemoveCmd.MarkFlagRequired("trunkVlan")

}

// trunkVlanShowCmd represents the trunk_vlan show Interface command
var trunkVlanShowCmd = &cobra.Command{
        Use:   "show",
        Short: "To show trunk data for a given interface",
        Long: `To show trunk data for a given interface`,
        Run: runTrunkVlanShowCmd,
}

func runTrunkVlanShowCmd(cmd *cobra.Command, args []string) {
        var resp []map[string]interface{}
        var err error
        if len(nxsFlags.intf) == 0 {
                cmd.Help()
                return
        }

        tw, client := startTransaction(cmd)
        defer tw.Flush()
        defer client.Logout()

        // Print the legend
        //fmt.Fprintf(tw, "ID\t\tNative\t\tTrunk\tMode\tState\tDescr\n")
        fmt.Fprintf(tw, "ID\t\tNative\tTrunk\tMode\tState\tDescr\n")

        resp, err = client.GetInterface(nxsFlags.intf)

        if err != nil {
                exitWithError(ExitError, err)
        }

        for _, r := range resp {
                //fmt.Fprintf(tw, "%s\t\t%s\t%s\t%s\t%s\t%s\n", r["id"], r["nativeVlan"], 
                //           r["trunkVlans"], r["mode"], r["adminSt"], r["descr"])
                tvlan := r["trunkVlans"]
                if r["trunkVlans"] == "" {
                    tvlan = "None"
                }
                fmt.Fprintf(tw, "%s\t\t%s\t%s\t%s\t%s\t%s\n", r["id"], r["nativeVlan"],
                           tvlan, r["mode"], r["adminSt"], r["descr"])
        }

}

// sendTrunkVlanCmd - Sends interface trunk vlan update commands to nexus switch
func sendTrunkVlanCmd(cmd *cobra.Command, args []string, 
                      trunkVlanPfx string) {

        var native string
        if len(nxsFlags.intf) == 0 {
                cmd.Help()
                return
        }
        if len(nxsFlags.trunkVlan) == 0 {
                cmd.Help()
                return
        }

        tw, client := startTransaction(cmd)
        defer tw.Flush()
        defer client.Logout()

        nxsFlags.trunkVlan = trunkVlanPfx + nxsFlags.trunkVlan
        native = nxsFlags.nativeVlan
        switch trunkVlanPfx {
        case "+":
            fallthrough
        case "":  // This is valid for replace operation
            native = nxsFlags.nativeVlan
        case "-":
            if native != "" {
                native = "None"
            }
        default:
            fmt.Fprintf(tw, `ERROR: Unexpected Trunk prefix received'
                            expected: "+", or "-", or ""`,
                            trunkVlanPfx)
            return
        }


        // Note client package uses add naming instead of create.
        err := client.AddTrunkVlan(nxsFlags.intf,
                                   nxsFlags.trunkVlan, native)
        if err != nil {
                exitWithError(ExitError, err)
        }

        if nxsFlags.nativeVlan != "" {
            fmt.Printf("Trunk Vlan %s native %s created for interface %s.\n", 
                       nxsFlags.trunkVlan, nxsFlags.nativeVlan, nxsFlags.intf)
        } else {
            fmt.Printf("Trunk Vlan %s created for interface %s.\n", 
                       nxsFlags.trunkVlan, nxsFlags.intf)
        }

}

// trunkVlanAddCmd represents the trunk_vlan add command
var trunkVlanAddCmd = &cobra.Command{
        Use:   "add",
        Short: "To add trunk and/or native vlan to interface",
        Long: `To add trunk and/or native vlan for a given
interface. For example:
'200' configures 'switchport trunk allowed vlan add 200',
'200-203,205' configures 'switchport trunk allowed vlan add 200-203,205'`,
        Run:  runTrunkVlanAdd,
}

func runTrunkVlanAdd(cmd *cobra.Command, args []string) {

    sendTrunkVlanCmd(cmd, args, "+")
}

// trunkVlanReplaceCmd represents the trunk_vlan replace command
var trunkVlanReplaceCmd = &cobra.Command{
        Use:   "replace",
        Short: "To replace all trunk/native vlans for a given interface",
        Long: `To replace all trunk vlans and/or native vlan for a 
given interface. For example:
'None' configures 'switchport trunk allowed vlan None',
'200' configures 'switchport trunk allowed vlan 200',
'200-203,205' configures 'switchport trunk allowed vlan 200-203,205'`,
        Run:  runTrunkVlanReplace,
}

func runTrunkVlanReplace(cmd *cobra.Command, args []string) {

    sendTrunkVlanCmd(cmd, args, "")
}

// trunkVlanRemoveCmd represents the trunk_vlan remove command
var trunkVlanRemoveCmd = &cobra.Command{
        Use:   "remove",
        Short: "To remove trunk and/or native vlan for a given interface",
        Long: `To remove trunk and/or native vlans for a given interface.
For example:
'200' configures 'switchport trunk allowed vlan remove 200',
'200-203,205' configures 'switchport trunk allowed vlan remove 200-203,205'`,
        Run:  runTrunkVlanRemove,
}

func runTrunkVlanRemove(cmd *cobra.Command, args []string) {

    sendTrunkVlanCmd(cmd, args, "-")
}
