#!/bin/bash

FILE=/run/secrets/STORAGE_CONNECTION_STRING     
if [ -f $FILE ]; then
    export STORAGE_CONNECTION_STRING=$(cat /run/secrets/STORAGE_CONNECTION_STRING)
fi 

IPv4_addresses=(`ip -4 addr show | grep -oP '(?<=inet\s)\d+(\.\d+){3}'`)
for address in $IPv4_addresses; do
    nohup /home/coap/coap-server -address $address -network udp4 -password "connect:anything" -dTLS &
    nohup /home/coap/coap-server -address $address -network udp4 &
done

IPv6_addresses=(`ip -6 addr show | grep -oP '(?<=inet6\s)[\da-fA-F:]+'`)
for address in $IPv6_addresses; do
    nohup /home/coap/coap-server -address [$address] -network udp6 -password "connect:anything" -dTLS &
    nohup /home/coap/coap-server -address [$address] -network udp6 &
done

sleep 1

tail -f nohup.out