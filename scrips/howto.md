Test Topologies
===============

# Simple Topology

Directory simple-topo contains scripts to create a simple topology. This creates four network namespaces
1. Two controller namespaces.
   - Sets up a `bridge` in the host namespace and creates two `veth` pairs connecting to controller namespaces
   - Enables IPv6 on these interfaces assign link local addresses
   - Adds two interfaces for connecting to nodes (might be through a bridge)
2. Two node namespaces
   - One end connecting to controller


It should simply be possible to do `ip netns exec <nsname> bash` and run commands. In particular, we'd want to run `gobgpd` on controller nodes (with unnumbered links). Run node agent and controller agent on our respective nodes.


