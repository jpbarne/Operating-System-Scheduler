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
	//for each task update values
	for i := 0; i < num_of_processes; i++ {
		simulation_load[i].completion_time = simulation_load[i].proccess_length + master_clock
		master_clock += simulation_load[i].proccess_length
		simulation_load[i].response_time = simulation_load[i].completion_time - simulation_load[i].arrival_time
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

//Shortest-Job-First
func SJF() {
	//work_load holds all jobs to be completed
	work_load := make([]Process, 0)
	work_load = append(work_load, simulation_load[0])
	simulation_load = simulation_load[1:]

	for {
		//update on_cpu, switches, and work_load
		sort.Sort(Shortest(work_load)) //sort by time_remaining
		on_cpu := work_load[0]
		work_load = work_load[1:]
		switches++

		//update clock and time remaining
		if on_cpu.time_remaining >= time_quantum {
			on_cpu.time_remaining -= time_quantum
			master_clock += time_quantum
		} else {
			master_clock += on_cpu.time_remaining
			on_cpu.time_remaining = 0
		}

		//check if job complete
		if on_cpu.time_remaining == 0 {
			//update task variables
			on_cpu.completion_time = master_clock
			on_cpu.response_time = on_cpu.completion_time - on_cpu.arrival_time
			simulation_load = append(simulation_load, on_cpu)
		} else { //if not complete, add task back to work load
			work_load = append(work_load, on_cpu)
		}

		//fill work load list and remove work load elements from simulation load
		new_load := make([]Process, 0)
		for _, proc := range simulation_load {
			if proc.arrival_time <= master_clock &&
				 proc.arrival_time > master_clock - time_quantum {
				 work_load = append(work_load, proc)
			 } else {
				 new_load = append(new_load, proc)
			 }
		}
		simulation_load = new_load

		//continue while jobs remain
		if len(work_load) == 0 {
			break
		}
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

	//chose task based on scheduling_policy
	if scheduling_policy == 0 {
		FIFO()
		sort.Sort(Processes(simulation_load))
  } else if scheduling_policy == 1 {
		SJF()
		sort.Sort(Processes(simulation_load))
  } else if scheduling_policy == 2 {
		RR()
		sort.Sort(Processes(simulation_load))
  }

	print_report()
}
