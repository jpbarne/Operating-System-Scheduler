package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

const MAX int = 1000

//data for each task
type Process struct{
  proccess_id int
  arrival_time int
  proccess_length int
  process_loaded int
  time_remaining int
  time_waiting int
  completion_time int
  response_time int
}

//global variables

var simulation_load = make([]Process, MAX)

var num_of_processes int
var scheduling_policy int
var preemption_policy int
var time_quantum int
var master_clock int
var processes_left int
var switches int

//scan in data
func read_data() {
	file, err := os.Open("simulation_load.txt");
  if err != nil {
    log.Fatal(err)
  }
	var temp string

  fmt.Fscanf(file, "%d %s %s\n", &scheduling_policy, &temp, &temp)
  fmt.Fscanf(file, "%d %s\n", &preemption_policy, &temp)
  fmt.Fscanf(file, "%d %s %s\n", &time_quantum, &temp, &temp)
  fmt.Fscanf(file, "%d %s %s %s\n\n", &num_of_processes, &temp, &temp, &temp)

  for i := 0; i < num_of_processes; i++ {
    fmt.Fscanf(file, "%d %s\n", &(simulation_load[i].proccess_id), &temp)
    fmt.Fscanf(file, "%d %s\n", &(simulation_load[i].proccess_length), &temp)
    fmt.Fscanf(file, "%d %s\n\n", &(simulation_load[i].arrival_time), &temp)

    simulation_load[i].process_loaded = 0
    simulation_load[i].time_remaining = simulation_load[i].proccess_length
    simulation_load[i].time_waiting = 0
    simulation_load[i].completion_time = -1
    simulation_load[i].response_time = -1
  }

	simulation_load = simulation_load[:num_of_processes]

  if err := file.Close(); err != nil {
    log.Fatal(err)
  }

  if scheduling_policy == 0 {
    preemption_policy = 0
  }
  if scheduling_policy == 2 {
    preemption_policy = 1
  }
  if preemption_policy < 0 || preemption_policy > 1 {
    preemption_policy = 1
  }

  master_clock = 0
  processes_left = 0
  switches = 0

}

//Print final results
func print_report()  {
  avg := 0

	//print schedule policy used
  if scheduling_policy == 0 {
    fmt.Println("Scheduling Policy: FIFO")
  } else if scheduling_policy == 1 {
    fmt.Println("Scheduling Policy: SJF")
  } else if scheduling_policy == 2 {
    fmt.Println("Scheduling Policy: RR")
  }

	//print preemption policy
  if preemption_policy == 0 {
    fmt.Println("Preemption: OFF")
  } else if preemption_policy == 1 {
    fmt.Println("Preemption: ON")
  }

  fmt.Println("Time Quantum: ", time_quantum)
  fmt.Println("Number of Process: ", num_of_processes)

	//print each tasks results
  for _, proc := range simulation_load {
    fmt.Println("\nProcess ID: ", proc.proccess_id)
    fmt.Println("   Arrival Time: ", proc.arrival_time)
    fmt.Println("   Process Length: ", proc.proccess_length)
    fmt.Println("   Completion Time: ", proc.completion_time)
    fmt.Println("   Response Time: ", proc.response_time)
    avg += proc.response_time
  }

  avg = avg / num_of_processes
  fmt.Println("\n Avg Response Time: ", avg)
  fmt.Println(" Number of Context Switches: ", switches)

}

// First in first out
func FIFO()  {
	for _, proc := range simulation_load {
		proc.completion_time = proc.proccess_length + master_clock
		master_clock += proc.proccess_length
		proc.response_time = proc.completion_time - proc.arrival_time
	}
	switches = num_of_processes
}

// for sorting by shortest job length
type Shortest []Process
func (p Shortest) Len() int {
	return len(p)
}
func (p Shortest) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Shortest) Less(i, j int) bool {
	return p[i].time_remaining < p[j].time_remaining
}

