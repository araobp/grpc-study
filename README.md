# DevOps for networking containers

This project re-uses outputs from my other project "[neutron-lan](https://github.com/araobp/neutron-lan)".

##Background and motivation

- As my hobby, I just want to develop a very simple DevOps framework for networking containers for several use cases.
- I don't have a chance to write code at work...
- I want a test bed (at a very low cost, under $100) to run YANG, grpc, docker, BGP/OSPF, Open vSwitch, OVSDB... Raspberry Pi is the best for such a purpose.
- I need to migrate from Java/Python to Golang for some reasons.

##Current status(2016/01/16)
- Still unstable...
- [Findings](./doc/FINDINGS.md)

####Direction from now on
- etcd and gRPC are not suitable for this project.
- I am going to use [tega](https://github.com/araobp/tega) to simplify the architecture, taking an event-driven/model-driven approach.
- I use protobuf as a schema lang.
-- I don's use YANG any longer, because YANG is not compatible with Python dict (thus incompatible with my tega db). YANG is for Java developers, since Java does not natively support hash-map.
- NLAN agents are written in Golang, whereas other components may be written in Python etc.

##Go version
1.5.3

##NLAN services
- PTN: Packet Transport Network (Layer 1 and Layer 2)
- DVR: Distributed Virtual Switch and Distributed Virtual Router (Layer 2 and Layer 3)
- vHosts: netns-based virtual hosts
- Router: Quagga configuration

To be added:
- Links: direct linking(tun/tap)
- Bridges: non-distributed virtual switch
- VRF: netns-based VRF
- container-macvlan direct linking (skipping docker0)
- [netlink](https://github.com/milosgajdos83/tenus) in addition to "ip" and "brctl"

##Target use cases

Use case 1 has already been implemented, and use case 2 is being planned at the moment.

### Use case 1: Network simulation

This use case makes use of NLAN's PTN, vHosts and Router services.
![WAN simulation](https://docs.google.com/drawings/d/1VKfKlwnzWQ2-ImfXeB5uNegGBK0BnaGU_4lS8h4Qpcw/pub?w=640&h=480)

####Declarative state representations:
- [ptn-bgp.yaml](./etc/ptn-bgp.yaml)
- [ptn-ospf.yaml](./etc/ptn-ospf.yaml)

####Running the simulated network on Raspberry Pi
This is sort of micro NFV(Network Function Virtualization) on a single Rapsberry Pi.
- Nine virtual routers (Linux containers)
- Sixteen virutal hosts (netns)

You can learn how routing protocols work on this simulated network.

[Setting up the software on Raspberry Pi](./doc/RPI.md)

Log in the virtual routers with ssh, and try "ip" or "vtysh" commands:
- ip route
- ip addr
- ip link
- ip netns
- vtysh: show run
- vtysh: show ip route
- vtysh: show ip bgp
     :

####Quagga and GoBGP:
This use case makes use of Quagga, but [gobgp](https://github.com/osrg/gobgp) may optionally be used as Route Reflector or Route Server on "RR" container in the fig above.
- [gobgpd.conf](./etc/gobgpd.conf)

You can also launch gobgpd from NLAN agent by including "EmbeddedBgp: true" in your NLAN state file:
```
      Router:
        Loopback: 10.1.1.5/32
        EmbeddedBgp: true
        Bgp:
          - As: 100
            Neighbors:
              - Peer: 10.200.1.101
                RemoteAs: 100
                RouteReflectorClient: true
              - Peer: 10.200.1.102
                RemoteAs: 100
                RouteReflectorClient: true
              - Peer: 10.200.1.103
                RemoteAs: 100
                RouteReflectorClient: true
              - Peer: 10.200.1.104
                RemoteAs: 100
                RouteReflectorClient: true
```
- [Using gobgp command](./doc/GOBGP.md)

### Use case 2: SOHO NFV (Network Functions Virtualization)

This is the next use case I am going to work on... (as my hobby: not so practical)

![SONO-NFV](https://docs.google.com/drawings/d/11fJUimZVrGxqAdq-hJK4abDu0ZThkfHGtbl_94zW0rQ/pub?w=640&h=480)

####Raspberry Pi cluster
You could make a super computer with a cluser of RPIs. Here is a example: 
* http://likemagicappears.com/projects/raspberry-pi-cluster/

Why not develop your own NFV infrastructure with a bunch of RPIs?

##Network simulation with Linux containers
I use Linux containers as virtual routers, and this tool will set up virtual links (L0/L1) and virtual switches (L2) over the containers. Then I will run Quagga/Zebra(L3) daemons over the virtual routers to study how legacy routing protocols work.
- [An example of such a network](https://camo.githubusercontent.com/3f15c9634b2491185ec680fa5bb7d19f6f01146b/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f31564b664b6c776e7a5751322d496d6658654235754e656747424b30426e6147555f346c53386834517063772f7075623f773d39363026683d373230)
- [Working with Docker for network simulation](https://camo.githubusercontent.com/77cf473ea9499432e57b06a951f5f5248419f9e1/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f313631426e383077384a5a4b513742586d496f306272377851346b71456442635f585a3235347a754f5253552f7075623f773d36383026683d343030)

##Architecture

[Architecture](./doc/ARCHITECTURE.md)

##NLAN model in YANG and protobuf
Go stub generation
```
                             ___________  
                            /rpc model / proto3
                            ~~~~~~~~~~~
                                 |
                               merge
                                 |
                                 V
 ___________                _______________                _____________
/YANG model/ == goyang ==> /Protobuf model/ == protoc ==> /Go gRPC stub/
~~~~~~~~~~~                ~~~~~~~~~~~~~~~                ~~~~~~~~~~~~~
                               proto3
```
- [YANG model](./model/nlan/nlan.yang)
- [protobuf model](./model/nlan/nlan.proto)
- [rpc model](./model/nlan/rpc.proto)
- [Go gRPC stub](./model/nlan/nlan.pb.go)

##Key technologies used for this project
- [go](https://github.com/golang/go)
- Open vSwitch, OVSDB/JSON-RPC([RFC7047](https://tools.ietf.org/html/rfc7047))
- [vxlan](https://tools.ietf.org/html/rfc7348)
- [docker](https://github.com/docker/docker)
- [YANG](https://tools.ietf.org/html/rfc6020)/[goyang](https://github.com/openconfig/goyang)/[pyang](https://github.com/mbj4668/pyang)
- JSON/YAML
- Protocol buffers and gRPC
- quagga
- [gobgp](https://github.com/osrg/gobgp)

##Interesting
- [Cumulus Linux](https://cumulusnetworks.com/)
- [socketplane](https://github.com/socketplane/socketplane)
- [goplane](https://github.com/osrg/goplane)
- [bgp sdn](https://tools.ietf.org/html/draft-lapukhov-bgp-sdn-00)
- [nsq](https://github.com/nsqio/nsq)

#NLAN installation

[Step 1] Make a Docker image named "router" following the instruction [here](./docker/SETUP.md).

[Step 2] Start tega db:
```
$ cd scripts
$ ./tegadb &
```
[Step 3] Assign (secondary) IP address to each container:
```
$ cd scripts
$ ./cli
tega CLI (q: quit, h:help)
[tega: 0] put nlan.ip.pe1
10.10.10.1/24

[tega: 1] put nlan.ip.pe2
10.10.10.2/24

[tega: 2] put nlan.ip.pe3
10.10.10.3/24

[tega: 3] put nlan.ip.pe4
10.10.10.4/24

[tega: 4] put nlan.ip.rr
10.10.10.5/2

[tega: 5] put nlan.ip.ce1
10.10.10.6/2

[tega: 6] put nlan.ip.ce2
10.10.10.7/2

[tega: 7] put nlan.ip.ce3
10.10.10.8/2

[tega: 8] put nlan.ip.ce4
10.10.10.9/2

[tega: 9] get nlan.ip
{ce1: 10.10.10.6/2, ce2: 10.10.10.7/2, ce3: 10.10.10.8/2, ce4: 10.10.10.9/2, pe1: 10.10.10.1/24,
  pe2: 10.10.10.2/24, pe3: 10.10.10.3/24, pe4: 10.10.10.4/24, rr: 10.10.10.5/2}

```
[Step 4] Run the following shell script to build Docker image with NLAN agent embedded and to start the containers:
```
$ ./setup.sh
```

[Step 5]
Try this to deploy "PTN/Vhost/Router" services:
```
$ ./master/master.sh ptn-bgp
```
The script sets up [this network](https://camo.githubusercontent.com/3f15c9634b2491185ec680fa5bb7d19f6f01146b/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f31564b664b6c776e7a5751322d496d6658654235754e656747424b30426e6147555f346c53386834517063772f7075623f773d39363026683d373230).

[Step 6]
Open ssh session to the containers:
```
$ cd scripts 
$ ./ssh.sh pe1
       :
$ ./ssh.sh ce1
       :
```
The password is "root".

Then you can do whatever you try.

#Development environment setup

##Building Golang and protobuf for 32bit Linux
I use a very old PC with a 32bit CPU at home, so I need to build 32bit binary from source codes:
- Go lang installation: https://golang.org/dl/
- Protobuf build and installation: https://github.com/google/protobuf/blob/master/INSTALL.txt
```
$ autoconf
$ ./autogen.sh
$ ./configure
$ make
$ make install
```
- Add /usr/local/lib to LD_LIBRARY_PATH
```
$ export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARLY_PATH

```

##Go plugin for vim

Install [vim-go](https://github.com/fatih/vim-go) to your vim.

