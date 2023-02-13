package utils

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func StringPtrToAddressPtr(hex *string) *common.Address {
	if hex == nil {
		return nil
	}

	return Ptr(common.HexToAddress(*hex))
}

func AddressPtrToStringPtr(addr *common.Address) *string {
	if addr == nil {
		return nil
	}

	return Ptr(AddressToString(*addr))
}

func AddressToString(addr common.Address) string {
	return strings.ToLower(addr.String())
}
