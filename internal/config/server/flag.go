package server

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
	"log"
	"os"
	"strconv"
)

var AddrServer = &config.NetAddr{
	Host: "localhost",
	Port: 8080,
}

var StorageInterval int
var FilePathStorage string
var Restore bool
var DatabaseDSN string

func ParseFlagServer() {
	flag.Var(AddrServer, "a", "localhost:8080")
	flag.StringVar(&FilePathStorage, "f", "/tmp/metrics-db.json", "the file path to save the storage")
	flag.IntVar(&StorageInterval, "i", 300, "interval save storage to file")
	flag.BoolVar(&Restore, "r", true, "upload previously saved ones")
	flag.StringVar(&DatabaseDSN, "d", "host=localhost port=5432 user=postgres password=postgres dbname=metrics sslmode=disable", "database dsn")
	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		err := AddrServer.Set(addr)
		if err != nil {
			log.Print(err)
		}
	}
	if storeInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInt, err := strconv.Atoi(storeInterval)
		if err != nil {
			log.Print(err)
		}
		StorageInterval = storeInt
	}

	if EnvDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		DatabaseDSN = EnvDatabaseDSN
	}

	if filePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		FilePathStorage = filePath
	}
	if restore, ok := os.LookupEnv("RESTORE"); ok {
		restBool, err := strconv.ParseBool(restore)
		if err != nil {
			log.Print(err)
		}
		Restore = restBool
	}

}
