package utils

import (
	"fmt"
	"io"
	"os"

	logger "github.com/golgoth31/release-installer/pkg/log"
	"sigs.k8s.io/yaml"
)

// Load loads yaml file and convert it to json
func Load(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, fmt.Errorf("%w", err)
	}

	return yaml.YAMLToJSONStrict(data)
}

func MoveFile(src, dst string, dirPerms os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		if ferr := in.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		if ferr := out.Close(); ferr != nil {
			logger.StdLog.Fatal().Err(ferr).Msg("Failed to close file")
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = out.Sync(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.Chmod(dst, dirPerms); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.Remove(src); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
