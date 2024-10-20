import random

import socket

import struct

import time

# Define CSP header and telemetry packet format (CSP header + telemetry)

# CSP header is packed into 2 bytes:

# - priority (2 bits) + destination (6 bits) are packed into 1st byte (B)

# - source (6 bits) + reserved (4 bits) + port (6 bits) + hmac (1 bit) + rdp (1 bit) are packed into 2nd byte (B)

# Then telemetry data:

# - satelliteID (uint32, 4 bytes)

# - temperature (float, 4 bytes)

# - batteryVoltage (float, 4 bytes)

# - altitude (float, 4 bytes)

csp_telemetry_format = (

    "B B I f f f"  # CSP header fields packed into two bytes, telemetry follows

)


def generate_random_telemetry():

    # Generate random CSP header and telemetry data

    priority = random.randint(0, 3)  # 2 bits for priority

    destination = random.randint(0, 63)  # 6 bits for destination address

    source = random.randint(0, 63)  # 6 bits for source address

    reserved = 0  # Reserved bits (4 bits)

    port = random.randint(0, 63)  # 6 bits for destination port

    hmac = random.randint(0, 1)  # 1 bit HMAC flag

    rdp = random.randint(0, 1)  # 1 bit RDP flag

    # Merge into two bytes (B format in struct)

    header_byte1 = (priority << 6) | destination

    header_byte2 = (

        (source << 2) | reserved | port >> 4

    )  # Shift port to fit remaining 2 bits in byte2

    header_byte3 = ((port & 0x3F) << 2) | (hmac << 1) | rdp  # Pack port and flags

    # Random telemetry data

    satellite_id = random.randint(1000, 9999)

    temperature = random.uniform(-100.0, 100.0)

    battery_voltage = random.uniform(0.0, 100.0)

    altitude = random.uniform(200.0, 400.0)

    return (

        header_byte1,

        header_byte2,

        satellite_id,

        temperature,

        battery_voltage,

        altitude,

    )


def send_csp_telemetry_packet(sock, address):

    while True:

        csp_telemetry_data = generate_random_telemetry()

        packet = struct.pack(csp_telemetry_format, *csp_telemetry_data)

        sock.sendto(packet, address)

        print(f"Sent CSP telemetry: {csp_telemetry_data}")

        time.sleep(1)  # Delay to simulate real-time telemetry


# Set up UDP socket

udp_ip = ""

udp_port = 5005
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

while True:
    try:
        # server must bind to an ip address and port
        sock.bind((udp_ip, udp_port))
        print("Listening on Port:", udp_port)
        break
    except Exception as e:
        print(e)
        print("ERROR: Cannot connect to Port:", udp_port)

while True:
    message, addr = sock.recvfrom(1024)  # OK someone pinged me.
    print(f"received from {addr}: {message}")
    break

# Start sending telemetry

send_csp_telemetry_packet(sock, addr)
