package main

import (
	"C"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/AdrK/pyspy_so/pkg/agent"
	"github.com/AdrK/pyspy_so/pkg/agent/pyspy"
	"github.com/AdrK/pyspy_so/pkg/agent/spy"
	"github.com/AdrK/pyspy_so/pkg/agent/types"
	"github.com/AdrK/pyspy_so/pkg/agent/upstream/remote"
	"github.com/AdrK/pyspy_so/pkg/config"
)

func processExists(pid int) bool {
	exists := nil == syscall.Kill(pid, 0)
	return exists
}

func startNewSession(cfg *config.Exec) error {
	// TODO: Removed some checks to simplify. Bring them back at the end.

	spyName := cfg.SpyName
	pid := cfg.Pid
	pyspy.Blocking = cfg.PyspyBlocking

	rc := remote.RemoteConfig{
		AuthToken:              cfg.AuthToken,
		UpstreamAddress:        cfg.ServerAddress,
		UpstreamThreads:        cfg.UpstreamThreads,
		UpstreamRequestTimeout: cfg.UpstreamRequestTimeout,
	}
	u, err := remote.New(rc, logrus.StandardLogger())
	if err != nil {
		return fmt.Errorf("new remote upstream: %v", err)
	}
	defer u.Stop()

	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(c)
		close(c)
	}()

	logrus.WithFields(logrus.Fields{
		"app-name":            cfg.ApplicationName,
		"spy-name":            spyName,
		"pid":                 pid,
		"detect-subprocesses": cfg.DetectSubprocesses,
	}).Debug("starting agent session")

	if cfg.SampleRate == 0 {
		cfg.SampleRate = types.DefaultSampleRate
	}

	sc := agent.SessionConfig{
		Upstream:         u,
		AppName:          cfg.ApplicationName,
		ProfilingTypes:   []spy.ProfileType{spy.ProfileCPU},
		SpyName:          spyName,
		SampleRate:       uint32(cfg.SampleRate),
		UploadRate:       10 * time.Second,
		Pid:              pid,
		WithSubprocesses: cfg.DetectSubprocesses,
	}
	session := agent.NewSession(&sc, logrus.StandardLogger())
	if err = session.Start(); err != nil {
		return fmt.Errorf("start session: %v", err)
	}
	defer session.Stop()

	waitForProcessToExit(c, pid)
	return nil
}

func waitForProcessToExit(c chan os.Signal, pid int) {
	if pid == -1 {
		<-c
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c:
			return
		case <-ticker.C:
			if !processExists(pid) {
				logrus.Debug("child process exited")
				return
			}
		}
	}
}

//export Start
func Start(ApplicationName *C.char, Pid C.int, SpyName *C.char, ServerAddress *C.char) {
	logrus.SetLevel(logrus.DebugLevel)

	// TODO: It might be more useful if it would be []pids instead of pid

	startNewSession(&config.Exec{
		SpyName:                C.GoString(SpyName),
		ApplicationName:        C.GoString(ApplicationName),
		SampleRate:             100,
		DetectSubprocesses:     true,
		LogLevel:               "debug",
		ServerAddress:          C.GoString(ServerAddress),
		AuthToken:              "",
		UpstreamThreads:        4,
		UpstreamRequestTimeout: time.Second * 10,
		NoLogging:              false,
		NoRootDrop:             false,
		Pid:                    int(Pid),
		UserName:               "",
		GroupName:              "",
		PyspyBlocking:          false,
	})
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	fmt.Println("app name:", os.Args[1], "pid: ", os.Args[2], "spy name: ", os.Args[3], "server address: ", os.Args[4])
	pid, _ := strconv.Atoi(os.Args[2])
	Start(C.CString(os.Args[1]), C.int(pid), C.CString(os.Args[3]), C.CString(os.Args[4]))
}
