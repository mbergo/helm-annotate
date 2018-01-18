package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
)

type getCmd struct {
	release      string
	client       helm.Interface
	timeout      int64
	wait         bool
	manifestType string
	namespace    string
}

// newGetCmd allows getting annotations of kubernetes manifests
func newGetCmd() *cobra.Command {

	gc := &getCmd{}

	cmd := &cobra.Command{
		Use:     "get RELEASE",
		Short:   fmt.Sprintf("gets annotation on a release"),
		PreRunE: setupConnection,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("This command neeeds 1 argument: release name")
			}
			if gc.manifestType == "" {
				return fmt.Errorf("You have to select a manifest type")
			}
			gc.release = args[0]
			gc.client = ensureHelmClient(gc.client)

			return gc.run()
		},
	}
	return cmd
}

func (e *getCmd) run() error {

	res, err := e.client.ReleaseContent(e.release)
	if err != nil {
		return err
	}
	values, err := chartutil.ReadValues([]byte(res.Release.Config.Raw))
	if err != nil {
		return errors.Wrap(err, "PB unable to de-serialize the release config to yaml")
	}
	y := values.AsMap()
	for k, v := range y {
		// We only care about our own values
		if strings.HasPrefix(k, "ANNO_") {
			key := strings.Replace(k, "ANNO_", "", -1)
			fmt.Printf("export %v=\"%v\"\n", key, v)
		}
	}
	fmt.Printf("# Run this command to configure your shell:\n")
	fmt.Printf("# eval $(helm annotate get %v)\n", e.release)

	return nil
}
