package conf

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type Logging struct {
	Level     int    `yaml:"level"`
	File_Name string `yaml:"filename"`
	Log_Num   int    `yaml:"log_num"`
	Max_Size  int    `yaml:"max_size"`
	Compress  bool   `yaml:"compress"`
}

type CpuMonitoring struct {
	Threshold              int           `yaml:"threshold"`
	Stack_Trace_Collection bool          `yaml:"stack_trace_collection"`
	Flame_Graph_Collection bool          `yaml:"flame_graph_collection"`
	Collection_Interval    time.Duration `yaml:"collection_interval"`
	Output_Dir             string        `yaml:"output_dir"`
	Bin_Dir                string        `yaml:"bin_dir"`
}

type ProcessMonitor struct {
	ProcessName   string        `yaml:"process_name"`
	ProcessType   string        `yaml:"process_type"`
	CpuMonitoring CpuMonitoring `yaml:"cpu_monitoring"`
}

type Reporting struct {
	Central_Server string `yaml:"central_server"`
	Port           int    `yaml:"port"`
}

type PerfCollectorConf struct {
	Log            Logging          `yaml:"logging"`
	ProcessMonitor []ProcessMonitor `yaml:"process_monitoring"`
	Report         Reporting        `yaml:"report"`
}

var (
	GlobalConfig PerfCollectorConf
	mu           sync.RWMutex
	Logger       = logrus.New()
)

func LoadConfig(config_path string) error {
	file, err := os.Open(config_path)
	if err != nil {
		return err
	}
	defer file.Close()

	mu.Lock()
	defer mu.Unlock()
	decoder := yaml.NewDecoder(file)
	//perfCollectorConf := PerfCollectorConf{}
	if err := decoder.Decode(&GlobalConfig); err != nil {
		return err
	}
	return nil
}

func GetConfig() *PerfCollectorConf {
	mu.Lock()
	defer mu.Unlock()
	return &GlobalConfig
}

func ReloadLogConfig(logging *Logging) {
	mu.Lock()
	defer mu.Unlock()
	InitLogger(logging)
}

func InitLogger(logging *Logging) {
	//Logger := logrus.New()
	// logger config
	Logger.SetOutput(&lumberjack.Logger{
		Filename:   logging.File_Name,
		MaxSize:    logging.Max_Size,
		Compress:   logging.Compress,
		MaxBackups: logging.Log_Num,
	})

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetLevel(logrus.Level(logging.Level))
}

func GetLogger() *logrus.Logger {
	mu.Lock()
	defer mu.Unlock()
	return Logger
}
