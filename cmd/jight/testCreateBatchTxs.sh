#!/bin/bash

send_1round() {
    ./jight --rpcport 9525  createbatchtxs --times=1000 --speed=10
    du -B K ../jightd/Jightdb ../jightd/JightdbMerging ../jightd/JightdbOthers
}
send_10rounds() {
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
    send_1round
}

send_100rounds() {
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
    send_10rounds
}

send_100rounds
send_100rounds
send_100rounds
