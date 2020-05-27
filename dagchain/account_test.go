package dagchain

import (
	"fmt"
	"testing"
)

func TestCreateNewAccount(t *testing.T) {
	acc := CreateNewAccount("158ewJ1itTAAE8gy1Hk45JdAmqvbpd0006", 100)
	accBytes := acc.Serialize()
	fmt.Printf("Length of accBytes: %d\n", len(accBytes))
	fmt.Println(accBytes)

	fmt.Printf("Length of accountBytes: %d\n", len(acc.Account))
	fmt.Println(acc.Account)
}