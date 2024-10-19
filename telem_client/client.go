package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type TelemEntity struct {
	SatelliteId int32
	Temperature float64
	Voltage     float64
	Altitude    float64
}

func main() {
	fmt.Println("Hello, World!")

	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		User:     "postgres",
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	// url := "127.0.0.1:5005"
	// url := "host.docker.internal:5005"
	// url := "telem_server:5005"
	url := os.Getenv("PNT_HOST")
	log.Printf("client for server url: %s\n", url)

	addr, err := net.ResolveUDPAddr("udp", url)

	if err != nil {
		println("error resolved addr")
	}

	println("resolved address")

	conn, err := net.DialUDP("udp", nil, addr)

	println("dialed udp")

	if err != nil {
		println("error dialing udp")
	}

	msg := make([]byte, 512)

	type telemData struct {
		HeaderByte1    [1]byte
		HeaderByte2    [1]byte
		Unused         [2]byte
		SatelliteId    [4]byte
		Temperature    [4]byte
		BatteryVoltage [4]byte
		Altitude       [4]byte
	}

	line := "send pnt"

	n, err := conn.Write([]byte(line))
	if err != nil {
		println("error writing")
	}
	log.Printf("sent %d bytes \n", n)

	for {
		n, err = conn.Read(msg)
		if err != nil {
			println("error reading")
			println(err)
		}
		log.Printf("server sent: %d bytes \n", n)

		data := telemData{}
		err = binary.Read(bytes.NewBuffer(msg[:]), binary.BigEndian, &data)
		if err != nil {
			panic(err)
		}

		localSatId := binary.LittleEndian.Uint32(data.SatelliteId[:])
		localTemp := bytesToFloat32(data.Temperature[:])
		localBattVoltage := bytesToFloat32(data.BatteryVoltage[:])
		localAltitude := bytesToFloat32(data.Altitude[:])

		fmt.Println("Header1: ", data.HeaderByte1)
		fmt.Println("Header2: ", data.HeaderByte2)
		fmt.Println("Satellite ID: ", localSatId)
		fmt.Println("Temp: ", localTemp)
		fmt.Println("Batt Voltage: ", localBattVoltage)
		fmt.Println("Alt: ", localAltitude)

		entityData := &TelemEntity{
			SatelliteId: int32(localSatId),
			Temperature: float64(localTemp),
			Voltage:     float64(localBattVoltage),
			Altitude:    float64(localAltitude),
		}
		_, err = db.Model(entityData).Insert()

		if err != nil {
			fmt.Println("could not insert")
			fmt.Println(err)
			panic(err)
		}
	}
}

func bytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	err := db.Model((*TelemEntity)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
		Temp:        false,
	})

	if err != nil {
		fmt.Println("could not create table")
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("created Table")
	}
	return nil
}
