# Install Goalng

 [Official document](https://go.dev/doc/install) 
 
```bash
$ wget 'https://go.dev/dl/go1.18.3.linux-amd64.tar.gz'

$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz

$ export PATH=$PATH:/usr/local/go/bin

$ go version
```

## Install dependecy package

```bash
$ go mod tidy
```

## Modify main.go.

```bash
$ nano main.go
```

```go
  // Find it and press your broker connection information.
	brokerHost := ""
	brokerPassword := ""
	brokerUsername := ""
```

## Compile

```bash
go build main.go -o SmartHost
```

## Execute

In Linux check your pzem serial port with 

```
ls /dev/ttyU*
```

In windwos check your pzem com port with device manager

```
COM4
```

Execute
```
./SmartHost [PZEM serial port] [Your Host Name]

./SmartHost /dev/ttyUSB0 myServer
```
