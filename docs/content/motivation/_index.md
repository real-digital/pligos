---
title: "Motivation"
date: 2019-09-30T17:20:32+02:00
weight: 2
---

##### Author: Abdul Jabbar
##### Created at: 19 Sep 2019


Following difference between the structure explains well that how pligos made the life easier to manage the cloud infrastructure. So, this is the motivation behind using pligos. 

#### Helm Directory Structure without using pligos

You can see in the following images that each microservice has its own deployment, service, ingress and etc. So, it will be really hectic to manage numerous services while using this usual approach.
![]({{< resource url="helm-structure-without-pligos.jpg" >}})

#### Helm Directory Structure with using pligos

And this image shows that by using pligos, you don't need to worry about the yaml files. You just have to configure the pligos.yaml for each service.
![]({{< resource url="helm-structure-with-pligos.jpg" >}})

**Further motivation behind pligos is given [here](https://github.com/real-digital/pligos/wiki/Opinions-behind-Pligos)**