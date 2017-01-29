# Presentation

## Networktopology

All elevators are connected to a dynamic master-slave TCP network. Every node in the network (elevator) can be in one of two states: Master or slave. The master is active and handles order synchronization, events and order assignments, for all nodes including itself. The slaves are passive and only forward all events and acts locally on orders assigned from master.

![alt text](network.png "A simple network of 4 nodes in it's idle state")

![alt text](network_example.png "Example of merging of multiple sub-networks")

The master is dynamically chosen to be the node on the local network with the lowest ip-address. Because the network may consist of multiple sub-networks (multiple masters), the master is chosen by merging the sub-networks using the following simple algoritm:



A network may consist of multiple sub-network and therefore multiple masters. All nodes are initially masters. The masters listen for TCP connections and repeatedly broadcast their ip-address over UDP. All masters listen for these packages. If a master receives such a package from a sender with ip-address lower than it's own, it closes it's connection to it's slaves and connects to the sender's ip. The node is now a slave until it is disconnected from it's master.

## Modularization

### Orders

* Global -- All orders in the network. All nodes should have identity global orders. The master handles synchronization of these orders.

* Local -- Orders the local node are assigned to do by the master

## Error handling

* Flipped bit -- In the network's communication, TCP has checksum to verify no flipped bits, and resends if checksum is not satisfied.

* Power loss -- If a node, master or slave, loses power, every other node has a copy of the global orders such that the system can recover. If the disconnected node is a slave, the master will simply detect the disconnection and redistribute all global orders to the remaining nodes (including itself). If the disconnected node is a master, all it's slaves will become masters and at first act on all the orders by their self. Eventually they will merge to a single master (by protocol) and the new master will redistribute the orders.

* Disconnect -- If a node disconnects, the network will handle the situation as if the node is dead (power loss). The single disconnected node will become master and act as if all other nodes are dead. In this situation multiple elevators may process the same orders, but no orders will be lost (not acted on).

* Elevator hangs -- If the physical elevator does not reach the next floor before a given timeout, the local program is terminated and must be investigated by an operator before manually restarted. The other nodes will act as if the given elevator is dead, and will redistribute all global orders.

* UDP packet never arrive -- This will cause no problem because UDP is only used for broadcasting repeatedly. If one does not arrive, eventually one of the following will.
