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

// vlanCmd represents the vlan command
var vlanCmd = &cobra.Command{
	Use:   "vlan",
	Short: "Manage Nexus VLAN modifications",
	Long: `Manage addition, removal, and display of VLAN
on Nexus Switch`,
}

func init() {
	RootCmd.AddCommand(vlanCmd)

        vlanCmd.AddCommand(vlanAddCmd)
        vlanAddCmd.Flags().StringVar(
            &nxsFlags.vlanid, "vlanid", "",
            "Vlan to configure")
        vlanAddCmd.Flags().StringVar(
            &nxsFlags.vnSeg, "vnSegment", "",
            "Vn segment id of the VLAN")
        vlanAddCmd.MarkFlagRequired("vlan")

        vlanCmd.AddCommand(vlanRemoveCmd)
        vlanRemoveCmd.Flags().StringVar(
            &nxsFlags.vlanid, "vlanid", "",
            "Vlan to remove ")
        vlanRemoveCmd.MarkFlagRequired("vlan")

        vlanCmd.AddCommand(vlanShowCmd)
        vlanShowCmd.Flags().StringVar(
            &nxsFlags.vlanid, "vlanid", "",
            "Vlan to show ")
}


// vlanShowCmd represents the show vlan operations
var vlanShowCmd = &cobra.Command{
        Use:   "show",
        Short: "To show vlan data for a given vlan",
        Long: `To show vlan data for a given vlan`,
        Run: runVlanShowCmd,
}

func runVlanShowCmd(cmd *cobra.Command, args []string) {
        var resp []map[string]interface{}
        var err error

        tw, client := startTransaction(cmd)
        defer tw.Flush()
        defer client.Logout()

        // Print the legend
        fmt.Fprintf(tw, "\tID\t\tVNI\t\tAdmin State\tOperState\tName\n")

        resp, err = client.GetVlan(nxsFlags.vlanid)

        if err != nil {
                exitWithError(ExitError, err)
        }

        for _, r := range resp {
                if r["BdOperName"] == "" {
                    continue
                }
                if r["accEncap"] == "unknown" {
                    fmt.Fprintf(tw, "\t%s\t\t%s\t\t%s\t%s\t%s\n", r["id"], r["accEncap"],
                               r["adminSt"], r["operSt"], r["BdOperName"])
                } else {
                    fmt.Fprintf(tw, "\t%s\t\t%s\t\t%s\t%s\t%s\n", r["id"], r["accEncap"],
                               r["adminSt"], r["operSt"], r["BdOperName"])
                }

        }

}

func sendVlanCmd(cmd *cobra.Command, args []string, action string) {

        var err error

        if nxsFlags.vlanid == "" {
            fmt.Printf("ERROR: VLANID not present\n")
            cmd.Help()
            return
        }

        tw, client := startTransaction(cmd)
        defer tw.Flush()
        defer client.Logout()

        switch action {
        case "add":
            err = client.AddVlan(nxsFlags.vlanid, nxsFlags.vnSeg)
        case "remove":
            err = client.DeleteVlan(nxsFlags.vlanid)
        default:
            err = fmt.Errorf(`ERROR: Unexpected action received "%s"
                            expected: "add", or "remove"`,
                            action)
        }

        // Note client package uses add naming instead of create.
        if err != nil {
                exitWithError(ExitError, err)
        }

        if nxsFlags.vnSeg != "" {
            fmt.Fprintf(tw, "%s Vlan %s with vn-segment %s successful.\n",
                       action, nxsFlags.vlanid, nxsFlags.vnSeg)
        } else {
            fmt.Fprintf(tw, "%s Vlan %s successful.\n",
                       action, nxsFlags.vlanid)
        }
}

// vlanAddCmd represents the vlan add command
var vlanAddCmd = &cobra.Command{
        Use:   "add",
        Short: "To add vlan to Nexus Switch",
        Long: `To add vlan to Nexus Switch`,
        Run:  runVlanAdd,
}

func runVlanAdd(cmd *cobra.Command, args []string) {

    sendVlanCmd(cmd, args, "add")
}

// vlanRemoveCmd represents the vlan remove command
var vlanRemoveCmd = &cobra.Command{
        Use:   "remove",
        Short: "To remove vlan from Nexus Switch",
        Long: `To remove vlan from Nexus Switch`,
        Run:  runVlanRemove,

}

func runVlanRemove(cmd *cobra.Command, args []string) {

    sendVlanCmd(cmd, args, "remove")
}
