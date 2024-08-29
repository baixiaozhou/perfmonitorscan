package collector

import (
	"encoding/json"
	"fmt"
	"github.com/baixiaozhou/backend/conf"
	"github.com/shirou/gopsutil/v3/process"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	C      = "c"
	PYTHON = "python"
	JAVA   = "java"
	GOLANG = "golang"
)

const (
	TIME_FLAG    = "TIME"
	COMMAND_FLAG = "COMMAND"
)

type ProcessThreadInfo struct {
	PID      int
	USER     string
	PRIORITY int
	NICE     int
	VIRT     string
	RES      string
	SHR      string
	STATE    string
	CPU_Per  string
	MEM_Per  string
	TIME     string
	COMMAND  string
}

func Monitor(monitors []conf.ProcessMonitor, worker_threads int) {
	// channel to pass services work
	//monitorChannel := make(chan conf.ProcessMonitor, len(monitors))

	// Done channel to signal when to stop monitoring
	doneChannel := make(chan bool)

	logger := conf.GetLogger()
	logger.Info("Monitor start")

	jsonData, _ := json.Marshal(monitors)
	logger.Info(string(jsonData))

	for _, monitor := range monitors {
		logger.Info("Monitor process monitor:", monitor)
		go monitorService(monitor, doneChannel)
	}
	//for i := 0; i < worker_threads; i++ {
	//	go monitorWorker(monitorChannel, doneChannel)
	//}
}

func monitorWorker(monitors <-chan conf.ProcessMonitor, doneChannel chan bool) {
	logger := conf.GetLogger()
	logger.Info("Monitor worker start")
	jsonData, _ := json.Marshal(monitors)
	logger.Info(string(jsonData))
	for monitor := range monitors {
		logger.Info("Monitor process monitor:", monitor)
		go monitorService(monitor, doneChannel)
	}
}

func monitorService(processMonitor conf.ProcessMonitor, doneChannel chan bool) {
	logger := conf.GetLogger()
	logger.Info("Monitor service start")

	ticker := time.NewTicker(processMonitor.CpuMonitoring.Collection_Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Debug("Monitor service tick start")
			//Get the CPU usage of a process
			pid, err := getPidByProcessName(processMonitor.ProcessName)
			if err != nil {
				return
			}
			p, err := process.NewProcess(int32(pid))
			if err != nil {
				conf.Logger.Errorf("monitor get process err:%v", err)
				return
			}
			cpuPercent, err := p.Percent(time.Second)
			if err != nil {
				conf.Logger.Errorf("monitor get process err:%v", err)
				return
			}
			if cpuPercent > float64(processMonitor.CpuMonitoring.Threshold) {
				conf.Logger.Debugf("the process :%v cpu usage is more than threshold, current:%v, threshold:%v",
					processMonitor.ProcessName, cpuPercent, processMonitor.CpuMonitoring.Threshold)
				processType := processMonitor.ProcessType
				// collect process thread
				procThreadInfos, err := getProcessTheadInfo(pid)
				if err != nil {
					conf.Logger.Errorf("get process thread info err:%v", err)
				}
				jsonData, _ := json.Marshal(procThreadInfos)
				conf.Logger.Debugf("get process thread info:%v", string(jsonData))
				// collect process info
				switch processType {
				case JAVA:
					if processMonitor.CpuMonitoring.Stack_Trace_Collection {
						stackInfo, fileName, err := CatchJavaStack(pid, processMonitor.CpuMonitoring)
						if err != nil {
							conf.Logger.Errorf("catch java stack err:%v", err)
						}
						jsonData, _ := json.Marshal(stackInfo)
						conf.Logger.Debugf("get process stack info:%v, fileName:%v", string(jsonData), fileName)

					}
					if processMonitor.CpuMonitoring.Flame_Graph_Collection {

					}
				default:

				}
			}

		case <-doneChannel:
			conf.Logger.Info("monitor service exit")
			return
		}
	}
}

func getPidByProcessName(processName string) (int, error) {
	out, err := exec.Command("pgrep", "-f", processName).Output()
	if err != nil {
		return 0, err
	}
	pids := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(pids) == 0 {
		conf.Logger.Warn("process '" + processName + "' not found")
		return 0, fmt.Errorf("process not found")
	}
	if len(pids) > 1 {
		conf.Logger.Warn("Multiple process IDs were found matching the name ‘%s’: %s. Using the first one.",
			processName, pids)
	}
	return strconv.Atoi(pids[0])
}

func getProcessTheadInfo(pid int) ([]ProcessThreadInfo, error) {
	var threadInfos []ProcessThreadInfo
	out, err := exec.Command("/usr/bin/top", "-Hp", strconv.Itoa(pid), "-n", "1", "-b").Output()
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	startFlag := false
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, TIME_FLAG) && strings.Contains(line, COMMAND_FLAG) {
			startFlag = true
			continue
		}
		if startFlag && (strings.TrimSpace(line) == "" || len(line) == 0) {
			break
		}
		currentThreadInfo := ProcessThreadInfo{}
		if startFlag {
			detail := strings.Fields(line)
			currentThreadInfo.PID, _ = strconv.Atoi(detail[0])
			currentThreadInfo.USER = detail[1]
			currentThreadInfo.PRIORITY, _ = strconv.Atoi(detail[2])
			currentThreadInfo.NICE, _ = strconv.Atoi(detail[3])
			currentThreadInfo.VIRT = detail[4]
			currentThreadInfo.RES = detail[5]
			currentThreadInfo.SHR = detail[6]
			currentThreadInfo.STATE = detail[7]
			currentThreadInfo.CPU_Per = detail[8]
			currentThreadInfo.MEM_Per = detail[9]
			currentThreadInfo.TIME = detail[10]
			currentThreadInfo.COMMAND = detail[11]
		}
		threadInfos = append(threadInfos, currentThreadInfo)
	}
	return threadInfos, nil
}
