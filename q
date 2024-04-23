[0;1;32m●[0m jobinfo.service - Prometheus Jobinfo
   Loaded: loaded (/usr/lib/systemd/system/jobinfo.service; disabled; vendor preset: disabled)
   Active: [0;1;32mactive (running)[0m since Sun 2024-04-21 11:09:16 CST; 3s ago
 Main PID: 3171868 (jobinfo)
    Tasks: 6 (limit: 3297702)
   Memory: 7.2M
   CGroup: /system.slice/jobinfo.service
           └─3171868 /usr/bin/jobinfo

Apr 21 11:09:16 compute133 systemd[1]: Started Prometheus Jobinfo.
Apr 21 11:09:16 compute133 jobinfo[3171868]: time="2024-04-21T11:09:16+08:00" level=info msg="serving cgroups from hierarchy root /sys/fs/cgroup"
Apr 21 11:09:16 compute133 jobinfo[3171868]: time="2024-04-21T11:09:16+08:00" level=info msg="listening on port 9821"
