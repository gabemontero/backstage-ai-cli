package util

import (
	"fmt"
	"k8s.io/cli-runtime/pkg/printers"
	"os"
	"sigs.k8s.io/yaml"
)

func PrintYaml(obj interface{}, addDivider bool) error {
	writer := printers.GetNewTabWriter(os.Stdout)
	output, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = writer.Write(output)
	if addDivider {
		fmt.Fprintln(os.Stdout, "---")
	}
	return err
}
