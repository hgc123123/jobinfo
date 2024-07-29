package collectors

import (
	"regexp"
	"strconv"
	"strings"
	//"fmt"
	ps "github.com/mitchellh/go-ps"
	cg "github.com/hgc123123/jobinfo/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/hgc123123/jobinfo/script"
	log "github.com/sirupsen/logrus"
)

type cgroupsSlurmCollector struct {
	cpuacctUsagePerCPUMetric *prometheus.Desc
	memoryUsageInBytesMetric *prometheus.Desc
	cpusetCPUsMetric         *prometheus.Desc
	gpuUsageMetric           *prometheus.Desc
	vramUsageMetric          *prometheus.Desc
	cgroupsRootPath          string
        desc     *prometheus.Desc
	metric   *prometheus.GaugeVec
	labelVal []string
}

func NewCgroupsSlurmCollector(cgroupsRootPath string) *cgroupsSlurmCollector {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "running_job_script",
			Help: "This is running job script",
		},
		[]string{"user_id", "job_id", "step_id", "task_id", "sbatch_jobname"}, // 这里定义了三个标签
	)

	// 创建一个描述符
	desc := prometheus.NewDesc(
		"running_job_script_desc",
		"This is running job script descripter",
		[]string{"user_id", "job_id", "step_id", "task_id", "sbatch_jobname"}, // 这里定义了与指标相关的三个标签
		nil,
	)
	return &cgroupsSlurmCollector{
		cpuacctUsagePerCPUMetric: prometheus.NewDesc("usage_cpu_each",
			"Per-nanosecond usage of each CPU in a cgroup",
			[]string{"user_id", "job_id", "step_id", "task_id", "cpu_id"}, nil,
		),
		memoryUsageInBytesMetric: prometheus.NewDesc("usage_of_memory",
			"Current memory used by the cgroup in bytes",
			[]string{"user_id", "job_id", "step_id", "task_id"}, nil,
		),
		cpusetCPUsMetric: prometheus.NewDesc("cpuset_all_cpus",
			"List of CPUs and whether or not they are in the cpuset cgroup",
			[]string{"user_id", "job_id", "step_id", "task_id", "cpu_id"}, nil,
		),
		gpuUsageMetric: prometheus.NewDesc("usage_gpu_each",
                        "Usage of each GPU in a cgroup",
                        []string{"user_id", "job_id", "step_id", "task_id", "gpu_id"}, nil,
		),
                vramUsageMetric: prometheus.NewDesc("vram_usage_gpu_each",
                        "Usage of each GPU memory in a cgroup",
                        []string{"user_id", "job_id", "step_id", "task_id", "gpu_id"}, nil,
                ),
		cgroupsRootPath: cgroupsRootPath,
		desc: desc,
		metric: metric,
	}
}

func (collector *cgroupsSlurmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.cpuacctUsagePerCPUMetric
	ch <- collector.memoryUsageInBytesMetric
	ch <- collector.cpusetCPUsMetric
	ch <- collector.gpuUsageMetric
	ch <- collector.vramUsageMetric
	ch <- collector.desc

}

