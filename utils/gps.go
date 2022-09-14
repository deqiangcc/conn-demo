package utils

import "time"

type GPS struct {
	Timestamp        *uint32 `protobuf:"varint,1,opt,name=Timestamp,def=0" json:"Timestamp,omitempty"`
	Latitude         *int32  `protobuf:"zigzag32,2,opt,name=Latitude,def=2000000000" json:"Latitude,omitempty"`
	Longitude        *int32  `protobuf:"zigzag32,3,opt,name=Longitude,def=2000000000" json:"Longitude,omitempty"`
	Altitude         *int32  `protobuf:"zigzag32,4,opt,name=Altitude,def=0" json:"Altitude,omitempty"`
	Quality          *int32  `protobuf:"varint,5,opt,name=Quality,def=0" json:"Quality,omitempty"`
	Mode             *int32  `protobuf:"varint,6,opt,name=Mode,def=0" json:"Mode,omitempty"`
	PDOP             *int32  `protobuf:"varint,7,opt,name=PDOP,def=0" json:"PDOP,omitempty"`
	HDOP             *int32  `protobuf:"varint,8,opt,name=HDOP,def=0" json:"HDOP,omitempty"`
	VDOP             *int32  `protobuf:"varint,9,opt,name=VDOP,def=0" json:"VDOP,omitempty"`
	SatellitesUsed   *int32  `protobuf:"varint,10,opt,name=SatellitesUsed,def=0" json:"SatellitesUsed,omitempty"`
	SatellitesInView *int32  `protobuf:"varint,11,opt,name=SatellitesInView,def=0" json:"SatellitesInView,omitempty"`
	Speed            *int32  `protobuf:"varint,12,opt,name=Speed,def=0" json:"Speed,omitempty"`
	Course           *int32  `protobuf:"zigzag32,13,opt,name=Course,def=0" json:"Course,omitempty"`
	HACC             *int32  `protobuf:"varint,14,opt,name=HACC,def=0" json:"HACC,omitempty"`
	VACC             *int32  `protobuf:"varint,15,opt,name=VACC,def=0" json:"VACC,omitempty"`
	FixTime          *int32  `protobuf:"varint,16,opt,name=FixTime,def=0" json:"FixTime,omitempty"`
	SampleType       *int32  `protobuf:"varint,20,opt,name=SampleType,def=0" json:"SampleType,omitempty"`
	GeoidAltitude    *int32  `protobuf:"zigzag32,21,opt,name=GeoidAltitude,def=0" json:"GeoidAltitude,omitempty"`
}

func NewGpss(n int) []GPS {
	now := uint32(time.Now().Unix())
	intDefault := int32(200)
	gps := GPS{
		Timestamp:        &now,
		Latitude:         &intDefault,
		Longitude:        &intDefault,
		Altitude:         &intDefault,
		Quality:          &intDefault,
		Mode:             &intDefault,
		PDOP:             &intDefault,
		HDOP:             &intDefault,
		VDOP:             &intDefault,
		SatellitesUsed:   &intDefault,
		SatellitesInView: &intDefault,
		Speed:            &intDefault,
		Course:           &intDefault,
		HACC:             &intDefault,
		VACC:             &intDefault,
		FixTime:          &intDefault,
		SampleType:       &intDefault,
		GeoidAltitude:    &intDefault,
	}

	var gpss []GPS
	for i := 1; i <= n; i++ {
		gpss = append(gpss, gps)
	}

	return gpss
}
