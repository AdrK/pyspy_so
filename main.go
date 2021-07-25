package main

import (
	"C"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/AdrK/pyspy_so/pkg/agent"
	"github.com/AdrK/pyspy_so/pkg/agent/spy"
	"github.com/AdrK/pyspy_so/pkg/agent/types"
	"github.com/AdrK/pyspy_so/pkg/agent/upstream/remote"
	"github.com/AdrK/pyspy_so/pkg/config"
)

func processExists(pid int) bool {
	exists := nil == syscall.Kill(pid, 0)
	return exists
}

type PyspySession struct {
	session *agent.ProfileSession
}

var pyspy_session = PyspySession{}

func (pys PyspySession) startNewSession(cfg *config.Exec) (*agent.ProfileSession, error) {
	// TODO: Removed some checks to simplify. Bring them back at the end.

	spyName := cfg.SpyName
	pid := cfg.Pid
	//pyspy.Blocking = cfg.PyspyBlocking // TODO: Set it somwhere in cfg

	rc := remote.RemoteConfig{
		AuthToken:              cfg.AuthToken,
		UpstreamAddress:        cfg.ServerAddress,
		UpstreamThreads:        cfg.UpstreamThreads,
		UpstreamRequestTimeout: cfg.UpstreamRequestTimeout,
	}
	u, err := remote.New(rc, logrus.StandardLogger())
	if err != nil {
		return nil, fmt.Errorf("new remote upstream: %v", err)
	}
	//defer u.Stop()

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
		return nil, fmt.Errorf("start session: %v", err)
	}
	//defer session.Stop()

	return session, nil
}

//export Pyroscope_start
func Pyroscope_start(ApplicationName *C.char, Pid C.int, SpyName *C.char, ServerAddress *C.char) int {
	logrus.SetLevel(logrus.DebugLevel)

	// TODO: It might be more useful if it would be []pids instead of pid

	s, err := pyspy_session.startNewSession(&config.Exec{
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

	if err != nil {
		fmt.Println(err.Error())
		return -1
	}

	pyspy_session.session = s

	return 0
}

//export Pyroscope_stop
func Pyroscope_stop(Pid C.int) {
	pyspy_session.session.Stop()
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	fmt.Println("app name:", os.Args[1], "pid: ", os.Args[2], "spy name: ", os.Args[3], "server address: ", os.Args[4])
	pid, _ := strconv.Atoi(os.Args[2])
	Pyroscope_start(C.CString(os.Args[1]), C.int(pid), C.CString(os.Args[3]), C.CString(os.Args[4]))
	time.Sleep(11 * time.Second)
	Pyroscope_stop(C.int(pid))
}