func (collector *cgroupsSlurmCollector) Collect(ch chan<- prometheus.Metric) {
	// Get a list of all processes
	procs, err := ps.Processes()
	if err != nil {
		log.Fatalf("unable to read process table: %v", err)
	}
	var slurmstepdIds []int
	for _, proc := range procs {
		if proc.Executable() == "slurmstepd" {
			slurmstepdIds = append(slurmstepdIds, proc.Pid())
		}
	}
	//collector.labelVal="#SBATCH --jobname=huhu"
	//collector.metric.WithLabelValues(collector.labelVal).Set(123456)
	//collector.metric.WithLabelValues(collector.labelVal).Collect(ch)
	// 发送指标
	//ch <- collector.metric
	// Filter processes by children of slurmstepd processes
	for _, ssid := range slurmstepdIds {
		for _, proc := range procs {
			if proc.PPid() == ssid {
				cgroups, err := cg.LoadProcessCgroups(proc.Pid(), collector.cgroupsRootPath)
				if err != nil {
					log.Fatalf("unable to read cgroups file: %v", err)
				}
				slurmRegex := regexp.MustCompile(`/slurm(?:/uid_([^/]+))?(?:/job_([^/]+))?(?:/step_([^/]+))?(?:/task_([^/]+))?`)
				matches := slurmRegex.FindStringSubmatch(string(cgroups.Cpuacct))
				/*
				exec_command := strings.Split(strings.Split(matches[0],"/")[3],"_")[1]
				final_command := "/var/spool/slurmd/job"+exec_command+"/slurm_script"
				command_exec_content,err := script.GetContentOfScript(final_command)
				collector.contentOfBatchScript=command_exec_content
				*/
				if err != nil{
					log.Fatalf("unable to read cpuacct usage per cpu: %v", err)
				}
				var (
					user_id string
					job_id  string
					step_id string
					task_id string
				)
				if len(matches) > 1 {
					user_id = matches[1]
				}
				if len(matches) > 2 {
					job_id = matches[2]
				}
				if len(matches) > 3 {
					step_id = matches[3]
				}
				if len(matches) > 4 {
					task_id = matches[4]
				}
				// cpuacctUsagePerCPUMetric
				usagePerCPU, err := cgroups.Cpuacct.GetUsagePerCPU()
				if err != nil {
					log.Fatalf("unable to read cpuacct usage per cpu: %v", err)
				}
				for cpuID, cpuUsage := range usagePerCPU {
					ch <- prometheus.MustNewConstMetric(collector.cpuacctUsagePerCPUMetric,
						prometheus.GaugeValue, float64(cpuUsage), user_id, job_id, step_id, task_id, strconv.Itoa(cpuID))                      
				}
				// usagePerGPU
                                usagePerGPU, err := cgroups.Devices.GetUsagePerGPU()
                                if err != nil {
                                        log.Fatalf("unable to read usage per gpu: %v", err)
                                }
				for gpuNumber, gpuUtil := range usagePerGPU {
                                	//fmt.Printf("usage of %d is %d\n",gpuNumber,gpuUtil)
					ch <- prometheus.MustNewConstMetric(collector.gpuUsageMetric,
                                                prometheus.GaugeValue, float64(gpuUtil), user_id, job_id, step_id, task_id, strconv.Itoa(gpuNumber))
				}
				// vramUsagePerGPU
                                vramUsagePerGPU, err := cgroups.Devices.GetVRAMUsagePerGPU()
                                if err != nil {
                                        log.Fatalf("unable to read usage per gpu: %v", err)
                                }
                                for gpuNumber, vramUtil := range vramUsagePerGPU {
                                        //fmt.Printf("usage of %d is %d\n",gpuNumber,gpuUtil)
                                        ch <- prometheus.MustNewConstMetric(collector.vramUsageMetric,
                                                prometheus.GaugeValue, float64(vramUtil), user_id, job_id, step_id, task_id, strconv.Itoa(gpuNumber))
                                }
				// memoryUsageInBytesMetric
				memoryUsageBytes, err := cgroups.Memory.GetUsageInBytes()
				if err != nil {
					log.Fatalf("unable to read memory usage in bytes: %v", err)
				}
				ch <- prometheus.MustNewConstMetric(collector.memoryUsageInBytesMetric,
					prometheus.GaugeValue, float64(memoryUsageBytes), user_id, job_id, step_id, task_id)
				// cpusetCPUsMetric
				cpusetCPUs, err := cgroups.Cpuset.GetCpus()
				if err != nil {
					log.Fatalf("unable to read cpuset CPUs: %v", err)
				}
				for _, cpuID := range cpusetCPUs {
					ch <- prometheus.MustNewConstMetric(collector.cpusetCPUsMetric,
						prometheus.GaugeValue, float64(1), user_id, job_id, step_id, task_id, strconv.Itoa(cpuID))
				}
				exec_command := strings.Split(strings.Split(matches[0],"/")[3],"_")[1]
				final_command := "/var/spool/slurmd/job"+exec_command+"/slurm_script"
				command_exec_content,err := script.GetContentOfScript(final_command)


				collector.labelVal=[]string{user_id, job_id, step_id, task_id, command_exec_content}
	        		metric := collector.metric.WithLabelValues(collector.labelVal[0],collector.labelVal[1],collector.labelVal[2],collector.labelVal[3],collector.labelVal[4])
				metric.Set(123)
				ch <- metric
			}
		}
	}
}
