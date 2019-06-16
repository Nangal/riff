/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands_test

import (
	"testing"

	"github.com/projectriff/riff/pkg/riff/commands"
	rifftesting "github.com/projectriff/riff/pkg/testing"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDoctorOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name:           "valid list",
			Options:        &commands.DoctorOptions{},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestDoctorCommand(t *testing.T) {
	table := rifftesting.CommandTable{
		{
			Name:         "installation is ok",
			Args:         []string{},
			GivenObjects: requiredNamespacesMocks(commands.RequiredNamespaces),
			ExpectOutput: `
Installation is OK
`,
		},
		{
			Name:         "istio-system is missing",
			Args:         []string{},
			GivenObjects: requiredNamespacesMocks(remove(commands.RequiredNamespaces, "istio-system")),
			ExpectOutput: `
Something is wrong!
missing istio-system
`,
		},
		{
			Name:         "multiple namespaces are missing",
			Args:         []string{},
			GivenObjects: requiredNamespacesMocks(remove(remove(commands.RequiredNamespaces, "istio-system"), "riff-system")),
			ExpectOutput: `
Something is wrong!
missing istio-system
missing riff-system
`,
		},
		{
			Name: "error",
			Args: []string{},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "namespaces"),
			},
			ShouldError:  true,
			ExpectOutput: `inducing failure for list namespaces`,
		},
	}

	table.Run(t, commands.NewDoctorCommand)
}

func requiredNamespacesMocks(namespaces []string) []runtime.Object {
	mockObjects := []runtime.Object{}
	for _, namespace := range namespaces {
		mockNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		mockObjects = append(mockObjects, mockNamespace)
	}

	return mockObjects
}

func remove(s []string, r string) []string {
	newList := []string{}
	for i, v := range s {
		if v != r {
			newList = append(newList, s[i])
		}
	}
	return newList
}