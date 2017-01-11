#!/bin/sh

temp=$(curl -s http://api.openweathermap.org/data/2.5/weather\?zip\=$1,us\&appid\=56b0b801438a43d1459ab2924781875f\&units\=imperial | jq .main.temp)

echo "The temperature in $1 is $temp"