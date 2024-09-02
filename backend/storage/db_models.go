package storage

type MonitoringCpuData struct {
	HostIp             string
	HostName           string
	Time               string
	Threshold          int
	ProcCpuPercent     float64
	ProcTopInfo        string
	ProcType           string
	StackInfo          string
	StackFilePath      string
	FlameGraphFilePath string
}
