package dbmetrics

import (
	"math"
	"sync"
	"sync/atomic"
)

type (
	Metrics struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int32
		invalidBytes   int32
	}
	dbMetrics map[int32]Metrics
)

var (
	once   sync.Once
	dbm    dbMetrics
	locks  []int64
	n, mod int
)

func Init() { initiate(18) }
func initiate(lockNum int) {
	once.Do(func() {
		dbm = make(dbMetrics)
		n = 1 << int(math.Log2(float64(lockNum)))
		locks = make([]int64, n)
		mod = n - 1
	})
}
func GetZeroMetrics() *Metrics { return &Metrics{0, 0, 0, 0} }

func DeleteMetrics(fd int) {
	for atomic.CompareAndSwapInt64(&locks[fd&mod], 0, -1) {
	}
	delete(dbm, int32(fd))
	locks[fd&mod] = 0

}

func PutMetrics(fd int, m *Metrics) {
	for atomic.CompareAndSwapInt64(&locks[fd&mod], 0, -1) {
	}
	dbm[int32(fd)] = *m
	locks[fd&mod] = 0

}

func GetMetrics(fd int) (m Metrics, ok bool) {
	for atomic.CompareAndSwapInt64(&locks[fd&mod], 0, -1) {
	}
	m, ok = dbm[int32(fd)]
	locks[fd&mod] = 0

	return
}

func (m *Metrics) UpdateValid(entriesChange, bytesChange int) {
	m.validEntries += int32(entriesChange)
	m.validBytes += int32(bytesChange)
}

func (m *Metrics) UpdateInvalid(entriesChange, bytesChange int) {
	m.invalidEntries += int32(entriesChange)
	m.invalidBytes += int32(bytesChange)
}

func (m *Metrics) Update(validEntriesChange, invalidEntriesChange, validBytesChange, invalidBytesChange int) {
	m.validEntries += int32(validEntriesChange)
	m.invalidEntries += int32(invalidEntriesChange)
	m.validBytes += int32(validBytesChange)
	m.invalidBytes += int32(invalidBytesChange)
}

func (m *Metrics) UpdateMetrics(change Metrics) {
	m.validEntries += change.validEntries
	m.invalidEntries += change.invalidEntries
	m.validBytes += change.validBytes
	m.invalidBytes += change.invalidBytes
}
