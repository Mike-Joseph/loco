# NOT MAINTAINED

Note that this is a project that was developed internally at The Mode Group
and has now been released under the Apache-2.0 license (see LICENSE and
NOTICE).  The Mode Group is not actively maintaining this project and is not
currently accepting pull requests, bug reports, or changes.  Users of this
project are welcome to fork it under the allowable terms of its license and
continue the project at their own discretion.

# loco
LoCo - Loop Control

## Building

Be sure to have `Go 1.9` or later

Build the binary:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o loco cmd/main.go
```

Then do the docker
```
docker build -t modehq/loco .
```
