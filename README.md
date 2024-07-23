#Install

JobInfo relies golang language, so you need to to install golang.

```
wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz
tar -C /usr/local/ -zxf go1.19.5.linux-amd64.tar.gz 
```

Add following contents in ~/.bashrc
```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
```
