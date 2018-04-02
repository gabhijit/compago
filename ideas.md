Let's start by asking a simple question -

1. There will be two types of nodes
   - controller nodes
   - worker/compute nodes

2. Controller nodes will run
   - A controller agent (gRPC)
   - GoBGP - For route advertisement etc.
   - etcd nodes

3. worker/compute nodes will run
   - A node agent


A controller agent
  - will listen for updates from the worker/node agents for following -
    - Worker registration (worker tells it's own reachable IP address)
    - New links/addresses/routes added
    - Other data like - metrics / statistics etc.
  - Send Following directives to worker/node agent
    - Create a tenant network
    - Policy Messages (details TBD)

A worker/node agent will -
  - Listen on one side for netlink messagges filter/appropriately format them
    if required and send them to controller.
  - Listen for actions to be performed. The actions may include
    - Setup a network (create a bridge and VTEP)
    - Apply policy to endpoints
  - Periodically gather statistics and send it to controller nodes


Why are we not simply taking project calico and work around it?

The reason being, project calico's felix agent is quite a heavy weight agent
running on every node and we want to have as much light weight agent on the
worker node as possible. Plus computation of policy to be applied on a worker
node can be done centrally and simply actions can be pushed. This saves compute
cycles on the worker node. Also, we can separate out the BGP node from the
calico node and make it kind of centralized.

So this is almost like opencontrail but we are using components from standard
Linux kernel, so no separate modules are required.

Ideally, it should be possible to support a mix and match of container and VM
end points on the same worker node.

Let's run through some simple use-cases and see how this can be achieved with
use-cases above -

1. Kubernetes pod is created -

Each kubernetes pod will run in it's own network namespace.

We will provide cni plugin that will do the setting up of network for container.

This should create a veth pair in the kernel and IP addresses using IPAM will
be assigned to one end of veth. This information will be picked up by the
worker agent through netlink messages and sent to the controller. The
controller would take right action about it eg. adding a BGP route to this.


2. When a VM is created -

VM will have one or more interfaces on one or more networks. Each network will
be identified by a VNI. So when a port is added to the VM, controller agent
will send following information to the worker/node agent. Port network - VM.

The worker/node agent using netlink will perform following - By default a
network for a VM will be opaque L2 network through an EVPN.

