kube-toolkit  - Toolkit for creating gRPC-based CLIs for Kubernetes in Go 
=========================================================================

![](https://raw.githubusercontent.com/ashleymcnamara/gophers/master/BATMAN_GOPHER.png)

> Image from [@ashleymcnamara](https://github.com/ashleymcnamara/gophers)'s gopher artwork - [license](https://github.com/ashleymcnamara/gophers/blob/master/LICENSE)


About
-----

If you ever used [Helm](https://github.com/kubernetes/helm) or [Draft](https://github.com/azure/draft), you know they are very cool command-line tools that connect to a Kubernetes cluster, more specifically to a server-side componend (Tiller in the case of Helm, Draftd for Draft) without exposing ports on the Internet, and allow you to interact with your cluster through gRPC-based services.

This repo aims to help you build similar services by allowing you to start from a pre-configured tool

Architecture
------------

[`kube-toolkit`](https://github.com/radu-matei/kube-toolkit) has two major components:

- `ktk` (short for Kubernetes ToolKit) - client that you install locally

- `ktkd` (short for Kuberentes ToolKit Daemon) - server-side component that is deployed on your Kubernetes cluster 


The `kube-toolkit` client (`ktk`) interacts with the server-side component (`ktkd`) using the Kubernetes API to create authenticated tunnels back to the cluster, using gRPC as the underlying communication protocol. The server runs as a pod in Kubernetes, and since it is only a starting point for future tools, it only knows how to return its version, and showcases how to stream data back to the client.

In order to communicate with the cluster you need to pass the `kubeconfig` file, and the tool will start a tunnel to the cluster for each command you execute, then will tear it down so there are no open connections to the cluster when no command is executed.


> Please note that there are still lots of things left to add, such as SSL, RBAC support or state management - you are more than welcome to contribute in any way to the project!

Disclaimer
----------
This is not an official Microsoft project, and all credits go to the awesome people building Helm and Draft, on which kube-toolkit is based.


Contributing
------------

Any idea (here, on Twitter - @Matei_Radu), issue or pull request is highly appreciated. Contribution guidelines will follow once there is a structure to this project.