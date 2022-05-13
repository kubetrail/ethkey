/*
Copyright © 2022 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/ethkey/pkg/flags"
	"github.com/kubetrail/ethkey/pkg/run"
	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode key",
	Long:  `Decode key to derive address`,
	RunE:  run.Decode,
	Args:  cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	f := decodeCmd.Flags()

	f.String(flags.Key, "", "Private or public key")
}
