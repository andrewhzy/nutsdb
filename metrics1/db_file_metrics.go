package metrics2

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type FileMetrics struct {
	ValidEntries   int32
	InvalidEntries int32
	ValidBytes     int64
	InvalidBytes   int64
}

var (
	once  sync.Once
	fmMap map[int32]atomic.Value
	lock  sync.Mutex
)

func Init() {
	once.Do(func() {
		fmMap = make(map[int32]atomic.Value)
	})
}

// this method is for test case reset the state.
func reset() {
	fmMap = make(map[int32]atomic.Value)
}

func DeleteFileMetrics(fd int32) {
	delete(fmMap, fd)
}

// UpdateFileMetrics
// you can start with a new &FileMetrics{0,0,0,0},
// then update it along the way you update the DB entries,
// then update it back to fmMap with UpdateFileMetrics.
func UpdateFileMetrics(fd int32, update *FileMetrics) error {
	if m, ok := fmMap[fd]; ok {
		for {
			fmOld := m.Load().(FileMetrics)
			fmNew := FileMetrics{
				fmOld.ValidEntries + update.ValidEntries,
				fmOld.InvalidEntries + update.InvalidEntries,
				fmOld.ValidBytes + update.ValidBytes,
				fmOld.InvalidBytes + update.InvalidBytes,
			}
			if m.CompareAndSwap(fmOld, fmNew) {
				return nil
			}
		}
	}
	return errors.Errorf("FileMetrics for fd: %d dese not exist, please Initiate it", fd)
}

func InitFileMetricsForFd(fd int32) {
	lock.Lock()
	if m, ok := fmMap[fd]; !ok {
		fmMap[fd] = atomic.Value{}
		m = fmMap[fd]
		m.Store(FileMetrics{0, 0, 0, 0})
	}
	lock.Unlock()
}

//
//func GetFileMetrics(fd int32) (FileMetrics, bool) {
//
//	if m, ok := fmMap[fd]; ok {
//		return m.Load().(FileMetrics),ok
//	}
//	return , ok
//}

func CountFileMetrics() int {
	return len(fmMap)
}

func GetFDsExceedThreshold(threshold float64) []int32 {
	var fds []int32
	for fd, fm := range fmMap {
		m := fm.Load().(FileMetrics)
		ratio1 := float64(m.InvalidEntries) / float64(m.InvalidEntries+m.ValidEntries)
		ratio2 := float64(m.InvalidBytes) / float64(m.InvalidBytes+m.ValidBytes)
		if ratio1 >= threshold || ratio2 >= threshold {
			fds = append(fds, fd)
		}
	}
	sort.Slice(fds, func(i, j int) bool { return fds[i] < fds[j] })
	return fds
}
