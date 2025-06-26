package backup

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func generateBackupFilename(config BackupConfig, dumper DatabaseDumper) string {
	timestamp := time.Now().Format("20060102_150405")
	
	dbName := config.DatabaseName
	if dbName == "" {
		dbName = "default"
	}

	extension := dumper.GetFileExtension()
	if config.Compression == "gz" && !strings.HasSuffix(extension, ".gz") {
		extension += ".gz"
	}

	return fmt.Sprintf("%s_%s_%s%s",
		config.DatabaseType,
		dbName,
		timestamp,
		extension,
	)
}

func logBackupSuccess(path string, size int64, dbType, storageType string) {
	log.Printf("✅ Backup completed successfully: %s", path)
	log.Printf("   Database: %s → Storage: %s", dbType, storageType)
	if size > 0 {
		log.Printf("   Size: %.2f MB", float64(size)/1024/1024)
	}
}