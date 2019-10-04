---
title: "Motivation"
date: 2019-09-30T17:20:32+02:00
weight: 1
---

#### Opinions Behind Pligos

The following article is going to describe my opinions on why Pligos is the superior Kubernetes manifest management tool and why Pligos is not even trying to be a manifest management tool. Managing Kubernetes manifests is hard and has many flaws, Joe Beda one of the original authors agrees with me and probably the whole community.

<div style="text-align: right">  

<blockquote class="twitter-tweet" data-lang="en-gb"><p lang="en" dir="ltr">As <a href="https://twitter.com/bryanl?ref_src=twsrc%5Etfw">@bryanl</a> says: YAML is for computers. When we started with YAML we never intended it to be the user facing solution. We saw it as &quot;assembly code&quot;. I&#39;m horrified that we are still interacting with it directly. That is a failure.</p>&mdash; Joe Beda (@jbeda) <a href="https://twitter.com/jbeda/status/994566252503810048?ref_src=twsrc%5Etfw">10 May 2018</a></blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
 
 
</div>


<https://twitter.com/jbeda/status/994566252503810048>

That's why an abundance of tools have been published to mitigate this problem. However I believe, while there are tools, such as `draft.sh` that are fairly close, there is no solution which fixes this issue in it's entireness. So let me explain:

In my opinion, with tools such as `helm`, `ksonnet` and `kapitan` we are inventing a lot of clever ways to make Kubernetes manifest templating bearable. While I agree, that some of the solutions (such as using jsonnet over Go templating) are more effective than others we are not addressing the real issue. Kubernetes manifests are a way to describe Kubernetes resources. While we are forced to describe Kubernetes resources, what we really want to do is describe configuration for our services. And because of that we are doing the same task over and over again, which is mapping everything our service requires to what the Kubernetes API offers. We often invent clever little strategies to do just that and often times, while defining the configuration for one service, we realize how everything could have been done so much better for the last few services. As a result of this we either let our configuration style diverge, or we take the time and refactor. Neither options are optimal and above all unnecessary.

The only way out of this mess is to **disconnect infrastructure configuration** from **Kubernetes manifest management**. What do I mean by that? When I am talking about infrastructure configuration I have something like this in mind:

    route:
      - name: http
        applicationPort: 8080
        exposeAs: 80
        protocol: TCP
        path: /myservice/api/v1

This configuration describes exactly how we want our service to behave once it's deployed inside a Kubernetes environment. We are stating that we have some service, which is going to listen on port 8080, while we want it to be exposed on port 80. Also, we want it to be accessible on the path <https://my.cluster/myservice/api/v1>. This configuration is universally true, regardless of how we map it to the Kubernetes API. In order to close the gap between what we defined in this little example and how our kubernetes manifests are defined we could use any of the tools already out there, the leading examples being Helm or ksonnet. The important part is, that you create as few Helm, ksonnet, whatever configuration instances as possible, such that you don't repeat yourself and can keep the refactoring to a minimum. In my experience, you could even stick to three configuration instances: `stateless web services`, `stateful applications` and `per node daemons`.

Pligos is currently implemented to use helm. It is based on a concept called **flavors**, which are basically helm starters with a public API, called **schema** which describes the structure in which your application should be configured. One example of such a configuration is shown above already. This allows the mapping of multiple services to a single helm chart. Once you configured your service based on the schema definition Pligos goes on and compiles the input configuration into something that maps to the configuration template. Additionally,  Pligos also treats the management of different contexts (`development`, `production`, `ci`) as a first class citizen. In conjunction with this, it's also planned to support `secret management` (TLS certificates, docker registry secrets, &#x2026; ) and application configuration management. Pligos already is designed with application configuration management in mind, however this is a whole topic on it's own.

When Joe Beda is talking about the horrible state we are in regarding Kubernetes yaml configurations I think people put too much focus on yaml being the issue. Yes, yaml might not be the final answer, however, I think we need to fight different issues. Pligos is here to minimize  the handling of Kubernetes manifests and instead tries to provide a configuration interface that let's us define behavior of our configuration in the cloud, not Kubernetes manifests.


#### Difference of Helm Directory Structure

Difference between the directory structure of helm explains well that how pligos made the life easier to manage the cloud infrastructure.

In the following image, you can see on `left side` that each microservice has its own deployment.yaml, service.yaml, ingress.yaml and etc. So, it will be really hectic to manage numerous services while using this usual approach.

And on the `right side`, it shows that by using pligos, you don't need to worry about the yaml files for configuration of each micro-service. You just have to configure the pligos.yaml for each service.
![]({{< resource url="with-without-plilgos.jpg" >}})
