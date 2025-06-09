package cdq

import (
	"strconv"
)

type FlagMap map[string]*FlagMapItem

type FlagMapItem struct {
	Value string
}

func (f FlagMap) String(long string) string {
	i := f[long]
	if i == nil {
		return ""
	}
	return i.Value
}

func (f FlagMap) Bool(long string) bool {
	i := f[long]
	if i == nil {
		return false
	}
	b, err := strconv.ParseBool(i.Value)
	if err != nil {
		return false
	}
	return b
}

func (f FlagMap) Int(long string) int {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.Atoi(i.Value)
	if err != nil {
		return 0
	}
	return v
}

func (f FlagMap) Int32(long string) int32 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseInt(i.Value, 10, 32)
	if err != nil {
		return 0
	}
	return int32(v)
}

func (f FlagMap) Int64(long string) int64 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseInt(i.Value, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func (f FlagMap) Uint(long string) uint {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseUint(i.Value, 10, 64)
	if err != nil {
		return 0
	}
	return uint(v)
}

func (f FlagMap) Uint32(long string) uint32 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseUint(i.Value, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

func (f FlagMap) Uint64(long string) uint64 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseUint(i.Value, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func (f FlagMap) Float32(long string) float32 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseFloat(i.Value, 32)
	if err != nil {
		return 0
	}
	return float32(v)
}

func (f FlagMap) Float64(long string) float64 {
	i := f[long]
	if i == nil {
		return 0
	}
	v, err := strconv.ParseFloat(i.Value, 64)
	if err != nil {
		return 0
	}
	return v
}
