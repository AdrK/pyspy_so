module github.com/AdrK/pyspy_so

go 1.14

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59
	github.com/blang/semver v3.5.1+incompatible
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/clarkduvall/hyperloglog v0.0.0-20171127014514-a0107a5d8004
	github.com/cosmtrek/air v1.27.3
	github.com/creack/pty v1.1.13 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/dgrijalva/lfu-go v0.0.0-20141010002404-f174e76c5138
	github.com/fatih/color v1.12.0
	github.com/felixge/fgprof v0.9.1
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/pprof v0.0.0-20210609004039-a478d1d731e9
	github.com/google/uuid v1.2.0
	github.com/iancoleman/strcase v0.1.3
	github.com/ianlancetaylor/demangle v0.0.0-20210406231658-61c622dd7d50 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josephspurrier/goversioninfo v1.2.0
	github.com/kardianos/service v1.2.0
	github.com/kisielk/godepgraph v0.0.0-20190626013829-57a7e4a651a9
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kyoh86/richgo v0.3.9
	github.com/kyoh86/xdg v1.2.0 // indirect
	github.com/markbates/pkger v0.17.1
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/goreman v0.3.7
	github.com/mgechev/revive v1.0.8
	github.com/mitchellh/go-ps v1.0.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/peterbourgon/ff/v3 v3.0.0
	github.com/prometheus/client_golang v1.11.0
	github.com/pyroscope-io/dotnetdiag v1.2.1
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/twmb/murmur3 v1.1.5
	github.com/wacul/ptr v1.0.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22
	golang.org/x/tools v0.1.4
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.2.0
)

replace github.com/mgechev/revive v1.0.3 => github.com/pyroscope-io/revive v1.0.6-0.20210330033039-4a71146f9dc1

replace github.com/dgrijalva/lfu-go => github.com/pyroscope-io/lfu-go v1.0.3
