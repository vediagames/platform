package domain

import (
	"net"
	"regexp"
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

func (i IP) IsValid() bool {
	regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !regex.MatchString(string(i)) {
		return false
	}

	parsedIPAddress := net.ParseIP(string(i))
	if parsedIPAddress == nil || parsedIPAddress.To4() == nil {
		return false
	}

	return true
}

func (ip IP) String() string {
	return string(ip)
}

type Device string

const (
	Laptop Device = "laptop"
	Tablet Device = "tablet"
	Phone  Device = "phone"
	TV     Device = "tv"
)

var AllDevice = []Device{
	Laptop,
	Tablet,
	Phone,
	TV,
}

func (d Device) IsValid() bool {
	switch d {
	case Laptop, Phone, TV, Tablet:
		return true
	}
	return false
}

func (d Device) String() string {
	return string(d)
}
