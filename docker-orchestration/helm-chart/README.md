# Kubernetes deployment by Helm Chart
Use [Helm Chart](https://helm.sh/docs/intro/using_helm/) for `osrm-backend` and `osrm-frontend` deployment on kubernetes cluster. These experiences are based on [AWS EKS](https://aws.amazon.com/eks), but should be similar on other platform.         

## Usage

```bash
$ helm version
version.BuildInfo{Version:"v3.0.2", GitCommit:"19e47ee3283ae98139d98460de796c1be1e3975f", GitTreeState:"clean", GoVersion:"go1.13.5"}
$ 
$ cd docker-orchestration/helm-chart
$ 
$ helm list -n routing-osrm
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
$
$ # install osrm-backend with specified within-data image
$ helm install osrm-backend --generate-name -n routing-osrm --set osrm.backend.image=telenavmap/osrm-backend-within-mapdata:no.86-20200114-master-telenav-e85d5ca-compile28-20200119T171422CST-nevada-latest
NAME: osrm-backend-1579678916
LAST DEPLOYED: Tue Jan 21 23:41:56 2020
NAMESPACE: routing-osrm
STATUS: deployed
REVISION: 1
TEST SUITE: None
$ 
$ # see the deployed osrm-backend
$ helm list -n routing-osrm
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
osrm-backend-1579678916 routing-osrm    1               2020-01-21 23:41:56.694024666 -0800 PST deployed        osrm-backend-0.4.0
$ 
$ # get service for this osrm-backend installation
$ kubectl get svc -n routing-osrm | grep osrm-backend-1579678916
osrm-backend-1579678916   LoadBalancer   172.20.236.194   internal-aaa74ad603cea11ea82ff0268b6e6a1a-222618767.us-west-2.elb.amazonaws.com    5000:30346/TCP                  6m46s
$ 
$ # install osrm-frontend target to this backend
$ helm install osrm-frontend --generate-name -n routing-osrm --set osrm.frontend.targetBackend='http://internal-aaa74ad603cea11ea82ff0268b6e6a1a-222618767.us-west-2.elb.amazonaws.com:5000',osrm.frontend.center='36.122314\,-115.111110'
$ 
$ # see the deployed osrm-frontend
$ helm list -n routing-osrm
NAME                            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
osrm-backend-1579678916         routing-osrm    1               2020-01-21 23:41:56.694024666 -0800 PST deployed        osrm-backend-0.4.0
osrm-frontend-1579679580        routing-osrm    1               2020-01-21 23:53:00.988146482 -0800 PST deployed        osrm-frontend-0.4.0
$ 
$ # get service for this osrm-frontend installation
$ kubectl get svc -n routing-osrm | grep osrm-frontend-1579679580
osrm-frontend-1579679580   LoadBalancer   172.20.40.115    internal-a36697e403cec11eab7a606e7ec8f570-2062152117.us-west-2.elb.amazonaws.com   9966:30997/TCP                  70s
$ 
$ # Now it's possible to access both frontend and backend via these installed services
$ curl "http://internal-a36697e403cec11eab7a606e7ec8f570-2062152117.us-west-2.elb.amazonaws.com:9966"
$ curl "http://internal-aaa74ad603cea11ea82ff0268b6e6a1a-222618767.us-west-2.elb.amazonaws.com:5000"
$ 
$ # uninstall after use
$ helm uninstall osrm-frontend-1579679580 -n routing-osrm
$ helm uninstall osrm-backend-1579678916 -n routing-osrm
```