package dbmetrics

import (
	"math"
	"sync"
)

type (
	Metrics struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int64
		invalidBytes   int64
	}
	dbMetrics map[int32]Metrics
)

var (
	once   sync.Once
	dbm    dbMetrics
	locks  []sync.RWMutex
	n, mod int
)

func Init() { initiate(8) }
func initiate(lockNum int) {
	once.Do(func() {
		dbm = make(dbMetrics)
		n = 1 << int(math.Log2(float64(lockNum)))
		mod = n - 1
		locks = make([]sync.RWMutex, n)
		for i := 0; i < n; i++ {
			locks[i] = sync.RWMutex{}
		}
	})
}
func GetZeroMetrics() *Metrics { return &Metrics{0, 0, 0, 0} }

func DeleteMetrics(fd int) {
	for !locks[fd&mod].TryLock() {
	}
	delete(dbm, int32(fd))
	locks[fd&mod].Unlock()

}

func PutMetrics(fd int, m *Metrics) {
	for !locks[fd&mod].TryLock() {
	}
	dbm[int32(fd)] = *m
	locks[fd&mod].Unlock()

}

func GetMetrics(fd int) (m Metrics, ok bool) {
	for !locks[fd&mod].TryRLock() {
	}
	m, ok = dbm[int32(fd)]
	locks[fd&mod].RUnlock()

	return
}

func (m *Metrics) UpdateValid(entriesChange, bytesChange int) {
	m.validEntries += int32(entriesChange)
	m.validBytes += int64(bytesChange)
}

func (m *Metrics) UpdateInvalid(entriesChange, bytesChange int) {
	m.invalidEntries += int32(entriesChange)
	m.invalidBytes += int64(bytesChange)
}

func (m *Metrics) Update(validEntriesChange, invalidEntriesChange, validBytesChange, invalidBytesChange int) {
	m.validEntries += int32(validEntriesChange)
	m.invalidEntries += int32(invalidEntriesChange)
	m.validBytes += int64(validBytesChange)
	m.invalidBytes += int64(invalidBytesChange)
}

func (m *Metrics) UpdateMetrics(change Metrics) {
	m.validEntries += change.validEntries
	m.invalidEntries += change.invalidEntries
	m.validBytes += change.validBytes
	m.invalidBytes += change.invalidBytes
}
