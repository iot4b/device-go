#!/bin/sh

ansible-playbook -i ansible/inventory/test.yml --limit=$1 ansible/playbooks/deploy_test.yml -vv