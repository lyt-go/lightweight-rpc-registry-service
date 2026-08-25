package store

import "fmt"

func StreamNodeNames(names []string, out chan<- string, errs chan<- error) {
	for _, name := range names {
		if name == "" {
			errs <- fmt.Errorf("node name is empty")
			return
		}
		out <- name
	}
}
