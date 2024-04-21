TARGET:=jobinfo
SOURCES:=main.go cgroups/cgroups.go cgroups/cpuset.go cgroups/cpuacct.go cgroups/memory.go collectors/slurm.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
