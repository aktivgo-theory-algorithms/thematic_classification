package config

import "runtime"

var (
	InputPath       = "../lenta-ru-news.csv"
	OutputPath      = "dataset.csv"
	GoroutinesCount = runtime.NumCPU()
)
