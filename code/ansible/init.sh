#!/bin/bash

if [ $# -ne 2 ]
then
    echo "Usage: ${0} IP_EC2_INSTANCE KEY_PATH"
    exit 1
fi

ansible-playbook -vvvv --private-key=${2} -i ${1}, deploy.yaml  
ssh -i ${2}  ubuntu@${1}       

