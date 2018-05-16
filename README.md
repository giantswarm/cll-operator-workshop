# cll-operator-workshop

This a repo pepared for https://continuouslifecycle.london/ operator workshop.
It demonstrates how to create a basic operator using [operatorkit][operatorkit]
and generated CRD client.

## Why Memcached?

[Memcached][memcached] was picked as an example as this shows how the operator
can be built to create and scale a distributed system running inside
Kubernetes. Memcached is easy to scale in particular as sharding is done on the
client side. This allows us to show a complete example which does not do
complex logic not really interesting from the example perspective.

## Implementation

This operator controller watches MemcachedConfig CRD (Custom Resource
Definition). Upon arrival of new CR (custom resource, also called object) it
creates the Memcached cluster according to the CR spec.

E.g. let's consider following MemcachedConfig CR:

```yaml
apiVersion: workshop.continuouslifecycle.london/v1alpha1
kind: MemcachedConfig
metadata:
  name: mycluster
  namespace: default
spec:
  replicas: 3
  memory: 100Mi
```

For this CR the operator will create Memcached cluster of 3 nodes 1G memory
each in the same namespace as the CR is defined in.

When `replicas` part of the spec incremented a new cluster node is added. When
`replicas` is decremented the youngest node is removed from the cluster.
When `memory` is changed all memcached nodes are rolled (not best design, but
good enough for workshop purposes).

For each cluster node the operator creates a corresponding Kubernetes Service.
Service names are the name of the CR plus index of the node staring from zero.
E.g. for CR with `name: mycluster` and `replicas: 3` there will be 3 Services
created with names `mycluster0000`, `mycluster0001`, and `mycluster0002`. This
is predictable and allows to connect easily to all nodes by the client with
knowledge of name and number of replicas. E.g. for [gomemcache][gomemcache]
connecting code in this case can look like:

```go
memcache.New("mycluster0000:11211", "mycluster0001:11211", "mycluster0002:11212")
```

[gomemcache]: https://github.com/bradfitz/gomemcache
[memcached]: https://memcached.org/
[operatorkit]: https://github.com/giantswarm/operatorkit
