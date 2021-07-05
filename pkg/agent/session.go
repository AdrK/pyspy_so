package agent

import (
	"sync"
	"time"

	// revive:disable:blank-imports Depending on configuration these packages may or may not be used.
	//   That's why we do a blank import here and then packages themselves register with the rest of the code.

	_ "github.com/AdrK/pyspy_so/pkg/agent/pyspy"
	"github.com/AdrK/pyspy_so/pkg/util/slices"
	"github.com/sirupsen/logrus"

	// revive:enable:blank-imports

	"github.com/AdrK/pyspy_so/pkg/agent/spy"
	"github.com/AdrK/pyspy_so/pkg/agent/upstream"
	"github.com/AdrK/pyspy_so/pkg/structs/transporttrie"
	"github.com/mitchellh/go-ps"
)

type ProfileSession struct {
	upstream   upstream.Upstream
	appName    string
	spyName    string
	sampleRate uint32
	uploadRate time.Duration
	pids       []int
	spies      []spy.Spy
	stopCh     chan struct{}
	trieMutex  sync.Mutex

	previousTries []*transporttrie.Trie
	tries         []*transporttrie.Trie

	profileTypes     []spy.ProfileType
	disableGCRuns    bool
	withSubprocesses bool

	startTime time.Time
	stopTime  time.Time

	logger Logger
}

type SessionConfig struct {
	Upstream         upstream.Upstream
	AppName          string
	ProfilingTypes   []spy.ProfileType
	DisableGCRuns    bool
	SpyName          string
	SampleRate       uint32
	UploadRate       time.Duration
	Pid              int
	WithSubprocesses bool
}

func NewSession(c *SessionConfig, logger Logger) *ProfileSession {
	ps := &ProfileSession{
		upstream:         c.Upstream,
		appName:          c.AppName,
		spyName:          c.SpyName,
		profileTypes:     c.ProfilingTypes,
		disableGCRuns:    c.DisableGCRuns,
		sampleRate:       c.SampleRate,
		uploadRate:       c.UploadRate,
		pids:             []int{c.Pid},
		stopCh:           make(chan struct{}),
		withSubprocesses: c.WithSubprocesses,
		logger:           logger,
	}

	ps.previousTries = make([]*transporttrie.Trie, 1)
	ps.tries = make([]*transporttrie.Trie, 1)

	return ps
}

func (ps *ProfileSession) takeSnapshots() {
	ticker := time.NewTicker(time.Second / time.Duration(ps.sampleRate))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			isdueToReset := ps.isDueForReset()
			// reset the profiler for spies every upload rate(10s), and before uploading, it needs to read profile data every sample rate
			if isdueToReset {
				for _, s := range ps.spies {
					if sr, ok := s.(spy.Resettable); ok {
						sr.Reset()
					}
				}
			}

			for _, s := range ps.spies {
				s.Snapshot(func(stack []byte, v uint64, err error) {
					if err != nil {
						// TODO: figure out what to do with these messages. A couple of considerations:
						// * We probably shouldn't just suppress these messages as they might be useful for users
						// * We probably want to throttle the messages because this is code that runs 100 times per second.
						//   If we don't throttle we risk upsetting users with a flood of messages
						// * In gospy case we need to add ability for users to bring their own logger, we can't just use logrus here
						return
					}
					if len(stack) > 0 {
						ps.trieMutex.Lock()
						defer ps.trieMutex.Unlock()
						ps.tries[0].Insert(stack, v, true)

					}
				})
			}

			// upload the read data to server and reset the start time
			if isdueToReset {
				ps.reset()
			}

		case <-ps.stopCh:
			// stop the spies
			for _, spy := range ps.spies {
				spy.Stop()
			}
			return
		}
	}
}

func (ps *ProfileSession) Start() error {
	ps.reset()

	s, err := spy.SpyFromName(ps.spyName, ps.pids[0])
	if err != nil {
		return err
	}

	ps.spies = append(ps.spies, s)
	go ps.takeSnapshots()

	return nil
}

