/*
Copyright 2018 The Fission Authors.

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

package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fission/fission/pkg/fission-cli/console"
	"github.com/fission/fission/pkg/utils"
)

type Resource interface {
	Dump(string)
}

func getFileName(dumpdir string, meta metav1.ObjectMeta) string {
	f := fmt.Sprintf("%v/%v_%v_%v.txt", dumpdir, meta.Namespace, meta.Name, meta.ResourceVersion)
	return filepath.Clean(f)
}

func getPodFileName(dumpdir string, pod metav1.ObjectMeta, containerName string) string {
	f := fmt.Sprintf("%v/%v_%v_%v_%v.txt", dumpdir, pod.Namespace, pod.Name, pod.ResourceVersion, containerName)
	return filepath.Clean(f)
}

func writeToFile(file string, obj interface{}) {
	bs, err := yaml.Marshal(obj)
	if err != nil {
		console.Error(fmt.Sprintf("Error encoding object: %v", err))
		return
	}

	// Due to unknown reason, the kubernetes objectMeta fields contain
	// empty byte and will fail os.Create/os.Openfile with error message
	// "open <file> invalid argument". To fix the problem, we need to
	// remove the empty byte from string.
	file = string(utils.RemoveZeroBytes([]byte(file)))

	err = os.WriteFile(file, bs, 0644)
	if err != nil {
		console.Error(fmt.Sprintf("Error writing file %v: %v", file, err))
	}
}
