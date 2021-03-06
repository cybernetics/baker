package filter

import (
	"fmt"
	"sync/atomic"

	"github.com/AdRoll/baker"
)

// NotNullDesc describes the NotNull filter.
var NotNullDesc = baker.FilterDesc{
	Name:   "NotNull",
	New:    NewNotNull,
	Config: &NotNullConfig{},
	Help:   "Discard the records having null (i.e empty) fields.\n",
}

// NotNullConfig holds configuration parameters for the NotNull filter.
type NotNullConfig struct {
	Fields []string `help:"Fields is the list of fields to check for null/empty values" required:"true"`
}

// NotNull is a baker filter that discards records having null fields.
type NotNull struct {
	numProcessedLines int64
	numFilteredLines  int64
	cfg               *NotNullConfig
	fields            []baker.FieldIndex
}

// NewNotNull creates and configures a new NotNull filter.
func NewNotNull(cfg baker.FilterParams) (baker.Filter, error) {
	if cfg.DecodedConfig == nil {
		cfg.DecodedConfig = &NotNullConfig{}
	}
	dcfg := cfg.DecodedConfig.(*NotNullConfig)

	f := &NotNull{cfg: dcfg}
	for _, field := range dcfg.Fields {
		if val, ok := cfg.FieldByName(field); ok {
			f.fields = append(f.fields, val)
		} else {
			return nil, fmt.Errorf("unknown field %q", field)
		}
	}
	return f, nil
}

// Stats implements baker.Filter.
func (v *NotNull) Stats() baker.FilterStats {
	return baker.FilterStats{
		NumProcessedLines: atomic.LoadInt64(&v.numProcessedLines),
		NumFilteredLines:  atomic.LoadInt64(&v.numFilteredLines),
	}
}

// Process implements baker.Filter.
func (v *NotNull) Process(l baker.Record, next func(baker.Record)) {
	atomic.AddInt64(&v.numProcessedLines, 1)
	for _, field := range v.fields {
		if l.Get(field) == nil {
			atomic.AddInt64(&v.numFilteredLines, 1)
			return
		}
	}

	next(l)
}
