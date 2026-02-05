package psql

import (
	"database/sql"
	"fmt"
	metrics "github.com/kornetvba/metrics-service/internal/model"
	_ "github.com/lib/pq"
	"strconv"
)

var DB *sql.DB

type DatabasePSQL struct {
	conn *sql.DB
}

func NewDatabasePSQL(db *sql.DB) *DatabasePSQL {
	return &DatabasePSQL{conn: db}
}

func (db *DatabasePSQL) BootStrap() error {

	_, err := db.conn.Exec(
		`
		CREATE TABLE IF NOT EXISTS metrics(
    	id SERIAL PRIMARY KEY,
    	type_metric VARCHAR(20) NOT NULL,
    	name_metric VARCHAR(200) NOT NULL UNIQUE,
    	value_metric NUMERIC);
`)
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabasePSQL) UpdateCounter(name string, value int64) (int64, error) {
	var v []byte
	err := db.conn.QueryRow(
		`
		INSERT INTO metrics (type_metric, name_metric,value_metric)
		VALUES($1,$2,$3)
		ON CONFLICT (name_metric)
		DO UPDATE 
		SET value_metric=metrics.value_metric + $3
		   RETURNING metrics.value_metric;
		   `,
		"counter", name, value).Scan(&v)

	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (db *DatabasePSQL) SetGauge(name string, value float64) (float64, error) {
	var val []byte

	_ = db.conn.QueryRow(
		`
		INSERT INTO metrics (type_metric, name_metric, value_metric)
		VALUES ($1,$2,$3)
		ON CONFLICT (name_metric)
		DO UPDATE
		   SET value_metric = $3
		   RETURNING metrics.value_metric
	
`, "gauge", name, value,
	).Scan(&val)

	valFloat, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		return 0, err
	}

	return valFloat, nil
}

func (db *DatabasePSQL) GetMetric(typeMetric string, nameMetric string) (interface{}, error) {
	row := db.conn.QueryRow(
		`
		SELECT (value_metric) 
		FROM metrics
		WHERE (type_metric=$1) AND (name_metric=$2)
				`, typeMetric, nameMetric,
	)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}

	switch typeMetric {
	case "counter":
		val, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	case "gauge":
		val, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return nil, err
		}
		return val, nil

	}

	return nil, fmt.Errorf("uncnown error")

}

func (db *DatabasePSQL) GetAllMetrics() (map[string]int64, map[string]float64) {
	counterMap := make(map[string]int64)
	gaugeMap := make(map[string]float64)

	rows, err := db.conn.Query(`
		SELECT type_metric, name_metric, value_metric FROM metrics;
`)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	for rows.Next() {
		var metric MetricPSQL
		var val []byte

		err = rows.Scan(&metric.MType, &metric.ID, &val)
		if err != nil {
			return nil, nil
		}

		switch metric.MType {
		case "counter":
			v, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return nil, nil
			}
			vl := int64(v)
			counterMap[metric.ID] = vl
		case "gauge":
			v, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				return nil, nil
			}
			gaugeMap[metric.ID] = v
		}

	}

	if err = rows.Err(); err != nil {
		return nil, nil
	}

	return counterMap, gaugeMap
}

func (db *DatabasePSQL) AppendMetrics(metrics []metrics.Metric) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmtCounter, err := tx.Prepare(
		`
			INSERT INTO metrics (type_metric, name_metric, value_metric)
			VALUES ($1,$2,$3)
			ON CONFLICT (name_metric)
			DO UPDATE
			   SET value_metric = metrics.value_metric + $3;`,
	)
	if err != nil {
		return err
	}
	defer stmtCounter.Close()
	stmtGauge, err := tx.Prepare(
		`
			INSERT INTO metrics (type_metric, name_metric, value_metric) 
			VALUES ($1,$2,$3)
			ON CONFLICT (name_metric)
			DO UPDATE
			   SET value_metric = $3
			`,
	)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()

	for _, metric := range metrics {
		if metric.ID == "" || metric.MType == "" {
			return fmt.Errorf("missing args")
		}

		if metric.Value == nil && metric.Delta == nil {
			return fmt.Errorf("missing args")
		}
		switch metric.MType {
		case "counter":
			_, err = stmtCounter.Exec(metric.MType, metric.ID, *metric.Delta)
			if err != nil {
				tx.Rollback()
				return err
			}
		case "gauge":
			_, err = stmtGauge.Exec(metric.MType, metric.ID, *metric.Value)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}
