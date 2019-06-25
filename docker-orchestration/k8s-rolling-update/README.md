# Kubernetes Rolling Update Deployment
Use kubernetes rolling update deployment strategy for timed replace container with new one. Latest traffic will be used during container startup.  

## Usage
```bash

# NOTE: prepare your image and fill into osrm.yaml
$ sed -i "" "s#TELENAV_OSRM_BACKEND_DOCKER_IMAGE#<your image>#g" ./osrm.yaml

# create deployment
$ kubectl create -f osrm.yaml  
deployment.extensions/osrm created

# OPTIONAL: checking deployment status, wait READY
$ kubectl get deployments 
NAME   READY   UP-TO-DATE   AVAILABLE   AGE
osrm   0/1     1            0           20s
$ kubectl get pods
NAME                    READY   STATUS    RESTARTS   AGE
osrm-67b55f46c9-tpfgs   0/1     Running   0          22s

# OPTIONAL: check deployment details for troubleshooting
$ kubectl describe deployments
...
$ kubectl describe pods
...

# OPTIONAL: check container running log for troubleshooting
$ kubectl logs -f --timestamps=true osrm-67b55f46c9-tpfgs
...

# start timed rolling update for traffic updating
$ crontab osrm_cron
$ crontab -l 
*/30 * * * * kubectl set env --env="LAST_MANUAL_RESTART=$(date -u +\%Y\%m\%dT\%H\%M\%S)" deploy/osrm >> ${HOME}/osrm_cron.log 2>&1


# OPTIONAL: if want to stop
$ crontab -r 
$ kubectl delete deployment osrm
deployment.extensions "osrm" deleted

```