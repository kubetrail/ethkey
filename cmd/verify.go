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

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify signature",
	Long: `This command verifies input hash signed using
private key.`,
	RunE: run.Verify,
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	f := verifyCmd.Flags()

	f.String(flags.Hash, "", "Hash of input data")
	f.String(flags.Sign, "", "Signature of hash")
	f.String(flags.PubKey, "", "Public key")
}
