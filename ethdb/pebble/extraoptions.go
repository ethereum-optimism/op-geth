package pebble

import "time"

type ExtraOptions struct {
	BytesPerSync                int
	L0CompactionFileThreshold   int
	L0CompactionThreshold       int
	L0StopWritesThreshold       int
	LBaseMaxBytes               int64
	MemTableStopWritesThreshold int
	MaxConcurrentCompactions    func() int
	DisableAutomaticCompactions bool
	WALBytesPerSync             int
	WALDir                      string
	WALMinSyncInterval          func() time.Duration
	TargetByteDeletionRate      int
	Experimental                ExtraOptionsExperimental
	Levels                      []ExtraLevelOptions
}

type ExtraOptionsExperimental struct {
	L0CompactionConcurrency   int
	CompactionDebtConcurrency uint64
	ReadCompactionRate        int64
	ReadSamplingMultiplier    int64
	MaxWriterConcurrency      int
	ForceWriterParallelism    bool
}

type ExtraLevelOptions struct {
	BlockSize      int
	IndexBlockSize int
	TargetFileSize int64
}
