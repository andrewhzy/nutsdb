package dbmetrics

import (
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
	once sync.Once
	dbm  dbMetrics
)

func Init()                    { once.Do(func() { dbm = make(dbMetrics) }) }
func GetZeroMetrics() *Metrics { return &Metrics{0, 0, 0, 0} }

func DeleteMetrics(fd int)                   { delete(dbm, int32(fd)) }
func PutMetrics(fd int, m Metrics)           { dbm[int32(fd)] = m }
func GetMetrics(fd int) (m Metrics, ok bool) { m, ok = dbm[int32(fd)]; return }

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
