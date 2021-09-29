SET PROJECT="onefinger"
SET BUILDFILE="../main.go"

: delete other binary
del /q %PROJECT%*

::build windows 
go build -ldflags "-w -s" -o %PROJECT%.exe %BUILDFILE%

::build linux amd64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-w -s" -o %PROJECT%-amd64  %BUILDFILE%

::build linux 386
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -ldflags "-s -w" -o %PROJECT%-386  %BUILDFILE%

::build linux arm
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -ldflags "-s -w" -o %PROJECT%-arm  %BUILDFILE%