/*
Copyright 2015 The Kubernetes Authors.

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

package kubectl

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodGenerate(t *testing.T) {
	tests := []struct {
		name      string
		podName   string
		images    []string
		expected  *v1.Pod
		expectErr bool
	}{
		{
			name:    "pod name and images ok",
			podName: "my-pod-ok",
			images:  []string{"nn/image1", "registry/nn/image2", "nn/image3:tag", "nn/image4@digest", "nn/image5@sha256:digest"},

			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "my-pod-ok",
					Labels: map[string]string{"app": "my-pod-ok"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{Name: "image1", Image: "nn/image1"},
						{Name: "image2", Image: "registry/nn/image2"},
						{Name: "image3", Image: "nn/image3:tag"},
						{Name: "image4", Image: "nn/image4@digest"},
						{Name: "image5", Image: "nn/image5@sha256:digest"},
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "no name",
			images:    []string{"nn/image1"},
			expectErr: true,
		},
		{
			name:      "empty name",
			podName:   "",
			images:    []string{"nn/image1"},
			expectErr: true,
		},
		{
			name:      "no images",
			podName:   "pod-no-images",
			expectErr: true,
		},
		{
			name:      "no images and no name",
			expectErr: true,
		},
		{
			name:      "empty images",
			podName:   "pod-empty-images",
			images:    []string{},
			expectErr: true,
		},
	}
	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			generator := PodGeneratorV1{
				Name:   tt.podName,
				Images: tt.images,
			}
			obj, err := generator.StructuredGenerate()
			switch {
			case tt.expectErr && err != nil:
				return // loop, since there's no output to check
			case tt.expectErr && err == nil:
				t.Errorf("%v: expected error and didn't get one", index)
				return // loop, no expected output object
			case !tt.expectErr && err != nil:
				t.Errorf("%v: unexpected error %v", index, err)
				return // loop, no output object
			}
			if !reflect.DeepEqual(obj.(*v1.Pod), tt.expected) {
				t.Errorf("\nexpected:\n%#v\nsaw:\n%#v", tt.expected, obj.(*v1.Pod))
			}
		})
	}
}
