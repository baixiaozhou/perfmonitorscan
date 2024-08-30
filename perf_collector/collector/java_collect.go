package collector

import (
	"bufio"
	"fmt"
	"github.com/baixiaozhou/backend/conf"
	"github.com/baixiaozhou/backend/utils"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	STATUS_DEADLOCK        = "Java-level deadlock"
	STATUS_RUNNABLE        = "RUNNABLE"
	STATUS_TIMED_WAITING   = "TIMED_WAITING"
	STATUS_BLOCKED         = "BLOCKED"
	STATUS_WAITING_PARKING = "WAITING (parking)"
	STATUS_WAITING_MONITOR = "WAITING (on object monitor)"
)

const (
	JSTACK_FLAG  = "Full thread dump Java HotSpot"
	DAEMON       = "daemon"
	PRIO         = "prio"
	OS_PRIO      = "os_prio"
	CPU          = "cpu"
	ELAPSED      = "elapsed"
	TID          = "tid"
	NID          = "nid"
	THREAD_STATE = "java.lang.Thread.State"
)

const (
	JVMTHREADPATTERN  = `^"[^\"]+" os_prio=\d+`
	THREADINFOPATTERN = `^"[^\"]+" #\d+ (?:daemon )?prio=\d+ os_prio=\d+`
)

const (
	JSTACK      = "jstack"
	TIME_FORMAT = "20060102150405"
	ASPROF      = "asprof"
	FLAME       = "flame"
	HTML        = ".html"
	EXIT_200    = "exit status 200"
)

type ThreadInfo struct {
	Name         string
	Id           string
	Daemon       bool
	Priority     int
	Os_prio      int
	Cpu_time     string
	Elapsed_time string
	Tid          string
	Nid          string
	State        string
	Addr         string
	StackTrace   []string
}

type DeadLockInfo struct {
	Info []string
}

type JstackInfo struct {
	ThreadsCount        int
	DeadLockCount       int
	RunnableCount       int
	WaitingParkCount    int
	WaitingMonitorCount int
	BlockedCount        int
	TimedWaitCount      int
	OtherCount          int
	Threads             []ThreadInfo
	Deadlocks           []DeadLockInfo
}

func CatchJavaStack(pid int, monitoring conf.CpuMonitoring) (JstackInfo, string, error) {
	binDir := monitoring.Bin_Dir
	if binDir == "" {
		binDir = "/usr/bin/"
	}
	if err := utils.CreateDir(monitoring.Output_Dir); err != nil {
		return JstackInfo{}, "", err
	}
	fileName := monitoring.Output_Dir + "/" + JSTACK + "_" + strconv.Itoa(pid) + "_" + time.Now().Format(TIME_FORMAT)
	out, err := exec.Command(binDir+JSTACK, "-l", strconv.Itoa(pid)).Output()
	if err != nil {
		conf.Logger.Warn("get jstack error:", err, ", pid:", pid)
		conf.Logger.Debug("to kill -3 ", pid)
		// jstack err, kill -3 pid
		exec.Command("/usr/bin/kill", "-3", strconv.Itoa(pid)).Run()
	} else {
		if err := utils.SaveToFile(out, fileName); err != nil {
			conf.Logger.Warn("save jstack info to file failed,err:", err)
			conf.Logger.Warn("jstackinfo:", string(out))
		}
	}
	jstackInfo, err := GenerJstackInfo(fileName)
	if err != nil {
		return JstackInfo{}, "", err
	}
	return jstackInfo, fileName, nil
}

