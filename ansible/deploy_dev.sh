#!/bin/sh

ansible-playbook -i ansible/inventory/dev.yml --limit=$1 ansible/playbooks/deploy_dev.yml -vv