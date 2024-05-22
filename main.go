package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	ps "github.com/mitchellh/go-ps"
	gopsutilnet "github.com/shirou/gopsutil/net"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	args := os.Args[1:]
	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("Error fetching processes:", err)
		return
	}

	results := aggregateResults(args, processes)

	if len(results) == 0 {
		fmt.Println("No matching processes found.")
	} else {
		for _, result := range results {
			printProcess(result)
		}
	}
}

func printHelp() {
	fmt.Println(`
  Syntax: pfinder <arguments>...
  
  Arguments can be:
    - a path    If argument is a path to an existing file or folder then
                whatever PID locking the resource (if any) will be reported.
    - a number  Report for such PID
    - a string  Running processes commands will be filtered case insensitive
                for containing the string
    - a regex   As above, in the string basic regex interpolation is supported.
                Regex, not wildcards: use "." (not "?") or ".+" (not "*")
    - a port    If argument is a port (prefixed with ':'), the processes 
                listening on that port will be reported.

  Multiple arguments can be provided, and the results will be aggregated.
  Should the argument contain whitespaces wrap it within "..."
`)
}

func aggregateResults(args []string, processes []ps.Process) map[int]ps.Process {
	results := make(map[int]ps.Process)
	for _, arg := range args {
		if strings.HasPrefix(arg, ":") {
			pids := handlePort(arg[1:], processes)
			for _, pid := range pids {
				if proc, found := findProcessByPID(pid, processes); found {
					results[pid] = proc
				}
			}
		} else if stat, err := os.Stat(arg); err == nil {
			if !stat.IsDir() {
				arg, err = filepath.Abs(arg)
				if err != nil {
					fmt.Println("Error getting absolute path:", err)
					continue
				}
			}
			pid, err := handlePath(arg, processes)
			if err == nil && pid != 0 {
				if proc, found := findProcessByPID(pid, processes); found {
					results[pid] = proc
				}
			}
		} else if pid, err := strconv.Atoi(arg); err == nil {
			if proc, found := findProcessByPID(pid, processes); found {
				results[pid] = proc
			}
		} else {
			pids := handleString(arg, processes)
			for _, pid := range pids {
				if proc, found := findProcessByPID(pid, processes); found {
					results[pid] = proc
				}
			}
		}
	}
	return results
}

func handlePath(path string, processes []ps.Process) (int, error) {
	var pid int
	var err error

	switch runtime.GOOS {
	case "darwin":
		pid, err = getFileLockingPIDMacOS(path, processes)
	case "linux":
		pid, err = getFileLockingPIDLinux(path)
	default:
		return 0, fmt.Errorf("unsupported OS")
	}

	return pid, err
}

func getFileLockingPIDMacOS(path string, processes []ps.Process) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	for _, proc := range processes {
		pid := proc.Pid()
		fdPath := fmt.Sprintf("/proc/%d/fd", pid)
		fds, err := os.ReadDir(fdPath)
		if err != nil {
			continue
		}

		for _, fdInfo := range fds {
			fdLink, err := os.Readlink(filepath.Join(fdPath, fdInfo.Name()))
			if err != nil {
				continue
			}

			if fdLink == path {
				return pid, nil
			}
		}
	}
	return 0, nil
}

func getFileLockingPIDLinux(path string) (int, error) {
	fdDir := "/proc"
	var wg sync.WaitGroup
	var mu sync.Mutex
	var foundPID int

	dirs, err := os.ReadDir(fdDir)
	if err != nil {
		return 0, err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		wg.Add(1)
		go func(dir fs.DirEntry) {
			defer wg.Done()
			pid, err := strconv.Atoi(dir.Name())
			if err != nil {
				return
			}
			fdPath := filepath.Join(fdDir, dir.Name(), "fd")
			fds, err := os.ReadDir(fdPath)
			if err != nil {
				return
			}

			for _, fd := range fds {
				fdLink, err := os.Readlink(filepath.Join(fdPath, fd.Name()))
				if err == nil && fdLink == path {
					mu.Lock()
					foundPID = pid
					mu.Unlock()
					return
				}
			}
		}(dir)
	}

	wg.Wait()
	return foundPID, nil
}

func handlePID(pid int, processes []ps.Process) {
	if proc, found := findProcessByPID(pid, processes); found {
		printProcess(proc)
	} else {
		fmt.Printf("No process found with PID %d\n", pid)
	}
}

func handleString(str string, processes []ps.Process) []int {
	var pids []int
	for _, p := range processes {
		if strings.Contains(strings.ToLower(p.Executable()), strings.ToLower(str)) {
			pids = append(pids, p.Pid())
		}
	}
	return pids
}

func handlePort(port string, processes []ps.Process) []int {
	conns, err := gopsutilnet.Connections("tcp")
	if err != nil {
		fmt.Println("Error fetching connections:", err)
		return nil
	}

	var pids []int
	for _, conn := range conns {
		if strconv.Itoa(int(conn.Laddr.Port)) == port {
			pids = append(pids, int(conn.Pid))
		}
	}
	return pids
}

func findProcessByPID(pid int, processes []ps.Process) (ps.Process, bool) {
	for _, p := range processes {
		if p.Pid() == pid {
			return p, true
		}
	}
	return nil, false
}

func printProcess(p ps.Process) {
	fmt.Printf("[PID:%d] [PPID:%d] [USER:%s] %s\n", p.Pid(), p.PPid(), getUser(p), p.Executable())
}

func getUser(p ps.Process) string {
	procStat := fmt.Sprintf("/proc/%d/status", p.Pid())
	data, err := os.ReadFile(procStat)
	if err != nil {
		return "unknown"
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			uid, err := strconv.Atoi(fields[1])
			if err == nil {
				user, err := user.LookupId(strconv.Itoa(uid))
				if err == nil {
					return user.Username
				}
			}
		}
	}
	return "unknown"
}
