package userutil

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/howeyc/crc16"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {

	st := sonyflake.Settings{
		MachineID: getMachineIDFromEnv,
		StartTime: time.Unix(1642764826, 0), // 设id的起始时间：2022-01-21 19:33:46; 切记不要修改这个时间值，否则有ID重复的风险!!!
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func getMachineIDFromEnv() (uint16, error) {
	dc := os.Getenv("DC")
	ip := os.Getenv("HOST_ADDR")
	if dc == "" || ip == "" {
		return getMachineID()
	}

	checksum := crc16.ChecksumCCITT([]byte(dc + ip))
	fmt.Printf("dc: %s; ip: %s; machineID: %d\n", dc, ip, checksum)

	return checksum, nil
}

func getMachineID() (uint16, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP
				mid := uint16(ip[14])<<8 + uint16(ip[15])
				fmt.Printf("machineID: %d from ip: %s\n", mid, ipnet.String())
				return mid, nil
			}
		}
	}
	return 0, errors.New("can't get ip")
}

func GenID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

func Base58(id uint64) string {
	bs := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(bs, id)

	return base58.Encode(bs[:n])
}
