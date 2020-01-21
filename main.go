// A utility simulates a memory leak for testing, diagnostic purposes
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	// define and set default command parameter flags
	var dFlag = flag.Int("d", 100, "Optional: delay is ms to adjust the leak rate; default is 100")
	var lFlag = flag.Int("l", 1048576, "Optional: limit the leak to this many MiBs")
	var hFlag = flag.Bool("h", false, "print usage information")

	// usage function that's executed if a required flag is missing or user asks for help (-h)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [-d <leak delay in ms; deafaults to 100> -l <leak limit in MiB>]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	//provide help (-h)
	if *hFlag == true {
		flag.Usage()
		os.Exit(0)
	}

	// a channel to use a hold when the memory limit it reached
	hold := make(chan bool)

	// a spinner that displays how much memory has leaked and when it holding
	go func(hold chan bool, lFlag int) {
		s := spinner.New(spinner.CharSets[35], 250*time.Millisecond)
		for {
			mem := memUsage()
			s.Prefix = fmt.Sprintf("Leaked: %d MiB ", mem)
			s.Start()
			s.Color("magenta")
			time.Sleep(1 * time.Second)
			s.Restart()
			// if we've reached the limit, update display and hold
			if mem >= uint64(lFlag) {
				s.Color("green")
				s.Prefix = fmt.Sprintf("Holding at %d MiB ", mem)
				s.UpdateCharSet(spinner.CharSets[28])
				s.UpdateSpeed(1 * time.Second)
				s.Restart()
				<-hold
			}
		}
	}(hold, *lFlag)

	// Although the "leak" var should contiue to grow, the GC is somehow getting in the way, disabling
	debug.SetGCPercent(-1)
	var leak string
	KB := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
			abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`

	// start leaking indefinitely unless a limit has been provided and met
	for {
		leak += KB
		time.Sleep(time.Duration(*dFlag) * time.Millisecond)
		mem := memUsage()
		// if we've reached the limit, hold
		if mem >= uint64(*lFlag) {
			<-hold
		}
	}
}

func memUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return bToMb(m.Alloc)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
