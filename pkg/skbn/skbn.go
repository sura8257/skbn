package skbn

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
)

// Copy copies files from src to dst
func Copy(src, dst string, parallel int, bufferSize int64) error {

	sourceURL, err := url.Parse(src)
	if err != nil {
		return err
	}

	destURL, err := url.Parse(dst)
	if err != nil {
		return err
	}

	if sourceURL.Scheme != "s3" && destURL.Scheme != "s3" {
		err = errors.New("usage: skbn cp <LocalPath> or <S3Uri> <LocalPath> or <S3Uri>")
		return err
	}

	if sourceURL.Scheme != "s3" {
		_, err = checkPath(sourceURL.Path)

		if err != nil {
			return err
		}
	}

	if sourceURL.Scheme != "s3" {
		err = copyFileToS3(src, dst, parallel, bufferSize)
		if err != nil {
			return err
		}
	} else {
		err = copyS3ToFile(src, dst, parallel, bufferSize)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check filesystem path
func checkPath(path string) (isDir bool, err error) {
	absPath, err := filepath.Abs(path)
	if err == nil {
		info, err := os.Lstat(absPath)
		if err == nil {
			return info.IsDir(), err
		}
	}
	return false, errors.New("The user provided path " + path + " does not exist")
}
