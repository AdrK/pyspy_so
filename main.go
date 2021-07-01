package main

import (
	"C"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/AdrK/pyspy_so/pkg/agent"
	"github.com/AdrK/pyspy_so/pkg/agent/pyspy"
	"github.com/AdrK/pyspy_so/pkg/agent/spy"
	"github.com/AdrK/pyspy_so/pkg/agent/types"
	"github.com/AdrK/pyspy_so/pkg/agent/upstream/remote"
	"github.com/AdrK/pyspy_so/pkg/config"
	"github.com/AdrK/pyspy_so/pkg/util/names"
)
import "strconv"

func processExists(pid int) bool {
	check := nil == syscall.Kill(pid, 0)
	logrus.Infof("checking if process exists: ", check)
	return check
}

func startNewSession(cfg *config.Exec) error {
	if !processExists(cfg.Pid) {
		return errors.New("process not found")
	}

	pyspy.Blocking = cfg.PyspyBlocking

	spyName := cfg.SpyName
	if spyName == "auto" {
		return fmt.Errorf("not supported")
	}

	logrus.Info("to disable logging from pyroscope, pass " + color.YellowString("-no-logging") + " argument to pyroscope exec")

	/*
		if err := performChecks(spyName); err != nil {
			return err
		}
	*/

	if cfg.ApplicationName == "" {
		logrus.Infof("we recommend specifying application name via %s flag or env variable %s",
			color.YellowString("-application-name"), color.YellowString("PYROSCOPE_APPLICATION_NAME"))
		cfg.ApplicationName = spyName + "." + names.GetRandomName(generateSeed())
		logrus.Infof("for now we chose the name for you and it's \"%s\"", color.GreenString(cfg.ApplicationName))
	}

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
	pid := cfg.Pid

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

func waitForSpawnedProcessToExit(c chan os.Signal, cmd *exec.Cmd) error {
	go func() {
		for s := range c {
			_ = cmd.Process.Signal(s)
		}
	}()
	return cmd.Wait()
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

func generateSeed() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "<unknown>"
	}
	return cwd + "|" + "&"
}

//export Start
func Start(ApplicationName *C.char, Pid C.int) {
	logrus.SetLevel(logrus.DebugLevel)

	conf := config.Exec{
		SpyName:                "pyspy",
		ApplicationName:        C.GoString(ApplicationName),
		SampleRate:             100,
		DetectSubprocesses:     true,
		LogLevel:               "debug",
		ServerAddress:          "http://192.168.5.16:4040",
		AuthToken:              "",
		UpstreamThreads:        4,
		UpstreamRequestTimeout: time.Second * 10,
		NoLogging:              false,
		NoRootDrop:             false,
		Pid:                    int(Pid),
		UserName:               "",
		GroupName:              "",
		PyspyBlocking:          false,
	}

	startNewSession(&conf)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	fmt.Println("app name: '", os.Args[1], "' pid: ", os.Args[2])
	pid, _ := strconv.Atoi(os.Args[2])
	Start(C.CString(os.Args[1]), C.int(pid))
}
