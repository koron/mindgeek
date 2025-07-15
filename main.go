package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func readline() (string, error) {
	r := bufio.NewReader(os.Stdin)
	l, p, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	if p {
		for {
			_, p, err := r.ReadLine()
			if err != nil {
				return "", err
				if !p {
					break
				}
			}
		}
	}
	return string(l), nil
}

func prompt(msg string, choices ...string) (string, error) {
	for {
		fmt.Printf("%s\n? ", msg)
		s, err := readline()
		if err != nil {
			return "", err
		}
		if len(choices) == 0 {
			return s, nil
		}
		for _, c := range choices {
			if s == c {
				return s, nil
			}
		}
		fmt.Println()
	}
}

func modeLoop(round int, name string, oneGame func() (bool, error)) error {
	var score = 0
	fmt.Printf("\n%s mode\n", name)
	for i := range round {
		fmt.Printf("Round #%d/%d, Score:%d\n", i+1, round, score)
		ok, err := oneGame()
		if err != nil {
			return err
		}
		if ok {
			score++
		}
	}
	fmt.Printf("%s total score: %d\n", name, score)
	return nil
}

func clairvoyanceOne() (bool, error) {
	want := strconv.Itoa(rand.IntN(5) + 1)
	got, err := prompt("See through the numbers 1 to 5", "1", "2", "3", "4", "5")
	if err != nil {
		return false, err
	}
	if want != got {
		fmt.Printf("Incollect: want=%s\n", want)
		return false, nil
	}
	fmt.Println("Correct")
	return true, nil
}

func psychokinesisOne() (bool, error) {
	move := &atomic.Bool{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ok := rand.IntN(2) != 0
			move.Store(ok)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	_, err := prompt("When you feel it move, press Enter")
	if err != nil {
		return false, err
	}
	if !move.Load() {
		fmt.Println("Stay")
		return false, nil
	}
	fmt.Println("Moved")
	return true, nil
}

func precognitionOne() (bool, error) {
	got, err := prompt("Predict the numbers 1 to 5", "1", "2", "3", "4", "5")
	if err != nil {
		return false, err
	}
	want := strconv.Itoa(rand.IntN(5) + 1)
	if want != got {
		fmt.Printf("Incollect: want=%s\n", want)
		return false, nil
	}
	fmt.Println("Correct")
	return true, nil
}

func run() error {
	for {
		c, err := prompt("Choose your training:\n  1. Clairvoyance\n  2. Psychokinesis\n  3. Precognition", "1", "2", "3")
		if err != nil {
			return err
		}
		switch c {
		case "1":
			if err := modeLoop(10, "Clairvoyance", clairvoyanceOne); err != nil {
				return err
			}
		case "2":
			if err := modeLoop(10, "Psychokinesis", psychokinesisOne); err != nil {
				return err
			}
		case "3":
			if err := modeLoop(10, "Precognition", precognitionOne); err != nil {
				return err
			}
		}
		fmt.Println()
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
