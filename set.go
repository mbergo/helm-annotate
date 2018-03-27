package main

import (
	"fmt"
	"io"
	"strings"

	"k8s.io/helm/pkg/chartutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
)

type setCmd struct {
	release     string
	out         io.Writer
	client      helm.Interface
	annotations []string
	namespace   string
}

// newSetCmd allows adding annotation to kubernetes manifests
func newSetCmd() *cobra.Command {

	edit := &setCmd{}

	cmd := &cobra.Command{
		Use:     "set [flags] RELEASE",
		Short:   fmt.Sprintf("sets annotation on a release"),
		PreRunE: setupConnection,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("There are %v arguments\n", len(args))
			for i, v := range args {
				fmt.Printf("The %vth arg is %v\n", i, v)
			}

			if len(args) != 1 {
				return fmt.Errorf("This command neeeds 1 argument: release name")
			}
			edit.release = args[0]
			edit.client = ensureHelmClient(edit.client)

			return edit.run()
		},
	}

	f := cmd.Flags()
	f.StringSliceVar(&edit.annotations, "annotations", nil, "list of annotations to add to the release")
	return cmd
}

func toMap(annos []string) (map[string]string, error) {
	annotationToApply := make(map[string]string)
	for _, v := range annos {
		if !strings.Contains(v, "=") {
			return nil, fmt.Errorf("All annotations should be in the format key=value was '%v'", v)
		}
		t := strings.Split(v, "=")
		annotationToApply[t[0]] = strings.Join(t[1:], "=")
	}
	return annotationToApply, nil
}

func (e *setCmd) run() error {

	annotationToApply, err := toMap(e.annotations)
	if err != nil {
		return errors.Wrap(err, "rror converting annotation to map")

	}
	res, err := e.client.ReleaseContent(e.release)
	if err != nil {
		return err
	}
	values, err := chartutil.ReadValues([]byte(res.Release.Config.Raw))
	if err != nil {
		return errors.Wrap(err, "unable to read values from release")
	}

	y, err := values.YAML()
	if err != nil {
		return errors.Wrap(err, "unable to convert values to YAML")
	}
	vm := chartutil.FromYaml(y)

	for k, v := range annotationToApply {
		vm["ANNO_"+k] = v
	}

	configYAML := chartutil.ToYaml(vm)
	_, err = e.client.UpdateReleaseFromChart(
		res.Release.Name,
		res.Release.Chart,
		helm.UpdateValueOverrides([]byte(configYAML)))
	if err != nil {
		return err
	}
	return nil
}