func ParseJstack(filePath string) ([]ThreadInfo, []DeadLockInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	var threads []ThreadInfo
	var deadlocks []DeadLockInfo
	var currentThreadInfo ThreadInfo
	var currentDeadLock DeadLockInfo
	commonRegex := regexp.MustCompile(THREADINFOPATTERN)
	jvmRegex := regexp.MustCompile(JVMTHREADPATTERN)
	scanner := bufio.NewScanner(file)
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		flag += 1
		if flag == 2 && !strings.Contains(line, JSTACK_FLAG) {
			return nil, nil, fmt.Errorf("Can not pasrse file: %s", filePath)
		}
		if commonRegex.MatchString(line) {
			parts := strings.Split(line, " ")
			currentThreadInfo = ThreadInfo{}
			currentThreadInfo.Name = parts[0]
			currentThreadInfo.Id = parts[1]
			for _, v := range parts {
				if v == DAEMON {
					currentThreadInfo.Daemon = true
				}
				if strings.HasPrefix(v, PRIO) {
					currentThreadInfo.Priority, _ = strconv.Atoi(strings.Split(v, "=")[1])
				}
				if strings.HasPrefix(v, OS_PRIO) {
					currentThreadInfo.Os_prio, _ = strconv.Atoi(strings.Split(v, "=")[1])
				}
				if strings.HasPrefix(v, CPU) {
					currentThreadInfo.Cpu_time = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, ELAPSED) {
					currentThreadInfo.Elapsed_time = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, TID) {
					currentThreadInfo.Tid = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, NID) {
					currentThreadInfo.Nid = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					currentThreadInfo.Addr = v
				}
			}
			currentThreadInfo.StackTrace = append(currentThreadInfo.StackTrace, line)
			scanner.Scan()
			line := scanner.Text()
			if strings.HasPrefix(strings.TrimSpace(line), THREAD_STATE) {
				currentThreadInfo.State = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" {
					threads = append(threads, currentThreadInfo)
					break
				}
				currentThreadInfo.StackTrace = append(currentThreadInfo.StackTrace, line)
			}
		} else if jvmRegex.MatchString(line) {
			parts := strings.Split(line, " ")
			currentThreadInfo = ThreadInfo{}
			currentThreadInfo.Name = parts[0]
			for _, v := range parts {
				if strings.HasPrefix(v, OS_PRIO) {
					currentThreadInfo.Os_prio, _ = strconv.Atoi(strings.Split(v, "=")[1])
				}
				if strings.HasPrefix(v, CPU) {
					currentThreadInfo.Cpu_time = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, ELAPSED) {
					currentThreadInfo.Elapsed_time = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, TID) {
					currentThreadInfo.Tid = strings.Split(v, "=")[1]
				}
				if strings.HasPrefix(v, NID) {
					currentThreadInfo.Nid = strings.Split(v, "=")[1]
				}
			}
			currentThreadInfo.StackTrace = append(currentThreadInfo.StackTrace, line)
			threads = append(threads, currentThreadInfo)
		}

		if strings.Contains(line, STATUS_DEADLOCK) {
			currentDeadLock = DeadLockInfo{}
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" {
					deadlocks = append(deadlocks, currentDeadLock)
					break
				}
				currentDeadLock.Info = append(currentDeadLock.Info, line)
			}
		}
	}
	return threads, deadlocks, scanner.Err()
}

func GenerJstackInfo(filePath string) (JstackInfo, error) {
	var generInfo JstackInfo
	threads, deadlocks, err := ParseJstack(filePath)
	if err != nil {
		return JstackInfo{}, err
	}
	generInfo.Deadlocks = deadlocks
	generInfo.Threads = threads
	generInfo.DeadLockCount = len(deadlocks)
	generInfo.ThreadsCount = len(threads)
	for _, thread := range generInfo.Threads {
		if thread.State == STATUS_RUNNABLE {
			generInfo.RunnableCount += 1
		} else if thread.State == STATUS_TIMED_WAITING {
			generInfo.TimedWaitCount += 1
		} else if thread.State == STATUS_BLOCKED {
			generInfo.BlockedCount += 1
		} else if thread.State == STATUS_WAITING_PARKING {
			generInfo.WaitingParkCount += 1
		} else if thread.State == STATUS_WAITING_MONITOR {
			generInfo.WaitingMonitorCount += 1
		} else {
			generInfo.OtherCount += 1
		}
	}
	return generInfo, nil
}

func CatchJavaFlameGraph(pid int, monitoring conf.CpuMonitoring) (fileName string, err error) {
	flame := monitoring.Flame_Graph_Collection
	asprof := flame.Bin_Dir + "/" + ASPROF
	if err := utils.CreateDir(monitoring.Output_Dir); err != nil {
		return "", err
	}
	file := monitoring.Output_Dir + "/" + FLAME + "_" + strconv.Itoa(pid) + "_" + time.Now().Format(TIME_FORMAT) + HTML
	cmd := exec.Command(
		asprof,
		"-d", flame.Collection_Duration.String(),
		"-e", "cpu",
		"-f", file,
		strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), EXIT_200) {
			exec.Command(asprof, "stop", strconv.Itoa(pid)).Run()
		}
		return "", err
	}
	return file, nil
}
