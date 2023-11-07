package validators

import (
	"github.com/ttacon/libphonenumber"
)

func MobileNumberValidator(mobile_number string) error {
	if _, err := libphonenumber.Parse(mobile_number, libphonenumber.CAPTURING_EXTN_DIGITS); err != nil {
		return err
	}
	return nil
}
