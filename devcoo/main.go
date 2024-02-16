package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func OnInterrupt() (context.Context, context.CancelFunc) {
	wrap := func(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
		ctx, closer := context.WithCancel(ctx)

		c := make(chan os.Signal, 1)
		signal.Notify(c, signals...)

		go func() {
			select {
			case <-c:
				closer()
			case <-ctx.Done():
			}
		}()

		return ctx, closer
	}

	return wrap(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	ctx, cancel := OnInterrupt()
	defer cancel()

	jobIndex := 0
	jobs := make([][]string, 0)

	taskIndex := 0
	tasks := make([][]string, 0)

	inGroup := false
	task := false
	for _, arg := range os.Args[1:] {
		if arg == "}" {
			if !inGroup {
				fmt.Fprintf(os.Stderr, "syntax error: unexpected close bracket\n")
				os.Exit(1)
			}

			cnt := 0
			if task {
				cnt = len(tasks[taskIndex])
				taskIndex = taskIndex + 1
			} else {
				cnt = len(jobs[jobIndex])
				jobIndex = jobIndex + 1
			}

			if cnt == 0 {
				fmt.Fprintf(os.Stderr, "syntax error: empty command\n")
				os.Exit(1)
			}

			inGroup = false
			task = false
		} else if arg == "{" {
			if inGroup {
				fmt.Fprintf(os.Stderr, "syntax error: unexpected open bracket\n")
				os.Exit(1)
			}
			inGroup = true

			if task {
				tasks = append(tasks, []string{})
			} else {
				jobs = append(jobs, []string{})

			}
		} else {
			if !inGroup {
				if !task && arg == "once" {
					task = true
					continue
				} else {

					fmt.Fprintf(os.Stderr, "syntax error: command should be enclosed in { }\n")
					os.Exit(1)
				}
			}
			if task {
				tasks[taskIndex] = append(tasks[taskIndex], arg)
			} else {
				jobs[jobIndex] = append(jobs[jobIndex], arg)
			}
		}
	}

	if task {
		fmt.Fprintf(os.Stderr, "syntax error: attribute without body\n")
		os.Exit(1)
	}

	if inGroup {
		fmt.Fprintf(os.Stderr, "syntax error: missing close bracket at the end of the arguments\n")
		os.Exit(1)
	}

	for index, args := range tasks {
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "%d{%s}: %s\n", index, strings.Join(args, " "), err)
			os.Exit(1)
		}
	}
	var wg sync.WaitGroup
	for index, args := range jobs {
		wg.Add(1)
		go func(index int, args []string) {
			defer wg.Done()

			for {
				cmd := exec.CommandContext(ctx, args[0], args[1:]...)
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%d{%s}: %s\n", index, strings.Join(args, " "), err)
					return
				}

			}
		}(index, args)
	}

	wg.Wait()
}
