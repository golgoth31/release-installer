package migration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	installv1 "github.com/golgoth31/release-installer/internal/migration/install/v1"
	"github.com/golgoth31/release-installer/pkg/config"
	"github.com/golgoth31/release-installer/pkg/output"
	"github.com/golgoth31/release-installer/pkg/release"
	"sigs.k8s.io/yaml"
)

var out output.Output

func Migrate(homedir string, version string, conf *config.Config) error {
	sem := semver.MustParse(version)

	switch sem.Major() {
	case 1:
		out.StepTitle(fmt.Sprintf("Migration from %s", version))

		if err := os.Rename(homedir+"/releases", homedir+"/references"); err != nil {
			return err
		}

		if err := os.Rename(homedir+"/install", homedir+"/releases"); err != nil {
			return err
		}

		var files []string

		if err := filepath.Walk(homedir+"/releases", func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}

			return nil
		}); err != nil {
			return err
		}

		for _, rel := range files {
			var relFiles []string

			if err := filepath.Walk(rel, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && info.Name() != "default" {
					relFiles = append(relFiles, path)
				}

				return nil
			}); err != nil {
				return err
			}

			for _, v := range relFiles {
				inst := installv1.Install{}
				data, _ := os.ReadFile(v)
				fmt.Printf("%s\n", v)
				yaml.Unmarshal(data, &inst)
				rel := release.New(conf, inst.Metadata.Release, inst.Spec.Version)
				rel.Rel.Spec.Arch = inst.Spec.Arch
				rel.Rel.Spec.Os = inst.Spec.Os
				rel.Rel.Spec.Path = inst.Spec.Path
				rel.Rel.Spec.Default = inst.Spec.Default
				rel.Write()
			}
		}
	}

	return nil
}
