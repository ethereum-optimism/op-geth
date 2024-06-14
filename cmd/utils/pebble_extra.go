package utils

import (
	"time"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethdb/pebble"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/urfave/cli/v2"
)

var (
	PebbleBytesPerSyncFlag = &cli.IntFlag{
		Name:     "pebble.bytes-per-sync",
		Category: flags.PebbleCategory,
	}
	PebbleL0CompactionFileThresholdFlag = &cli.IntFlag{
		Name:     "pebble.l0-compaction-file-threshold",
		Category: flags.PebbleCategory,
	}
	PebbleL0CompactionThresholdFlag = &cli.IntFlag{
		Name:     "pebble.l0-compaction-threshold",
		Category: flags.PebbleCategory,
	}
	PebbleL0StopWritesThresholdFlag = &cli.IntFlag{
		Name:     "pebble.l0-stop-writes-threshold",
		Category: flags.PebbleCategory,
	}
	PebbleLBaseMaxBytesFlag = &cli.Int64Flag{
		Name:     "pebble.l-base-max-bytes",
		Category: flags.PebbleCategory,
	}
	PebbleMemTableStopWritesThresholdFlag = &cli.IntFlag{
		Name:     "pebble.mem-table-stop-writes-threshold",
		Category: flags.PebbleCategory,
	}
	PebbleMaxConcurrentCompactionsFlag = &cli.IntFlag{
		Name:     "pebble.max-concurrent-compactions",
		Category: flags.PebbleCategory,
	}
	PebbleDisableAutomaticCompactionsFlag = &cli.BoolFlag{
		Name:     "pebble.disable-automatic-compactions",
		Category: flags.PebbleCategory,
	}
	PebbleWALBytesPerSyncFlag = &cli.IntFlag{
		Name:     "pebble.wal-bytes-per-sync",
		Category: flags.PebbleCategory,
	}
	PebbleWALDirFlag = &cli.StringFlag{
		Name:     "pebble.wal-dir",
		Category: flags.PebbleCategory,
	}
	PebbleWALMinSyncIntervalFlag = &cli.DurationFlag{
		Name:     "pebble.wal-min-sync-interval",
		Category: flags.PebbleCategory,
	}
	PebbleTargetByteDeletionRateFlag = &cli.IntFlag{
		Name:     "pebble.target-byte-deletion-rate",
		Category: flags.PebbleCategory,
	}

	// TODO: PebbleLevelOptions

	// Experimental

	PebbleL0CompactionConcurrencyFlag = &cli.IntFlag{
		Name:     "pebble.l0-compaction-concurrency",
		Category: flags.PebbleCategory,
	}
	PebbleCompactionDebtConcurrencyFlag = &cli.Uint64Flag{
		Name:     "pebble.compaction-debt-concurrency",
		Category: flags.PebbleCategory,
	}
	PebbleReadCompactionRateFlag = &cli.Int64Flag{
		Name:     "pebble.read-compaction-rate",
		Category: flags.PebbleCategory,
	}
	PebbleReadSamplingMultiplierFlag = &cli.Int64Flag{
		Name:     "pebble.read-sampling-multiplier",
		Category: flags.PebbleCategory,
	}
	PebbleMaxWriterConcurrencyFlag = &cli.IntFlag{
		Name:     "pebble.max-writer-concurrency",
		Category: flags.PebbleCategory,
	}
	PebbleForceWriterParallelismFlag = &cli.BoolFlag{
		Name:     "pebble.force-writer-parallelism",
		Category: flags.PebbleCategory,
	}

	PebbleFlags = []cli.Flag{
		PebbleBytesPerSyncFlag,
		PebbleL0CompactionFileThresholdFlag,
		PebbleL0CompactionThresholdFlag,
		PebbleL0StopWritesThresholdFlag,
		PebbleLBaseMaxBytesFlag,
		PebbleMemTableStopWritesThresholdFlag,
		PebbleMaxConcurrentCompactionsFlag,
		PebbleDisableAutomaticCompactionsFlag,
		PebbleWALBytesPerSyncFlag,
		PebbleWALDirFlag,
		PebbleWALMinSyncIntervalFlag,
		PebbleTargetByteDeletionRateFlag,
		// Experimental
		PebbleL0CompactionConcurrencyFlag,
		PebbleCompactionDebtConcurrencyFlag,
		PebbleReadCompactionRateFlag,
		PebbleReadSamplingMultiplierFlag,
		PebbleMaxWriterConcurrencyFlag,
		PebbleForceWriterParallelismFlag,
	}
)

func setPebbleExtraOptions(ctx *cli.Context, cfg *ethconfig.Config) {
	peos := new(pebble.ExtraOptions)

	if flag := PebbleBytesPerSyncFlag.Name; ctx.IsSet(flag) {
		peos.BytesPerSync = ctx.Int(flag)
	}
	if flag := PebbleL0CompactionFileThresholdFlag.Name; ctx.IsSet(flag) {
		peos.L0CompactionFileThreshold = ctx.Int(flag)
	}
	if flag := PebbleL0CompactionThresholdFlag.Name; ctx.IsSet(flag) {
		peos.L0CompactionThreshold = ctx.Int(flag)
	}
	if flag := PebbleL0StopWritesThresholdFlag.Name; ctx.IsSet(flag) {
		peos.L0StopWritesThreshold = ctx.Int(flag)
	}
	if flag := PebbleLBaseMaxBytesFlag.Name; ctx.IsSet(flag) {
		peos.LBaseMaxBytes = ctx.Int64(flag)
	}
	if flag := PebbleMemTableStopWritesThresholdFlag.Name; ctx.IsSet(flag) {
		peos.MemTableStopWritesThreshold = ctx.Int(flag)
	}
	if flag := PebbleMaxConcurrentCompactionsFlag.Name; ctx.IsSet(flag) {
		peos.MaxConcurrentCompactions = func() int { return ctx.Int(flag) }
	}
	if flag := PebbleDisableAutomaticCompactionsFlag.Name; ctx.IsSet(flag) {
		peos.DisableAutomaticCompactions = ctx.Bool(flag)
	}
	if flag := PebbleWALBytesPerSyncFlag.Name; ctx.IsSet(flag) {
		peos.WALBytesPerSync = ctx.Int(flag)
	}
	if flag := PebbleWALDirFlag.Name; ctx.IsSet(flag) {
		peos.WALDir = ctx.String(flag)
	}
	if flag := PebbleWALMinSyncIntervalFlag.Name; ctx.IsSet(flag) {
		peos.WALMinSyncInterval = func() time.Duration { return ctx.Duration(flag) }
	}
	if flag := PebbleTargetByteDeletionRateFlag.Name; ctx.IsSet(flag) {
		peos.TargetByteDeletionRate = ctx.Int(flag)
	}

	// Experimental

	if flag := PebbleL0CompactionConcurrencyFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.L0CompactionConcurrency = ctx.Int(flag)
	}
	if flag := PebbleCompactionDebtConcurrencyFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.CompactionDebtConcurrency = ctx.Uint64(flag)
	}
	if flag := PebbleReadCompactionRateFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.ReadCompactionRate = ctx.Int64(flag)
	}
	if flag := PebbleReadSamplingMultiplierFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.ReadSamplingMultiplier = ctx.Int64(flag)
	}
	if flag := PebbleMaxWriterConcurrencyFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.MaxWriterConcurrency = ctx.Int(flag)
	}
	if flag := PebbleForceWriterParallelismFlag.Name; ctx.IsSet(flag) {
		peos.Experimental.ForceWriterParallelism = ctx.Bool(flag)
	}

	cfg.PebbleExtraOptions = peos
}
