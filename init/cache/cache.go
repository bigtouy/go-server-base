package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"time"

	"go-server-base/global"
	"go-server-base/init/cache/badger_db"
)

func Init() {
	c := global.CONF.System.Cache

	options := badger.Options{
		Dir:                c,
		ValueDir:           c,
		ValueLogFileSize:   102400000,
		ValueLogMaxEntries: 100000,
		VLogPercentile:     0.1,

		MemTableSize:                  64 << 20,
		BaseTableSize:                 2 << 20,
		BaseLevelSize:                 10 << 20,
		TableSizeMultiplier:           2,
		LevelSizeMultiplier:           10,
		MaxLevels:                     7,
		NumGoroutines:                 8,
		MetricsEnabled:                true,
		NumCompactors:                 4,
		NumLevelZeroTables:            5,
		NumLevelZeroTablesStall:       15,
		NumMemtables:                  5,
		BloomFalsePositive:            0.01,
		BlockSize:                     4 * 1024,
		SyncWrites:                    false,
		NumVersionsToKeep:             1,
		CompactL0OnClose:              false,
		VerifyValueChecksum:           false,
		BlockCacheSize:                256 << 20,
		IndexCacheSize:                0,
		ZSTDCompressionLevel:          1,
		EncryptionKey:                 []byte{},
		EncryptionKeyRotationDuration: 10 * 24 * time.Hour, // Default 10 days.
		DetectConflicts:               true,
		NamespaceOffset:               -1,
	}

	cache, err := badger.Open(options)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	global.CACHE = badger_db.NewCacheDB(cache)
	global.LOG.Info("init cache successfully")
}
