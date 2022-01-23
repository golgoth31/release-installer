// Package reference ...
package reference

import (
	"context"
	"fmt"
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

	reference, _, err := client.Repositories.ListReleases(ctx, r.Spec.Repo.Owner, r.Spec.Repo.Name, opts)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to get version list")
	}

	for _, val := range reference {
		out = append(out, val.GetTagName())
	}

	return out
}
