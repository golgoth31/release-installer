// Package reference ...
package reference

import (
	"context"
	"fmt"

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
}

// ListVersions string list.
func (r *Reference) ListVersions(num int) []string {
	ctx := context.Background()
	client := github.NewClient(nil)
	opts := &github.ListOptions{
		Page:    1,
		PerPage: num,
	}
	out := []string{}

	reference, _, err := client.Repositories.ListReleases(ctx, r.Ref.Spec.Repo.GetOwner(), r.Ref.Spec.Repo.GetName(), opts)
	if err != nil {
		logger.StdLog.Fatal().Err(err).Msg("Unable to get version list")
	}

	for _, val := range reference {
		out = append(out, val.GetTagName())
	}

	return out
}
