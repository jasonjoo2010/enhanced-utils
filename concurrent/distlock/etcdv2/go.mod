module github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/redis

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/jasonjoo2010/enhanced-utils v0.0.0-20200608071141-0f10b99c6fe4
	github.com/onsi/ginkgo v1.12.3 // indirect
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
