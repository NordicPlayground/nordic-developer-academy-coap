#!/bin/bash

FILE=/run/secrets/STORAGE_CONNECTION_STRING     
if [ -f $FILE ]; then
export STORAGE_CONNECTION_STRING=$(cat /run/secrets/STORAGE_CONNECTION_STRING)
fi 

# Enumerate all network interfaces and get their assigned IP addresses
interfaces=$(ip -o addr show | awk '{print $2}' | cut -d':' -f1)

for interface in $interfaces; do
    IP=$(ip -o addr show $interface | awk '{print $4}' | cut -d'/' -f1)
    echo "Interface: $interface, Assigned IP: $IP"

    nohup /home/coap/coap-server -address "${IP}" -password "connect:anything" -dTLS &
    nohup /home/coap/coap-server -address "${IP}" &
done

sleep infinity