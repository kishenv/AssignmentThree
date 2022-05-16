##  Asignment 3
##### 1. Create a Client-Server service that reports heartbeats in intervals of 3 seconds, the CPU and memory usage profiles every 20 seconds.

1. Set up a HTTP listener on the Server node that listens to the heartbeats and the compute metrics of the client node.
2. Setup a client service that reports the metrics and the heartbeats to the Server node through REST API calls.
3. Build the required docker images for both, client and server and deploy in Minikube environment.
4. Configure services for the deployments to allow exchange of data between the pods.



##### Commands to deploy the Client-Server pods to exhange Data.

```bash
kubectl apply -f metrics-yaml/
kubectl logs -f <pod-name>
kubectl delete -f metrics-yaml/
```
Validation environment:
* minikube : v1.25.2
* docker : 20.10.12
* kubectl :v1.24.0
* Ubuntu : 22.04

**Apply the configuration and monitor the status of deployments and services.:**
[![DeployClientServer](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/Deploy.png?raw=true "DeployClientServer")](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/Deploy.png?raw=true"DeployClientServer")

**Validate the service functionality from logs:**
The exchange of data with regards to the heartbeat and metrics that are sent by the client are acknowledged by the server with a HTTP 200 OK and a response.

**Logs from client**
[![clientLogs](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/ClientLog.png?raw=true "clientLogs")](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/ClientLog.png?raw=true "clientLogs")

**Logs from server**
[![Server Logs ](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/Server.png?raw=true "Server Logs ")](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/Server.png?raw=true "Server Logs ")

** Behavior of the server on termination of the client pod by deleting  deployment.**

On deleting the metrics-client deployment, the server rolls back to reporting the lost connection once every 10 seconds, after 3 consecutive heartbeats are not received by the server.

[![Disconnection](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/LostConnection.png?raw=true "Disconnection")](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/LostConnection.png?raw=true "Disconnection")

** Behavior of server on re-deploying the client pod.**
On re-deploying the metrics-client deployment, the heartbeats are received by the server with the metrics re-captured.

[![DisconnRecon](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/DisconnectAndReconnect.png?raw=true "DisconnRecon")](https://github.com/kishenv/AssignmentThree/blob/main/Screenshots/DisconnectAndReconnect.png?raw=true "DisconnRecon")
