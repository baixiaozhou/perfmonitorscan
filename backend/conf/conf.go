package conf

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Logging struct {
	Level     int    `yaml:"level"`
	File_Name string `yaml:"filename"`
	Log_Num   int    `yaml:"log_num"`
	Max_Size  int    `yaml:"max_size"`
	Compress  bool   `yaml:"compress"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	DbName   string `yaml:"dbname"`
}

type BackendConfig struct {
	Logging Logging  `yaml:"logging"`
	DB      DBConfig `yaml:"db"`
}

var (
	mu           sync.RWMutex
	GlobalConfig *BackendConfig
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

func ReloadLogConfig(logging *Logging) {
	mu.Lock()
	defer mu.Unlock()
	InitLogger(logging)
}
