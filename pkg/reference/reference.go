<<<<<<< HEAD
// Package reference ...
=======
// Package reference is used to manipulate references.
>>>>>>> ce9013b (Feat: V2)
package reference

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"os"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// New ...
func New(rel string) *Reference {
	return loadYaml(rel)
}

func loadYaml(file string) *Reference {
	referencePath := fmt.Sprintf(
		"%s/%s",
		viper.GetString("homedir"),
		viper.GetString("releases.dir"),
	)

	if file == "myself" {
		return myself
	}

	r := &Reference{} // nolint: exhaustivestruct

	data, err := os.ReadFile(fmt.Sprintf("%s/%s.yaml", referencePath, file))
	if err != nil {
		logger.StdLog.Debug().Err(err).Msg("Unable to read release definition")
		logger.StdLog.Fatal().Msg("Unknown release")
	}

	if err := yaml.Unmarshal(data, r); err != nil {
		logger.StdLog.Fatal().Err(err).Msg("")
	}

	return r
=======

	logger "github.com/golgoth31/release-installer/pkg/log"
	"github.com/golgoth31/release-installer/pkg/utils"
	"github.com/google/go-github/v32/github"
	"google.golang.org/protobuf/encoding/protojson"
)

// Load reference data from yaml manifest.
func (r *Reference) Load() error {
	jsonData, err := utils.Load(r.File)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	logger.StdLog.Debug().Msgf("Reference json data: %s", jsonData)

	if errUnmarshall := protojson.Unmarshal(jsonData, &r.Ref); errUnmarshall != nil {
		return fmt.Errorf("%w", errUnmarshall)
	}

	return nil
>>>>>>> ce9013b (Feat: V2)
}

// ListVersions ...
func (r *Reference) ListVersions(num int) []string {
	ctx := context.Background()
	client := github.NewClient(nil)
	opts := &github.ListOptions{
		Page:    1,
		PerPage: num,
	}
	out := []string{}

<<<<<<< HEAD
	reference, _, err := client.Repositories.ListReleases(ctx, r.Spec.Repo.Owner, r.Spec.Repo.Name, opts)
=======
	reference, _, err := client.Repositories.ListReleases(ctx, r.Ref.Spec.Repo.GetOwner(), r.Ref.Spec.Repo.GetName(), opts)
>>>>>>> ce9013b (Feat: V2)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to get version list")
	}

	for _, val := range reference {
		out = append(out, val.GetTagName())
	}

	return out
}
