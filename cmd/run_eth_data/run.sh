#!/bin/bash
echo "--------------- Process 0th file ----------------"
./run_eth_data 0 9625 > log0.txt 2>&1

echo "--------------- Process 1st file ----------------"
./run_eth_data 1 9625 > log1.txt 2>&1

echo "--------------- Process 2nd file ----------------"
./run_eth_data 2 9625 > log2.txt 2>&1

echo "--------------- Process 3rd file ----------------"
./run_eth_data 3 9625 > log3.txt 2>&1

echo "--------------- Process 4th file ----------------"
./run_eth_data 4 9625 > log4.txt 2>&1

echo "--------------- Process 5th file ----------------"
./run_eth_data 5 9625 > log5.txt 2>&1

echo "--------------- Process 6th file ----------------"
./run_eth_data 6 9625 > log6.txt 2>&1

echo "--------------- Process 7th file ----------------"
./run_eth_data 7 9625 > log7.txt 2>&1

echo "--------------- Process 8th file ----------------"
./run_eth_data 8 9625 > log8.txt 2>&1
