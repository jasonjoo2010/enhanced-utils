enhanced-utils
===

Useful utilities collection.

# Concurrent

## Distributed lock
It provides interface definition and memory(mock) and redis implementations.

```go
mutex_lock := distlock.NewMutex("project-namespace", 60*time.Second, redis.New([]string{"127.0.0.1:6379"}))

reentry_lock := distlock.NewReentry("project-namespace", 60*time.Second, redis.New([]string{"127.0.0.1:6379"}))
```

### Storage Supported for Lock

* Mock(memory)
* Redis
* Etcdv2
* Etcdv3
* Zookeeper
* Database

# String

## Convertion
```go
strutils.ToUnderscore("AppleWatch") -> "apple_watch"
strutils.ToCamel("apple_watch") -> "appleWatch"
```

## Random
```go
strutils.RandNumbers(10) -> "0001234567"
strutils.RandString(4) -> "aZb0"
strutils.RandLowCased(4) -> "aaz0"
strutils.RandHash(4) -> "03a0"
strutils.RandPrintable(4) -> "03~0"
```

## Validation
```go
strutils.IsURL("http://www.google.com") -> true
strutils.IsEmail("jack@google.com") -> true
```
