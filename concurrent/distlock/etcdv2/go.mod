module github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/etcdv2

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/jasonjoo2010/enhanced-utils v0.0.2
	github.com/onsi/ginkgo v1.12.3 // indirect
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
