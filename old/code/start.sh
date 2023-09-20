#!/bin/bash

if [ $# -ne 4] 
then 
    echo "Usage: ${0} NMAPPERS MAXITER THRESHOLD(1 = 0.001%)"
    exit 1
fi

if [ $1 -lt 0 -o $1 -gt 99 ]
then
	echo "The number of mappers isn't correct or is too high."
	exit 1
fi


if [ $2 -lt 0 -o $2 -gt 1000 ]
then
	echo "The number of max iterations isn't correct or is too high."
	exit 1
fi

if [ $3 -lt 0 -o $3 -gt 999 ]
then
	echo "DeltaThreshold must be included in [0,100]."
	exit 1
fi

NUMMAP=${1} MAXITER=${2} THRESHOLD=${3} docker compose up

