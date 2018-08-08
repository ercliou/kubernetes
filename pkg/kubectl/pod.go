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
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PodGeneratorV1 supports stable generation of a pod
type PodGeneratorV1 struct {
	Name   string
	Images []string
}

// Ensure it supports the generator pattern that uses parameters specified during construction
var _ StructuredGenerator = &PodGeneratorV1{}

// StructuredGenerate outputs a pod object using the configured fields
func (g *PodGeneratorV1) StructuredGenerate() (runtime.Object, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}

	labels := map[string]string{}
	labels["app"] = g.Name
	pod := &v1.Pod{}
	pod.ObjectMeta = metav1.ObjectMeta{
		Name:   g.Name,
		Labels: labels,
	}
	//pod.Name = g.Name
	pod.Spec = buildPodSpec2(g.Images)
	return pod, nil
}

// validate validates required fields are set to support structured generation
func (g *PodGeneratorV1) validate() error {
	if len(g.Name) == 0 {
		return fmt.Errorf("name must be specified")
	}
	if len(g.Images) == 0 {
		return fmt.Errorf("image must be specified")
	}

	return nil
}

func buildPodSpec2(images []string) v1.PodSpec {
	podSpec := v1.PodSpec{Containers: []v1.Container{}}
	for _, imageString := range images {
		// Retain just the image name
		imageSplit := strings.Split(imageString, "/")
		name := imageSplit[len(imageSplit)-1]
		// Remove any tag or hash
		if strings.Contains(name, ":") {
			name = strings.Split(name, ":")[0]
		}
		if strings.Contains(name, "@") {
			name = strings.Split(name, "@")[0]
		}
		podSpec.Containers = append(podSpec.Containers, v1.Container{Name: name, Image: imageString})
	}
	return podSpec
}
