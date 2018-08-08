/*
Copyright 2014 The Kubernetes Authors.

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
	"fmt"
	"net/http"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	cmdtesting "k8s.io/kubernetes/pkg/kubectl/cmd/testing"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
)

func TestCreatePod(t *testing.T) {
	podObject := &v1.Pod{}
	podObject.Name = "my-pod"
	tf := cmdtesting.NewTestFactory().WithNamespace("test")
	defer tf.Cleanup()

	codec := legacyscheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)
	ns := legacyscheme.Codecs

	tf.Client = &fake.RESTClient{
		GroupVersion:         schema.GroupVersion{Version: "v1"},
		NegotiatedSerializer: ns,
		Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			panic("wtfpod")

			fmt.Println("ericcc")
			switch p, m := req.URL.Path, req.Method; {
			case p == "/podsx" && m == "POST":
				return &http.Response{StatusCode: http.StatusCreated, Header: defaultHeader(), Body: objBody(codec, podObject)}, nil
			default:
				t.Fatalf("unexpected request: %#v\n%#v", req.URL, req)
				return nil, nil
			}
		}),
	}
	tf.ClientConfigVal = &restclient.Config{}

	ioStreams, _, buf, _ := genericclioptions.NewTestIOStreams()
	cmd := NewCmdCreatePod(tf, ioStreams)
	cmd.Flags().Set("dry-run", "false")
	cmd.Flags().Set("output", "name")
	cmd.Flags().Set("image", "pancakes/sweet.potato:v2")
	cmd.Run(cmd, []string{podObject.Name})
	expectedOutput := "pod/" + podObject.Name + "\n"
	if buf.String() != expectedOutput {
		t.Errorf("expected output: %s, but got: %s", expectedOutput, buf.String())
	}
}

func TestCreatePodNoImage(t *testing.T) {
	tf := cmdtesting.NewTestFactory()
	defer tf.Cleanup()

	ioStreams := genericclioptions.NewTestIOStreamsDiscard()
	cmd := NewCmdCreatePod(tf, ioStreams)
	cmd.Flags().Set("output", "name")

	options := &PodOpts{
		CreateSubcommandOptions: &CreateSubcommandOptions{
			PrintFlags: genericclioptions.NewPrintFlags("created").WithTypeSetter(legacyscheme.Scheme),
			DryRun:     true,
			IOStreams:  ioStreams,
		},
	}

	err := options.Complete(tf, cmd, []string{"my-pod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = options.Run()
	fmt.Println("erico")
	fmt.Println(err)
	if err == nil {
		t.Errorf("expected error and didn't get one")
	}
}

func TestCreatePodUnsupportedGenerator(t *testing.T) {
}
