package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	metrics "github.com/kornetvba/metrics-service/internal/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileStorage struct {
	Storage
	FilePath            string
	StorageIntervalSave int
	Restore             bool
}

func NewFileStorage(FilePath string, StorageInterval int, Restore bool, storage Storage) *FileStorage {
	return &FileStorage{
		Storage:             storage,
		FilePath:            FilePath,
		StorageIntervalSave: StorageInterval,
		Restore:             Restore,
	}
}

func (f *FileStorage) Save() error { //open create and write data
	directories := filepath.Dir(f.FilePath)
	if len(directories) > 2 {
		err := os.MkdirAll(directories, 0777)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(f.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	counter, gauge := f.GetAllMetrics()
	if len(counter) != 0 {
		for k, v := range counter {
			metric := metrics.Metric{ID: k, MType: "counter", Delta: &v}
			err = json.NewEncoder(file).Encode(metric)
			if err != nil {
				return err
			}
		}
	}
	if len(gauge) != 0 {
		for k, v := range gauge {
			metric := metrics.Metric{ID: k, MType: "gauge", Value: &v}
			err = json.NewEncoder(file).Encode(metric)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (f *FileStorage) Load() error { //read file run, and read all data file
	if !f.Restore {
		return errors.New("FileLoad off")
	}

	file, err := os.OpenFile(f.FilePath, os.O_RDONLY, 0777)

	if err != nil {
		return err
	}

	bufferReader := bufio.NewScanner(file)
	for bufferReader.Scan() {
		var metric = metrics.Metric{}
		err = json.Unmarshal(bufferReader.Bytes(), &metric)
		if err != nil {
			return err
		}
		if metric.MType == "gauge" {
			f.SetGauge(metric.ID, *metric.Value)
		}
		if metric.MType == "counter" {
			f.UpdateCounter(metric.ID, *metric.Delta)
		}

	}

	return err
}

func (f *FileStorage) SaveToFile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(writer, request)
		if f.StorageIntervalSave == 0 {
			err := f.Save()
			if err != nil {
				log.Print(err)
			}
		}
	})
}
