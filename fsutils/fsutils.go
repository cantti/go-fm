package fsutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type DstExistsAction int

const (
	DstExistsActionOverWrite DstExistsAction = iota
	DstExistsActionSkip
)

func Copy(src, dst string, onDstExists func() DstExistsAction) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get stat : %w", err)
	}
	_, err = os.Stat(dst)
	if err == nil {
		if onDstExists != nil {
			action := onDstExists()
			if action == DstExistsActionOverWrite {
				os.RemoveAll(dst)
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	if srcStat.IsDir() {
		entries, _ := readDirRecursively(src, dst, "/")
		for _, e := range entries {
			srcStat, err := os.Stat(e.src)
			if err != nil {
				return fmt.Errorf("failed to get stat : %w", err)
			}
			_, dstStatErr := os.Stat(e.dst)
			if dstStatErr == nil {
				return fmt.Errorf("file already exists : %w", err)
			}
			if srcStat.IsDir() {
				// dir needs to be writable
				os.Mkdir(e.dst, os.ModePerm)
			} else {
				copyFile(e.src, e.dst)
			}
			// set perm for directories
			for _, e := range entries {
				stat, err := os.Stat(e.src)
				if err != nil {
					return fmt.Errorf("failed to get stat : %w", err)
				}
				if stat.IsDir() {
					os.Chmod(e.dst, stat.Mode().Perm())
				} else {
					continue
				}
			}
		}
		return nil
	} else {
		copyFile(src, dst)
		return nil
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	os.Chmod(dst, sourceFileStat.Mode().Perm())

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

type entry struct {
	src      string
	dst      string
	relative string
}

func readDirRecursively(srcBase string, dstBase string, relativePath string) ([]entry, error) {
	var result []entry
	entries, _ := os.ReadDir(filepath.Join(srcBase, relativePath))
	result = append(result, entry{
		src:      filepath.Join(srcBase, relativePath),
		dst:      filepath.Join(dstBase, relativePath),
		relative: relativePath})
	for _, e := range entries {
		if e.IsDir() {
			children, _ := readDirRecursively(srcBase, dstBase, filepath.Join(relativePath, e.Name()))
			result = append(result, children...)
		} else {

			result = append(result, entry{
				src:      filepath.Join(srcBase, relativePath, e.Name()),
				dst:      filepath.Join(dstBase, relativePath, e.Name()),
				relative: relativePath})
		}
	}
	return result, nil
}
