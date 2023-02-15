package file

import (
	"bufio"
	"fmt"
	"os"
)

func SaveLines(name string, lines []string) (retError error) {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil && retError == nil {
			retError = err
		}
	}()
	for _, line := range lines {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	return nil
}

func ForEachLine(name string, callback func(text string) error) error {
	in, err := os.Open(name)
	if err != nil {
		return err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err := callback(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}
