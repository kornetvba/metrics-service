package backup

import (
	"bufio"
	"encoding/json"
	metrics "github.com/kornetvba/metrics-service/internal/model"
	"github.com/kornetvba/metrics-service/internal/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Backup interface {
	Save() error
	Load() error
}

type BackupManager struct {
	store           repository.Storage
	Restore         bool
	StorageInterval int
	FilePath        string
}

func NewBackupManager(store repository.Storage, restore bool, storageInterval int, filePath string) *BackupManager {
	return &BackupManager{store: store, Restore: restore, StorageInterval: storageInterval, FilePath: filePath}
}

func (bm *BackupManager) Load() error {
	file, err := os.Open(bm.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffioReader := bufio.NewScanner(file)
	var metric metrics.Metric

	for buffioReader.Scan() {
		err := json.Unmarshal(buffioReader.Bytes(), &metric)
		if err != nil {
			return err
		}
		switch metric.MType {
		case "counter":
			_, err := bm.store.UpdateCounter(metric.ID, *metric.Delta)
			if err != nil {
				return err
			}
		case "gauge":
			_, err := bm.store.SetGauge(metric.ID, *metric.Value)
			if err != nil {
				return err
			}

		}
	}
	return buffioReader.Err()

}

func (bm *BackupManager) Save() error {
	directories := filepath.Dir(bm.FilePath)
	if len(directories) > 2 {
		err := os.MkdirAll(directories, 0777)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(bm.FilePath, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	counters, guages := bm.store.GetAllMetrics()
	var metric metrics.Metric

	for k, v := range counters {
		metric.MType, metric.ID, metric.Delta = "counter", k, &v

		if err := json.NewEncoder(file).Encode(metric); err != nil {
			return err
		}
	}
	for k, v := range guages {
		var metric metrics.Metric

		metric.MType, metric.ID, metric.Value = "gauge", k, &v

		if err := json.NewEncoder(file).Encode(metric); err != nil {
			return err
		}
	}
	return nil

}

func (bm *BackupManager) SaveFileSync(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if bm.StorageInterval == 0 {
			err := bm.Save()
			if err != nil {
				log.Print(err)
			}
		}
	})
}
