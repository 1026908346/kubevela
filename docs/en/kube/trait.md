# Attach Traits to Kube Based Components

Most traits in KubeVela can be attached to Kube based component seamlessly. 
In this sample application below, we add two traits,
[scaler](https://github.com/oam-dev/kubevela/blob/master/charts/vela-core/templates/defwithtemplate/manualscale.yaml)
and
[virtualgroup](https://github.com/oam-dev/kubevela/blob/master/docs/examples/kube-module/virtual-group-td.yaml), to a Kube based component.

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: myapp
  namespace: default
spec:
  components:
    - name: mycomp
      type: kube-worker
      properties: 
        image: nginx:1.14.0
      traits:
        - type: scaler
          properties:
            replicas: 2
        - type: virtualgroup
          properties:
            group: "my-group1"
            type: "cluster"
```

## Verify traits work correctly

Deploy the application and verify traits work.

Check the scaler trait.
```shell
$ kubectl get manualscalertrait
NAME                            AGE
demo-podinfo-scaler-3x1sfcd34   2m
```
```shell
$ kubectl get deployment mycomp-v1 -o json | jq .spec.replicas
2
```

Check the virtualgroup trait.
```shell
$ kubectl get deployment mycomp-v1 -o json | jq .spec.template.metadata.labels
{
  "app.cluster.virtual.group": "my-group1",
  "app.kubernetes.io/name": "myapp"
}
```

## Update an Application

After the application is deployed and workloads/traits are created successfully,
you can update the application, and corresponding changes will be applied to the
workload.

Let's make several changes on the configuration of the sample application.

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: myapp
  namespace: default
spec:
  components:
    - name: mycomp
      type: kube-worker
      properties: 
        image: nginx:1.14.1 # 1.14.0 => 1.14.1
      traits:
        - type: scaler
          properties:
            replicas: 4 # 2 => 4
        - type: virtualgroup
          properties:
            group: "my-group2" # my-group1 => my-group2
            type: "cluster"
```

Apply the new configuration and check the results after several seconds.

> After updating, the workload name is changed from `mycomp-v1` to `mycomp-v2`.

Check the new parameter works.
```shell
$ kubectl get deployment mycomp-v2 -o json | jq '.spec.template.spec.containers[0].image'
"nginx:1.14.1"
```

Check the scaler trait.
```shell
$ kubectl get deployment mycomp-v2 -o json | jq .spec.replicas
4
```

Check the virtualgroup trait.
```shell
$ kubectl get deployment mycomp-v2 -o json | jq .spec.template.metadata.labels
{
  "app.cluster.virtual.group": "my-group2",
  "app.kubernetes.io/name": "myapp"
}
```