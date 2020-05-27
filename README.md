## Jight

This is a simplified version of our DagChain system. Specifically,
- All the wallet addresses are maintained by a genesis node (#1 node)
- All the transactions are issued by the genesis node (#1 node), and then are broadcast to the other nodes
- Each address is selected to issue a new transaction one by one, with the account 1 as the receiver.
- To simulate the different views of tips due to the network latency, the tips cited by a transaction will not be 
deleted from the tip collection immediately. Instead, these tips will be marked by their cited flags and be deleted 
from the tip collection later, by an outer cli command

## Branch to test data reduction approach with all the TxContent stored together and solely 

This branch is created to test est data reduction approach with all the TxContent stored together and solely.
The common experimental configuration of the experiments are as follows:
- The number of nodes is set as 10 (`config.GENESIS_ADDR_COUNT`), which can be adjusted later.
- Every second, each node will issue a new transaction, by the function `send_round` in `testCmd.sh`
- The network latency is set as 10 seconds (`config.NETWORK_LATENCY`). Thus the tips will be clean every 10 seconds by
   `./jight --rpcport 9525 refreshtips` in the function `send_10rounds` in `testCmd.sh`.

Before each group of experiments, the direcroty `Jightdb`, `JightdbMerging`, and `JightdbOthers`  must be removed.

- step 1: `jightd` is firstly run by the command `./jighd`.
- step 2: `jight` is run by the command `./testCmd.sh`