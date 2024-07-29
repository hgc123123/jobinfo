TARGET:=jobinfo
SOURCES:=main.go metrics/cgroups.go metrics/cpuset.go metrics/cpuacct.go metrics/memory.go collectors/slurm.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
