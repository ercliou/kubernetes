/*
Copyright 2016 The Kubernetes Authors.

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

package create

import (
	"github.com/spf13/cobra"

	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

var (
	podLong = templates.LongDesc(i18n.T(`
	Create a pod with the specified name.`))

	podExample = templates.Examples(i18n.T(`
	# Create a new pod named my-pod that runs the busybox image.
	kubectl create pod my-pod --image=busybox`))
)

type PodOpts struct {
	CreateSubcommandOptions *CreateSubcommandOptions
}

// NewCmdCreatePod is a macro command to create a new Pod.
func NewCmdCreatePod(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	options := &PodOpts{
		CreateSubcommandOptions: NewCreateSubcommandOptions(ioStreams),
	}

	cmd := &cobra.Command{
		Use: "pod NAME --image=image [--dry-run]", // TODO add --cmd option?
		DisableFlagsInUseLine: true,
		Aliases:               []string{"po"},
		Short:                 i18n.T("Create a Pod with the specified name."),
		Long:                  podLong,
		Example:               podExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Complete(f, cmd, args))
			cmdutil.CheckErr(options.Run())
		},
	}

	options.CreateSubcommandOptions.PrintFlags.AddFlags(cmd)

	cmdutil.AddApplyAnnotationFlags(cmd)
	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddGeneratorFlags(cmd, cmdutil.PodV1GeneratorName)
	cmd.Flags().StringSlice("image", []string{}, "Image name to run.")
	cmd.MarkFlagRequired("image")
	return cmd
}

func (o *PodOpts) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	name, err := NameFromCommandArgs(cmd, args)
	if err != nil {
		return err
	}

	// TODO fallback shit

	imageNames := cmdutil.GetFlagStringSlice(cmd, "image")
	var generator kubectl.StructuredGenerator
	switch generatorName := cmdutil.GetFlagString(cmd, "generator"); generatorName {
	case cmdutil.PodV1GeneratorName:
		generator = &kubectl.PodGeneratorV1{Name: name, Images: imageNames}
	default:
		return errUnsupportedGenerator(cmd, generatorName)
	}

	return o.CreateSubcommandOptions.Complete(f, cmd, args, generator)
}

func (o *PodOpts) Run() error {
	return o.CreateSubcommandOptions.Run()
}
