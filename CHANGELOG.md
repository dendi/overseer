# Changelog

## [2019/12/11] cmaster11/overseer:1.12.1

* More reasonable tag names
* You can now perform [period-tests](./README.md#period-tests):

> What if you want to test how many times your web service fails in 1 minute?

## [2019/12/10] cmaster11/overseer:release-1.12

* Overseer runs now by default multiple tests in parallel (defaults to num of CPUs). This behavior is tunable with the cli flag `-parallel`:

```
overseer worker -parallel 9
```

## [2019/10/31] cmaster11/overseer:release-1.11-5

* Added support to accept multiple status codes when performing HTTP checks:

```
https://www.google.com must run http with status 200,301
```

## [2019/10/26] cmaster11/overseer:release-1.11-4

* Added support for queue-bridge filters. Destination queues (`-dest-queue=overseer.results.email`) can now be filtered
by various test result tags: `-dest-queue=overseer.results.cronjobs[type=k8s-event,target=my-namespace/Job]`. More in 
the [Kubernetes example](example-kubernetes/overseer-bridge-queue.optional.yaml) and 
[source code](bridges/queue-bridge/filter.go).