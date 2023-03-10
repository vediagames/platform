package domain

import (
	"fmt"
	"net"
	"time"
)

type Session struct {
	ID         string
	PageURL    string
	IP         IP
	Device     Device
	CreatedAt  time.Time
	InsertedAt time.Time
}

type IP string

func (i IP) Validate() error {

	parsedIPAddress := net.ParseIP(string(i))
	if parsedIPAddress == nil || parsedIPAddress.To4() == nil {
		return fmt.Errorf("invalid value: %q", i)
	}

	return nil
}

func (ip IP) String() string {
	return string(ip)
}

type Device string

const (
	DeviceLaptop  = Device("laptop")
	DeviceDesktop = Device("desktop")
	DeviceTablet  = Device("tablet")
	DevicePhone   = Device("phone")
	DeviceTV      = Device("tv")
)

func (d Device) Validate() error {
	switch d {
	case DeviceLaptop, DeviceDesktop, DeviceTablet, DevicePhone, DeviceTV:
		return nil
	default:
		return fmt.Errorf("invalid value: %q", d)
	}
}

func (d Device) String() string {
	return string(d)
}
