# What is Krossboard
Krossboard provides an advanced centralized Kubernetes resource usage analytics and accounting for multiple Kubernetes clusters. The clusters can be located on premises and/or in the cloud (i.e. self-deployed or managed).

Key features:

* **Multi-Kubernetes Data Collection**: Krossboard periodically collects raw metrics related to containers, pods and nodes from each Kubernetes cluster it handles. The built-in data collection period is 5 minutes.
* **Powerful Analytics Processing**: Krossboard internally processes raw metrics to produce insightful Kubernetes usage accounting and analytics metrics. The analytics data are tracked on a hourly-basis, per namespace, per cluster, and globally.
* **Insightful Usage Accounting**: Krossboard periodically processes usage accounting, per namespace and per cluster. By the default, the UI displays accounting the following periods without any additioanl configuration: daily accounting for the last 14 days, monthly for the ast 12 months.
* **Easy to getting started**: Krossboard concepts are intuitive, the Kubernetes clusters it manages are set through standard KUBECONFIG resources. Thanks to its Kubernetes operator, Krossboard is easy to install and to operate. 
* **REST API to download reports**: Krossboard enables REST API to expose the analytics and accounting data it generates to third-parties analytics systems in CSV or JSON format.


![](../krossboard-architecture-overview.png)


# Getting Started

Use our ready-to-use [Krossboard Kubernetes Operator](https://github.com/2-alchemists/krossboard-kubernetes-operator) to easily setup and operate your instance of Krossboard.


# Additional Resources

* [Krossboard Website](https://krossboard.app/)
* [Krossboard Enterprise Support](https://krossboard.app/#pricing) 
