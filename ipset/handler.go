package ipset

import (
	"fmt"
)

func CheckSetExists(name string) error {
	for entry := range ch {
		fmt.Println(entry)
	}
	return nil
}
