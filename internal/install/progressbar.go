package install

import (
	"io"

	logger "github.com/golgoth31/release-installer/internal/log"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

func (cpb *progressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	pb := mpb.New(mpb.WithWidth(60))
	// Parameters of th new progress bar
	bar := pb.AddBar(totalSize,
		mpb.PrependDecorators(
			decor.OnComplete(
				decor.Name(logger.LineStart()), logger.OkStatus(),
			),
		),
		mpb.AppendDecorators(
			decor.Name(src),
			decor.Name(" "),
			decor.AverageSpeed(decor.UnitKB, "(% .2f)"),
		),
	)
	// bar := pb.AddBar(totalSize,
	// 	mpb.PrependDecorators(
	// 		decor.Name(src),
	// 		decor.Name(" "),
	// 		decor.CountersKiloByte("% .2f / % .2f "),
	// 		decor.AverageSpeed(decor.UnitKB, "(% .2f)"),
	// 	),
	// 	mpb.AppendDecorators(
	// 		decor.Percentage(),
	// 		decor.Name(" - "),
	// 		decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 4}),
	// 		decor.Name(" - "),
	// 		decor.OnComplete(
	// 			decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
	// 		),
	// 	),
	// )

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

func (c *readCloser) Close() error { return c.close() }