func SJF_PREE(work_load []Process, last_proc Process)  {
	work_load = append(work_load, simulation_load[0])
	simulation_load = simulation_load[1:]

	for {
		//update on_cpu and work_load
		sort.Sort(Shortest(work_load))
		on_cpu := work_load[0]
		work_load = work_load[1:]

		//update switches
		/*if last_proc.proccess_id != on_cpu.proccess_id {
			fmt.Println(last_proc, "		", on_cpu)
			switches++
		}*/

		// for time Quantum
/*		for i := 0; i < time_quantum; i++ {
			//check current job is the shortest
			var shorter Process
			for _, proc := range work_load {
				if proc.time_remaining < on_cpu.time_remaining &&
					 proc.arrival_time == master_clock {
					 shorter = proc
				 }
			}

			// shorter job arrived -> preempt
			if shorter != (Process{}){
				work_load = append(work_load, on_cpu)
				on_cpu = shorter
				switches++
			}

			on_cpu.time_remaining--
			master_clock++
			if on_cpu.time_remaining == 0 {
				on_cpu.completion_time = master_clock
				on_cpu.response_time = on_cpu.completion_time - on_cpu.arrival_time
				simulation_load = append(simulation_load, on_cpu)
				break
			}
		} */

		for i := 0; i < time_quantum; i++ {
			//fill work load list
			new_load := make([]Process, 0)
			for _, proc := range simulation_load {
				if proc.arrival_time == master_clock {
					 work_load = append(work_load, proc)
				 } else {
					 new_load = append(new_load, proc)
				 }
			}
			simulation_load = new_load
			sort.Sort(Shortest(work_load))

			if len(work_load) > 0 {
				if work_load[0].time_remaining < on_cpu.time_remaining &&
					 work_load[0].arrival_time >= master_clock 						 {
					 fmt.Println(work_load[0])
					 i = 0
					 work_load = append(work_load, on_cpu)
					 on_cpu = work_load[0]
					 work_load = work_load[1:]
					 sort.Sort(Shortest(work_load))
					 switches++
				}
			}
			on_cpu.time_remaining--
			master_clock++

			if on_cpu.time_remaining == 0 {
				switches++
				on_cpu.completion_time = master_clock
				on_cpu.response_time = on_cpu.completion_time - on_cpu.arrival_time
				simulation_load = append(simulation_load, on_cpu)
			}
		}
		//task didn't finish
		if on_cpu.time_remaining > 0 {
			work_load = append(work_load, on_cpu)
		}

		if len(work_load) == 0 {
			break
		}
	}
}

func SJF_NON(work_load []Process, last_proc Process)  {

}

//Shortest-Job-First
func SJF() {
	work_load := make([]Process, 0)
	last_proc := Process{-1, 0, 0, 0, 0, 0, 0, 0}

	if preemption_policy == 1 {
		SJF_PREE(work_load, last_proc)
	} else {
		SJF_NON(work_load, last_proc)
	}
}

//Round-robin
func RR()  {
	work_queue := make([]Process, 0)

	// main loop
	for {
		// append processes arriving to work queue
		if len(simulation_load) > 0 {
			new_load := make([]Process, 0)
			for _, proc := range simulation_load {
				if proc.arrival_time >= master_clock &&
					 proc.arrival_time < master_clock+time_quantum {
					work_queue = append(work_queue, proc)
				} else {
					new_load = append(new_load, proc)
				}
			}
			simulation_load = new_load
		}

		switches++
		on_cpu := work_queue[0] 		// load current process
		work_queue = work_queue[1:] // pop current proccess off queue

	 // not going to finish in time_quantum
		if on_cpu.time_remaining > time_quantum {
			on_cpu.time_remaining -= time_quantum
			master_clock += time_quantum
			work_queue = append(work_queue, on_cpu)
		} else if on_cpu.time_remaining <= time_quantum { // will finish
			master_clock += on_cpu.time_remaining
			on_cpu.completion_time = master_clock
			on_cpu.time_remaining = 0
			on_cpu.response_time = on_cpu.completion_time - on_cpu.arrival_time

			simulation_load = append(simulation_load, on_cpu)
		}

		// terminate loop
		if len(work_queue) == 0 {
			break
		}
	}
}

// sorting functions
type Processes []Process
func (p Processes) Len() int {
	return len(p)
}
func (p Processes) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Processes) Less(i, j int) bool {
	return p[i].proccess_id < p[j].proccess_id
}

func main() {
  read_data()

	if scheduling_policy == 0 {
		FIFO()
  } else if scheduling_policy == 1 {
		SJF()
		sort.Sort(Processes(simulation_load))
  } else if scheduling_policy == 2 {
		RR()
		sort.Sort(Processes(simulation_load))
  }

	print_report()
}
