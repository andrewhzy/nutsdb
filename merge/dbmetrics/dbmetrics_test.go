package dbmetrics

import (
	"reflect"
	"testing"
)

func TestDeleteMetrics(t *testing.T) {
	Init()
	type args struct {
		fd int
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{3}},
		{"test", args{4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteMetrics(tt.args.fd)
		})
	}
}

func TestGetMetrics(t *testing.T) {
	Init()
	PutMetrics(1, Metrics{1, 1, 1, 1})
	PutMetrics(2, Metrics{2, 2, 2, 2})
	type args struct {
		fd int
	}
	tests := []struct {
		name   string
		args   args
		wantM  Metrics
		wantOk bool
	}{
		{"", args{1}, Metrics{1, 1, 1, 1}, true},
		{"", args{2}, Metrics{2, 2, 2, 2}, true},
		{"", args{3}, Metrics{0, 0, 0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, gotOk := GetMetrics(tt.args.fd)
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("GetMetrics() gotM = %v, want %v", gotM, tt.wantM)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetMetrics() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestInitMergeMetrics(t *testing.T) {
	Init()
	tests := []struct {
		name string
	}{
		{""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
		})
	}
}

func TestPutMetrics(t *testing.T) {
	Init()
	type args struct {
		fd int
		m  Metrics
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{1, Metrics{1, 1, 1, 1}}},
		{"", args{12, Metrics{2, 2, 2, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutMetrics(tt.args.fd, tt.args.m)
		})
	}
}

func TestMetrics_Update(t *testing.T) {
	type fields struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int64
		invalidBytes   int64
	}
	type args struct {
		validEntriesChange   int
		invalidEntriesChange int
		validBytesChange     int
		invalidBytesChange   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Metrics
	}{
		{"", fields{1, 1, 1, 1}, args{1, 1, 1, 1}, &Metrics{2, 2, 2, 2}},
		{"", fields{1, 1, 1, 1}, args{-1, -1, -1, -1}, &Metrics{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				validEntries:   tt.fields.validEntries,
				invalidEntries: tt.fields.invalidEntries,
				validBytes:     tt.fields.validBytes,
				invalidBytes:   tt.fields.invalidBytes,
			}
			if m.Update(tt.args.validEntriesChange, tt.args.invalidEntriesChange, tt.args.validBytesChange, tt.args.invalidBytesChange); !reflect.DeepEqual(m, tt.want) {
				t.Errorf("GetZeroMetrics() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestMetrics_UpdateInvalid(t *testing.T) {
	type fields struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int64
		invalidBytes   int64
	}
	type args struct {
		entries int
		bytes   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Metrics
	}{
		{"", fields{1, 1, 1, 1}, args{1, 1}, &Metrics{1, 2, 1, 2}},
		{"", fields{1, 1, 1, 1}, args{-1, -1}, &Metrics{1, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				validEntries:   tt.fields.validEntries,
				invalidEntries: tt.fields.invalidEntries,
				validBytes:     tt.fields.validBytes,
				invalidBytes:   tt.fields.invalidBytes,
			}
			if m.UpdateInvalid(tt.args.entries, tt.args.bytes); !reflect.DeepEqual(m, tt.want) {
				t.Errorf("GetZeroMetrics() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestMetrics_UpdateMetrics(t *testing.T) {
	type fields struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int64
		invalidBytes   int64
	}
	type args struct {
		change Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Metrics
	}{
		{"", fields{1, 1, 1, 1}, args{Metrics{1, 1, 1, 1}}, &Metrics{2, 2, 2, 2}},
		{"", fields{1, 1, 1, 1}, args{Metrics{-1, -1, -1, -1}}, &Metrics{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				validEntries:   tt.fields.validEntries,
				invalidEntries: tt.fields.invalidEntries,
				validBytes:     tt.fields.validBytes,
				invalidBytes:   tt.fields.invalidBytes,
			}
			if m.UpdateMetrics(tt.args.change); !reflect.DeepEqual(m, tt.want) {
				t.Errorf("GetZeroMetrics() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestMetrics_UpdateValid(t *testing.T) {
	type fields struct {
		validEntries   int32
		invalidEntries int32
		validBytes     int64
		invalidBytes   int64
	}
	type args struct {
		entries int
		bytes   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Metrics
	}{
		{"", fields{1, 1, 1, 1}, args{1, 1}, &Metrics{2, 1, 2, 1}},
		{"", fields{1, 1, 1, 1}, args{-1, -1}, &Metrics{0, 1, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				validEntries:   tt.fields.validEntries,
				invalidEntries: tt.fields.invalidEntries,
				validBytes:     tt.fields.validBytes,
				invalidBytes:   tt.fields.invalidBytes,
			}
			if m.UpdateValid(tt.args.entries, tt.args.bytes); !reflect.DeepEqual(m, tt.want) {
				t.Errorf("GetZeroMetrics() = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestGetZeroMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		{"", &Metrics{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetZeroMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZeroMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
