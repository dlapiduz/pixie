#!/bin/sh

temp=$(curl -s http://api.openweathermap.org/data/2.5/weather\?zip\=$1,us\&appid\=b1b15e88fa797225412429c1c50c122a\&units\=imperial | jq .main.temp)

echo "The temperature in $1 is $temp"