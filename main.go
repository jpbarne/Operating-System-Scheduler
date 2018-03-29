package main

import (
	"fmt"
	"log"
	"os"
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
var work_queue = make([]Process, MAX)
var on_cpu Process

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
  for i := 0; i < num_of_processes; i++ {
    fmt.Println("Process ID: ", simulation_load[i].proccess_id)
    fmt.Println("   Arrival Time: ", simulation_load[i].arrival_time)
    fmt.Println("   Process Length: ", simulation_load[i].proccess_length)
    fmt.Println("   Completion Time: ", simulation_load[i].completion_time)
    fmt.Println("   Response Time: ", simulation_load[i].response_time)
    avg += simulation_load[i].response_time
  }

  avg = avg / num_of_processes
  fmt.Println(" Avg Response Time: ", avg)
  fmt.Println(" Number of Context Switches: ", switches)

}

// First in first out
func FIFO()  {
	for i := 0; i < num_of_processes; i++ {
		simulation_load[i].completion_time = simulation_load[i].proccess_length + master_clock
		master_clock += simulation_load[i].proccess_length
		simulation_load[i].response_time = simulation_load[i].completion_time - simulation_load[i].arrival_time
	}
	switches = num_of_processes
}

func SJF()  {

}

func RR()  {
	done := false

	for !done {

		if len(simulation_load) > 0 {
			for i := 0; i < num_of_processes; i++ {
				if simulation_load[i].arrival_time > master_clock ||
				simulation_load[i].arrival_time <= master_clock+time_quantum {
				work_queue = append(work_queue, simulation_load[i])
				simulation_load = append(simulation_load[:i], simulation_load[i+1:]...)
				}
			}
		}

		switches++
		on_cpu = work_queue[0]

	 //not going to finish in time_quantum
		if on_cpu.time_remaining > time_quantum {
			on_cpu.time_remaining -= time_quantum
			master_clock += time_quantum
			work_queue = append(work_queue, on_cpu)
		} else if on_cpu.time_remaining <= time_quantum {
			master_clock += on_cpu.time_remaining
			on_cpu.completion_time = master_clock
			on_cpu.time_remaining = 0
			on_cpu.response_time = on_cpu.completion_time - on_cpu.arrival_time
		}

		if len(work_queue) == 0 {
			done = true
		}
	}
}


func main() {
  read_data()

	if scheduling_policy == 0 {
		FIFO()
  } else if scheduling_policy == 1 {
		SJF()
  } else if scheduling_policy == 2 {
		RR()
  }

	print_report()
}
