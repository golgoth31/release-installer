package utils

import (
	"context"
	"os"

	"github.com/golgoth31/release-installer/pkg/progressbar"
	"github.com/hashicorp/go-getter"
)

var defaultProgressBar getter.ProgressTracker = progressbar.New()

func Download(src, dst string, progress bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build the client
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	opts := []getter.ClientOption{}
	if progress {
		opts = append(opts, getter.WithProgress(defaultProgressBar))
	}
	client := &getter.Client{
		Ctx:     ctx,
		Src:     src,
		Dst:     dst,
		Pwd:     pwd,
		Mode:    getter.ClientModeAny,
		Options: opts,
	}

	return client.Get()
}
