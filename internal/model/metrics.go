package metrics

type Metric struct {
	ID    string   `json:"ID"`              //name metric
	MType string   `json:"type"`            //type metric
	Delta *int64   `json:"delta,omitempty"` //counter
	Value *float64 `json:"value,omitempty"` //gauge
}

type Metrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	PollCount     int64
	RandomValue   float64
}

func (g *Metrics) ToMap() map[string]interface{} {

	return map[string]interface{}{
		"Alloc":         g.Alloc,
		"BuckHashSys":   g.BuckHashSys,
		"Frees":         g.Frees,
		"GCCPUFraction": g.GCCPUFraction,
		"GCSys":         g.GCSys,
		"HeapAlloc":     g.HeapAlloc,
		"HeapIdle":      g.HeapIdle,
		"HeapInuse":     g.HeapInuse,
		"HeapObjects":   g.HeapObjects,
		"HeapReleased":  g.HeapReleased,
		"HeapSys":       g.HeapSys,
		"LastGC":        g.LastGC,
		"Lookups":       g.Lookups,
		"MCacheInuse":   g.MCacheInuse,
		"MCacheSys":     g.MCacheSys,
		"MSpanInuse":    g.MSpanInuse,
		"MSpanSys":      g.MSpanSys,
		"Mallocs":       g.Mallocs,
		"NextGC":        g.NextGC,
		"NumForcedGC":   g.NumForcedGC,
		"NumGC":         g.NumGC,
		"OtherSys":      g.OtherSys,
		"PauseTotalNs":  g.PauseTotalNs,
		"StackInuse":    g.StackInuse,
		"StackSys":      g.StackSys,
		"Sys":           g.Sys,
		"TotalAlloc":    g.TotalAlloc,
		"PollCount":     g.PollCount,
		"RandomValue":   g.RandomValue,
	}
}
