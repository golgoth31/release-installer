// Package progressbar ...
package progressbar

import (
	"io"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

const pbWidth = 200

// TrackProgress generates a progress bar.
func (cpb *ProgressBar) TrackProgress(
	src string,
	currentSize,
	totalSize int64,
	stream io.ReadCloser,
) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	pb := mpb.New(mpb.WithWidth(pbWidth))
	// Parameters of th new progress bar
	bar := pb.AddBar(totalSize,
		mpb.PrependDecorators(
			decor.OnComplete(
				decor.Name(logger.InfoStatus()), logger.OkStatus(),
			),
		),
		mpb.AppendDecorators(
			decor.Name(src),
			decor.Name(" "),
			decor.AverageSpeed(decor.UnitKB, "(% .2f)"),
		),
	)

	reader := bar.ProxyReader(stream)

	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			pb.Wait()

			return nil
		},
	}
}

// Close ...
func (c *readCloser) Close() error { return c.close() }
