package storage

import "syscall"

func DiskUsage(path string) (uint64, uint64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, 0
	}
	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	return all, free
}

const (
	_ = 1.0 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)