func (ps *ProfileSession) isDueForReset() bool {
	// TODO: duration should be either taken from config or ideally passed from server
	now := time.Now().Truncate(ps.uploadRate)
	start := ps.startTime.Truncate(ps.uploadRate)

	return !start.Equal(now)
}

// the difference between stop and reset is that reset stops current session
// and then instantly starts a new one
func (ps *ProfileSession) reset() {
	ps.trieMutex.Lock()
	defer ps.trieMutex.Unlock()

	now := time.Now()
	// upload the read data to server
	ps.uploadTries(now)

	// reset the start time
	ps.startTime = now

	if ps.withSubprocesses {
		ps.addSubprocesses()
	}
}

func (ps *ProfileSession) Stop() {
	logrus.Debug("session finished")
	ps.trieMutex.Lock()
	defer ps.trieMutex.Unlock()

	ps.stopTime = time.Now()
	close(ps.stopCh)
	// TODO: wait for stopCh consumer to finish!

	// before stopping, upload the tries
	ps.uploadTries(time.Now())

	ps.upstream.Stop()
}

// upload the read profile data about 10s to server
func (ps *ProfileSession) uploadTries(now time.Time) {
	for i, trie := range ps.tries {
		skipUpload := false

		if trie != nil {
			endTime := now.Truncate(ps.uploadRate)

			uploadTrie := trie
			if ps.profileTypes[i].IsCumulative() {
				previousTrie := ps.previousTries[i]
				if previousTrie == nil {
					skipUpload = true
				} else {
					// TODO: Diff doesn't remove empty branches. We need to add that at some point
					uploadTrie = trie.Diff(previousTrie)
				}
			}

			if !skipUpload {
				name := ps.appName + "." + string(ps.profileTypes[i])
				ps.upstream.Upload(&upstream.UploadJob{
					Name:            name,
					StartTime:       ps.startTime,
					EndTime:         endTime,
					SpyName:         ps.spyName,
					SampleRate:      ps.sampleRate,
					Units:           ps.profileTypes[i].Units(),
					AggregationType: ps.profileTypes[i].AggregationType(),
					Trie:            uploadTrie,
				})
			}
			if ps.profileTypes[i].IsCumulative() {
				ps.previousTries[i] = trie
			}
		}
		ps.tries[i] = transporttrie.New()
	}
}

func (ps *ProfileSession) addSubprocesses() {
	newPids := findAllSubprocesses(ps.pids[0])
	for _, newPid := range newPids {
		if !slices.IntContains(ps.pids, newPid) {
			ps.pids = append(ps.pids, newPid)
			newSpy, err := spy.SpyFromName(ps.spyName, newPid)
			if err != nil {
				if ps.logger != nil {
					ps.logger.Errorf("failed to initialize a spy %d [%s]", newPid, ps.spyName)
				}
			} else {
				if ps.logger != nil {
					ps.logger.Debugf("started spy for subprocess %d [%s]", newPid, ps.spyName)
				}
				ps.spies = append(ps.spies, newSpy)
			}
		}
	}
}

func findAllSubprocesses(pid int) []int {
	res := []int{}

	childrenLookup := map[int][]int{}
	processes, err := ps.Processes()
	if err != nil {
		// TODO: handle
		return res
	}
	for _, p := range processes {
		ppid := p.PPid()
		if _, ok := childrenLookup[ppid]; !ok {
			childrenLookup[ppid] = []int{}
		}
		childrenLookup[ppid] = append(childrenLookup[ppid], p.Pid())
	}

	todo := []int{pid}
	for len(todo) > 0 {
		parentPid := todo[0]
		todo = todo[1:]

		if children, ok := childrenLookup[parentPid]; ok {
			for _, childPid := range children {
				res = append(res, childPid)
				todo = append(todo, childPid)
			}
		}
	}

	return res
}
