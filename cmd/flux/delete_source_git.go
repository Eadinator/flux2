/*
Copyright 2020 The Flux authors

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

package main

import (
	"context"
	"fmt"

	"github.com/fluxcd/flux2/internal/utils"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

var deleteSourceGitCmd = &cobra.Command{
	Use:   "git [name]",
	Short: "Delete a GitRepository source",
	Long:  "The delete source git command deletes the given GitRepository from the cluster.",
	Example: `  # Delete a Git repository
  flux delete source git podinfo
`,
	RunE: deleteSourceGitCmdRun,
}

func init() {
	deleteSourceCmd.AddCommand(deleteSourceGitCmd)
}

func deleteSourceGitCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("git name is required")
	}
	name := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), rootArgs.timeout)
	defer cancel()

	kubeClient, err := utils.KubeClient(rootArgs.kubeconfig, rootArgs.kubecontext)
	if err != nil {
		return err
	}

	namespacedName := types.NamespacedName{
		Namespace: rootArgs.namespace,
		Name:      name,
	}

	var git sourcev1.GitRepository
	err = kubeClient.Get(ctx, namespacedName, &git)
	if err != nil {
		return err
	}

	if !deleteArgs.silent {
		prompt := promptui.Prompt{
			Label:     "Are you sure you want to delete this source",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err != nil {
			return fmt.Errorf("aborting")
		}
	}

	logger.Actionf("deleting source %s in %s namespace", name, rootArgs.namespace)
	err = kubeClient.Delete(ctx, &git)
	if err != nil {
		return err
	}
	logger.Successf("source deleted")

	return nil
}
