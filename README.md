Joker
=====
![](https://raw.githubusercontent.com/ashleymcnamara/gophers/master/BATMAN_GOPHER.png)

> Image from [@ashleymcnamara](https://github.com/ashleymcnamara/gophers)'s gopher artwork - [license](https://github.com/ashleymcnamara/gophers/blob/master/LICENSE)

About
-----

This (very much work in progress) project is an experiment on creating and managing cloud resources from Kubernetes.

Here's the main idea behind this project: **most people that use Kubernetes with a cloud provider most often end up using other services from that cloud provider - databases, caches, queues and others**.

The problem with this flow - in order to use any of the before-mentioned services, you need to first deploy / enable that service from a portal / CLI / deployment file, then get the access keys for that service, then create a Kubernetes secret out of them. After this, you can go and create your Kubernetes deployments that use those services.

> The steps above are literally [described in the Google Cloud Platform on using Kubernetes with Pub/Sub](https://cloud.google.com/container-engine/docs/tutorials/authenticating-to-cloud-platform) and in [the Azure documentation on monitoring Kubernetes using OMS](https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-kubernetes-oms)

The idea is to eliminate [Yak Shaving](https://www.hanselman.com/blog/YakShavingDefinedIllGetThatDoneAsSoonAsIShaveThisYak.aspx) and just have your resources deployed and have all connection strings / keys as objects that you, as someone who knows Kubernetes, can understand - secrets. Then, just delete them just as you would delete other Kubernetes resources.


Moreover, if you happen to use more than one cloud provider, you would want the same tool to deploy and manage your services.

**The end game is to deploy and manage your cloud resources in a similar manner to Kubernetes services**.


How I envision working with `joker`
----------------------------------

Working with `joker` must be familiar if you used `helm` or `draft` - so should be the architecture.

Components:

- `joker` - client that you install locally
- `gotham` - server-side component that lives inside your cluster and has pods you can configure to deploy cloud resources

`joker init` - deploys `gotham` to your cluster

`joker configure {cloud}` - uses the native authentication of the cloud provider to configure secrets with cloud credentials (client IDs, client secrets, tenant IDs, keys)

`joker create {cloud} resource` - creates a cloud service and configures a Kubernetes secret with the access keys

> Example: `joker create az redis --name {} --location {} --resource-group {} --sku {} --shard-count {}`

`joker ls {cloud}` - list all {cloud} services

`joker delete {cloud} {resource}` - delete a {cloud} service


> {cloud} - aws/az/gcloud

> This is where the discussion needs to take place on how to pass parameters for your service - they might be the parameters you pass to `aws`, `az` and `gcloud`, it can be a Terraform file, ARM template or other deployment mechanism - **initially this will support the command line arguments you pass to `aws`, `az`, `gcloud`**.

> Another discussion is declarative deployments


How about [Kubernetes Service Brokers](https://github.com/kubernetes-incubator/service-catalog)?
-------------------------------------

> Watch this [introductory video on Service Brokers and Service Catalog](https://www.youtube.com/watch?v=0aLqc-o256w&app=desktop)

> A service broker is an endpoint that manages a set of software offerings called services. The end-goal of the service-catalog project is to provide a way for Kubernetes users to consume services from brokers and easily configure their applications to use those services, without needing detailed knowledge about how those services are created or managed.


While the service brokers extend the `kubectl` functionality, and seem to be more generic (and extensible to any sort of service), `joker` limits itself to working with Azure, GCP and AWS and making use of services from these cloud providers as easy as possible.

Service Catalog is an awesome feature (that I am glad I discovered), but at this moment it seems like you need to understand a lot of concepts before integrating with Azure, GCP and AWS - and it is a lot more powerful than what `joker` wants to be.

Why the name?
-------------

In some card games, the Joker is the card you can use to your advantage and can substitute other cards - continuing our analogy - if you happen to use Google Cloud services with Kubernetes, use Joker to deploy and manage them; if you use Azure services, use Joker to deploy and manage them - you get the point?

Of course, the Joker is one of my favourite characters, so there's that.

> Now having decided on this name, I'm hoping at some point [@ashleymcnamara](https://github.com/ashleymcnamara/gophers) will be kind enough to draw an awesome logo (as all her others!) for Joker. Until then, I will be using [her Batman gopher (you get it, Batman, Joker :D)](https://github.com/ashleymcnamara/gophers/blob/master/BATMAN_GOPHER.png) - here's [the license from her gophers artwork](https://github.com/ashleymcnamara/gophers/blob/master/LICENSE)

Initial thoughts
----------------

- this tool has to be cloud agnostic - you want to use it regardless of what services you use and from what cloud provider
- it doesn't need to reinvent the wheel in deploying cloud resources - there needs to be a discussion on what will be supported - various cloud CLIs (`az`, `gcloud`, `aws`), SDKs, native deployment files, Terraform - there are tons of ways of deploying cloud services. **This tool should support the ones people use and not reinvent new ones!**


Considerations
--------------

- while I work for Microsoft, **this is not an official Microsoft project**
- this project is highly experimental (at the time of writing this document, not a single line of code was written)
- this project wants to explore the other side of developing Kubernetes apps with a cloud provider - so far, the discussion has been: "how do we deploy Kubernetes itself on this cloud provider?" rather than: "how do we use services from this cloud provider together with Kubernetes in the most convenient way?"

- some people will argue that this is way out of scope for Kubernetes - but I would say that Kubernetes already make use of great platform features (load balancers, storage) - why not expand that to databases, caches and cool platform services people want to use with the cloud?

- if at any point you have any idea regarding this project, your contribution is most welcome (as long as it is respectful and mindful of others).

Documentation, installing, using
--------------------------------

As the project is in very early stages, things will change pretty fast - so if you see the docs out of data, please ping or create a pull request. 

There will be docs, but since there is nothing to document yet, this is just a placeholder - thanks for your patience :D

Contributing
------------

Any idea (here, on Twitter - @Matei_Radu), issue or pull request is highly appreciated. Contribution guidelines will follow once there is a structure to this project.