#!/bin/bash

#
# $1 = node number from Talisman BBS
#

cd /home/robbiew/go/src/github.com/robbiew/autosig
exec ./autosig $1
