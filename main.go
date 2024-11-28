package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	fmt.Println("Press Enter to quit the program...")
	fmt.Println("Do you want to continue monitoring? (y/n)")
	var choice string
	for i := 0; i < 10; i++ {
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		choice = strings.ToLower(strings.TrimSpace(choice))
		if choice == "y" {
			break
		} else if choice == "n" {
			fmt.Println("Program terminated.")
			return
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}

	fmt.Print("\033[2J\033[H")
	done := make(chan bool)

	go displayStats(done)

	fmt.Scanln()

	done <- true

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Program terminated.")
}

func displayStats(done chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	moveUp := func(n int) {
		fmt.Printf("\033[%dA", n)
	}

	clearDown := "\033[J"

	fmt.Println(strings.Repeat("\n", 100))

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			moveUp(100)

			mem, err := mem.VirtualMemory()
			if err == nil {
				used := mem.Used
				total := mem.Total

				usedGB := float64(used) / (1024 * 1024 * 1024)
				totalGB := float64(total) / (1024 * 1024 * 1024)

				fmt.Printf("Memory Usage: %.2f%% (%.2f GB / %.2f GB)\n", mem.UsedPercent, usedGB, totalGB)
				displayMemoryBar(mem.UsedPercent)
			}

			cpuPercent, err := cpu.Percent(0, false)
			if err == nil && len(cpuPercent) > 0 {
				fmt.Printf("Overall CPU Usage: %.2f%%\n", cpuPercent[0])
				displayUsageBar(cpuPercent[0])
			}

			corePercent, err := cpu.Percent(0, true)
			if err == nil {
				for i := 0; i < len(corePercent); i += 2 {
					if i+1 < len(corePercent) {
						fmt.Printf("CPU Core %d: %.2f%%\t\tCPU Core %d: %.2f%%\n", i, corePercent[i], i+1, corePercent[i+1])
						displayDualUsageBar(corePercent[i], corePercent[i+1])
					} else {
						fmt.Printf("CPU Core %d: %.2f%%\n", i, corePercent[i])
						displayUsageBar(corePercent[i])
					}
				}
			}

			fmt.Print(clearDown)
		}
	}
}

func displayUsageBar(usage float64) {
	const barLength = 50
	filledLength := int(usage / 100 * barLength)

	bar := make([]rune, barLength)
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar[i] = '#'
		} else {
			bar[i] = '-'
		}
	}

	fmt.Printf("[%s]\n\n", string(bar))
}

func displayDualUsageBar(usage1, usage2 float64) {
	const barLength = 25
	filledLength1 := int(usage1 / 100 * barLength)
	filledLength2 := int(usage2 / 100 * barLength)

	bar1 := make([]rune, barLength)
	bar2 := make([]rune, barLength)

	for i := 0; i < barLength; i++ {
		if i < filledLength1 {
			bar1[i] = '█'
		} else {
			bar1[i] = '░'
		}
		if i < filledLength2 {
			bar2[i] = '█'
		} else {
			bar2[i] = '░'
		}
	}

	fmt.Printf("[%s]  [%s]\n\n", string(bar1), string(bar2))
}

func displayMemoryBar(usage float64) {
	const barLength = 50
	filledLength := int(usage / 100 * barLength)

	bar := make([]rune, barLength)
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar[i] = '█'
		} else {
			bar[i] = '░'
		}
	}

	fmt.Printf("[%s]\n\n", string(bar))
}
