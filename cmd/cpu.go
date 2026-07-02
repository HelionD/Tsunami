/*
Copyright © 2026 NAME HERE tsunami.project.dev@gmail.com
*/
package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ── cpu command ───────────────────────────────────────────────────────────────
//
// This is the main "tsunami cpu" command.
// You can run it like:
//
//	./tsunami cpu                        → stress all cores at 100% for 10s
//	./tsunami cpu -l 50 -d 30            → stress all cores at 50% for 30s
//	./tsunami cpu -c 2 -l 75 -d 20       → stress only 2 cores at 75% for 20s
//	./tsunami cpu cores                  → just print how many cores you have
//
// How the load percentage works:
//   We break time into tiny 100ms windows.
//   Inside each window the goroutine spins (burns CPU) for `load` ms
//   and then sleeps for the remaining `100 - load` ms.
//
//   Example with --load 30:
//     spin  30ms  <-- eating CPU
//     sleep 70ms  <-- doing nothing
//   That gives roughly 30% CPU usage on that core.
//
//   Example with --load 100:
//     spin  100ms <-- eating CPU the whole time
//     sleep   0ms <-- no rest at all
//   That gives 100% CPU usage on that core.

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Stress the CPU at a given percentage",
	Long: `Stress one or more CPU cores at a chosen load percentage.

Flags:
  -l, --load      How hard to hit each core, from 1 to 100  (default 100)
  -d, --duration  How many seconds to run the stress test   (default 10)
  -c, --cores     How many cores to stress                  (default: all)

Subcommands:
  cores           Print how many CPU cores this machine has

Examples:
  ./tsunami cpu                      stress all cores at 100% for 10s
  ./tsunami cpu -l 50                stress all cores at 50% for 10s
  ./tsunami cpu -c 2 -l 80 -d 30    stress only 2 cores at 80% for 30s
  ./tsunami cpu cores                show core count and exit`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// Read the three flags the user passed in
		duration := viper.GetInt("duration") // how many seconds to run
		load     := viper.GetInt("load")     // how hard to hit the CPU (1-100)
		cores    := viper.GetInt("cores")    // how many cores to stress

		// Ask Go how many CPU cores this machine has
		maxCores := runtime.NumCPU()

		// --- Check that the numbers make sense ---

		if load < 1 || load > 100 {
			return fmt.Errorf("--load must be between 1 and 100, you gave: %d", load)
		}

		if duration <= 0 {
			return fmt.Errorf("--duration must be bigger than 0, you gave: %d", duration)
		}

		// If the user did not pass -c, default to stressing all cores
		if cores <= 0 {
			cores = maxCores
		}

		// You cannot stress more cores than the machine actually has
		if cores > maxCores {
			return fmt.Errorf("this machine only has %d cores, you asked for %d", maxCores, cores)
		}

		// --- Tell the user what is about to happen ---
		fmt.Println("Starting CPU stress...")
		fmt.Printf("  Total cores on machine : %d\n", maxCores)
		fmt.Printf("  Cores being stressed   : %d\n", cores)
		fmt.Printf("  Load per core          : %d%%\n", load)
		fmt.Printf("  Duration               : %d seconds\n\n", duration)

		// --- Calculate spin and sleep times ---
		// Each window is 100ms total.
		// spinTime  = how long each goroutine burns CPU  (e.g. 30ms for 30% load)
		// sleepTime = how long it rests after that       (e.g. 70ms for 30% load)
		spinTime  := time.Duration(load)     * time.Millisecond
		sleepTime := time.Duration(100-load) * time.Millisecond

		// running is a simple true/false switch.
		// All goroutines keep looping while running == true.
		// When we set it to false they will stop on their next check.
		running := true

		// --- Start one goroutine for each core we want to stress ---
		// A goroutine is just a lightweight function running in the background.
		// Think of it like telling a worker "go do this job".
		for i := 0; i < cores; i++ {
			go func() {
				for running {
					// SPIN: burn CPU for spinTime
					// Keep looping until the spin deadline is reached
					spinDeadline := time.Now().Add(spinTime)
					for time.Now().Before(spinDeadline) {
						_ = 1 * 1 // dummy math so the compiler keeps this loop
					}

					// SLEEP: rest for sleepTime so we don't use 100% when load < 100
					// time.Sleep just pauses this goroutine for that amount of time
					if sleepTime > 0 {
						time.Sleep(sleepTime)
					}
				}
			}()
		}

		// --- Wait for the full duration ---
		// The main program just sits here while the goroutines do their work
		time.Sleep(time.Duration(duration) * time.Second)

		// --- Stop all goroutines ---
		// Setting running to false tells every goroutine to stop next time it checks
		running = false

		fmt.Println("Done! CPU stress finished.")
		return nil
	},
}

// ── cores subcommand ──────────────────────────────────────────────────────────
//
// This is a subcommand that lives under "cpu".
// Running "./tsunami cpu cores" just prints how many cores your machine has
// and then exits. No stress, no timers, just info.

var coresCmd = &cobra.Command{
	Use:   "cores",
	Short: "Show how many CPU cores this machine has",
	Run: func(cmd *cobra.Command, args []string) {
		// runtime.NumCPU() returns the number of logical CPU cores
		total := runtime.NumCPU()
		fmt.Printf("This machine has %d CPU core(s) available.\n", total)
		fmt.Printf("You can stress anywhere from 1 to %d core(s) with: tsunami cpu -c <number>\n", total)
	},
}

// ── init ──────────────────────────────────────────────────────────────────────
//
// init() is a special Go function that runs automatically when the program starts.
// We use it to:
//   1. Register the "cpu" command under the root "tsunami" command
//   2. Register the "cores" subcommand under "cpu"
//   3. Define all the flags (-l, -d, -c)
//   4. Tell viper about the flags so config.yaml values work too

func init() {
	// Register "cpu" so "./tsunami cpu" works
	rootCmd.AddCommand(cpuCmd)

	// Register "cores" under "cpu" so "./tsunami cpu cores" works
	cpuCmd.AddCommand(coresCmd)

	// -d / --duration  →  how many seconds to run (default: 10)
	cpuCmd.Flags().IntP("duration", "d", 10, "How many seconds to stress the CPU")

	// -l / --load  →  percent load per core (default: 100)
	cpuCmd.Flags().IntP("load", "l", 100, "How hard to stress each core, from 1 to 100")

	// -c / --cores  →  how many cores to stress (default: 0 which means all)
	cpuCmd.Flags().IntP("cores", "c", 0, "How many cores to stress (default: all cores)")

	// Connect each flag to viper so the config.yaml file can also set these values
	viper.BindPFlag("duration", cpuCmd.Flags().Lookup("duration"))
	viper.BindPFlag("load",     cpuCmd.Flags().Lookup("load"))
	viper.BindPFlag("cores",    cpuCmd.Flags().Lookup("cores"))
}
