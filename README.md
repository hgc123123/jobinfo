# Install

JobInfo relies golang language, so you need to install golang firstly.

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

Install JobInfo
```
make
```

# Enable service

Copy jobinfo file to /usr/bin/, and then create file jobinfo.service file under the directory
`/usr/lib/systemd/system/`

The content is as follows:
```
[Unit]
Description=Jobinfo

[Service]
ExecStart=/usr/bin/jobinfo
Restart=always
RestartSec=15

[Install]
WantedBy=multi-user.target
```

```
systemctl enable jobinfo
systemctl start jobinfo
```
