package mail

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateNumericOTP returns a cryptographically secure 6-digit numeric string
func GenerateNumericOTP() (string, error) {
	// Max value for 6 digits is 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// Formats numbers like 42 into "000042" to maintain strict 6-digit sizing
	return fmt.Sprintf("%06d", n.Int64()), nil
}
