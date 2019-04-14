## Load test

```
mvn -N io.takari:maven:wrapper
mvnw gatling:execute -Dusers=1000 -Dduration=30 -DbaseUrl=http://localhost:8080
```

```
> ps -eo size,pid,ppid,cmd,%mem,%cpu --sort=-%mem | head

### 1 ###
# Idle process
SIZE   PID  PPID CMD                         %MEM %CPU
   0   280     4 go run main.go               0.1  0.2
   0   418   280 /tmp/go-build833902426/b001  0.0  0.0

# Under load test
SIZE   PID  PPID CMD                         %MEM %CPU
   0   280     4 go run main.go               0.1  0.2
   0   418   280 /tmp/go-build833902426/b001  0.3  1.1

### 2 ####
# Idle process
SIZE   PID  PPID CMD                         %MEM %CPU
   0   869     4 go run main.go               0.1  0.9
   0   996   869 /tmp/go-build904335587/b001  0.0  0.0

# Under load test
SIZE   PID  PPID CMD                         %MEM %CPU
   0   869     4 go run main.go               0.1  0.7
   0   996   869 /tmp/go-build904335587/b001  0.4  6.2

### 3 ###
# service containerized
> docker stats $(docker ps | awk '{if(NR>1) print $NF}')
# Idle process
CONTAINER ID        NAME                                         CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
076f1a57b097        accountservice.1.q2xop6g2lcyr45hx86evyffe4   0.00%               6.328MiB / 2.771GiB   0.22%               4.49kB / 3.53kB     0B / 1.72MB         5

# Under load test
CONTAINER ID        NAME                                         CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
076f1a57b097        accountservice.1.q2xop6g2lcyr45hx86evyffe4   25.41%              33.26MiB / 2.771GiB   1.17%               10.5MB / 10.5MB     0B / 1.72MB         8

### 4 ###
# docker swarm mode - services scaled up
CONTAINER ID        NAME                                         CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O           PIDS
fc7f07a7a19b        accountservice.4.nu7vc6x1b40u9koxoni7kyr3a   75.61%              15.27MiB / 2.771GiB   0.54%               2.41MB / 3.12MB     0B / 1.79MB         38
28b82b7d2cd0        accountservice.3.ktfyk3p02zd65cby8u6ic0n47   64.26%              15.31MiB / 2.771GiB   0.54%               2.41MB / 3.13MB     0B / 1.8MB          41
21ecdeca6614        accountservice.2.8yooxp89wxuwiqrrqsv1uyffe   60.41%              15.28MiB / 2.771GiB   0.54%               2.42MB / 3.13MB     0B / 1.79MB         47
2e062d0ff429        accountservice.1.dqjxrkg9o760imjglnw1ro7yn   61.77%              16.86MiB / 2.771GiB   0.59%               2.44MB / 3.17MB     0B / 1.8MB          31

```

## Docker Swarm

**Configure Manager Node**

```
> docker swarm init --advertise-addr 192.168.99.100
Swarm initialized: current node (eltsfe59whab9d71bjsyeflpa) is now a manager.

To add a worker to this swarm, run the following command:

    docker swarm join --token SWMTKN-1-0ljalif9jf9grmlm1hbewplgj44wh6t2wt44n1wfrrxypxlo7o-55jcgul49od8nuvloo9vu9zs2 192.168.99.100:2377

To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
```

**Current State of the Swarm**

```
docker info
```

**Information about nodes**

```
> docker node ls

ID                            HOSTNAME                STATUS              AVAILABILITY        MANAGER STATUS      ENGINE VERSION
eltsfe59whab9d71bjsyeflpa *   localhost.localdomain   Ready               Active              Leader              18.06.3-ce

```

**Retrieve the join command for a worker**

```
docker swarm join-token worker
To add a worker to this swarm, run the following command:

    docker swarm join --token SWMTKN-1-0ljalif9jf9grmlm1hbewplgj44wh6t2wt44n1wfrrxypxlo7o-55jcgul49od8nuvloo9vu9zs2 192.168.99.100:2377
```

**Create an overlay network**

```
docker network create --driver overlay my_network
```

**Deploying the Account Service**

```
docker service create --name=account-service --replicas=1 --network=my_network -p=8080:8080 maxsuelmarinho/golang-microservices-example:accountservice-0.0.1

docker service create --name=quotes-service --replicas=1 --network=my_network -p=9090:8080 maxsuelmarinho/golang-microservices-example:quotesservice-0.0.1
```

**Service status**

```
> docker service ls

ID                  NAME                MODE                REPLICAS            IMAGE                                                              PORTS
sc9ofnr14ox6        account-service      replicated          1/1                 maxsuelmarinho/golang-microservices-example:accountservice-0.0.1   *:8080->8080/tcp
```

**Remove a service**

```
docker service rm account-service
```

**Scale up**

```
docker service scale account-service=4
```

**Deploy the stack to the swarm**

```
docker stack deploy --compose-file docker-compose.yml stackdemo
```

**Stack status**

```
> docker stack services stackdemo

ID                  NAME                        MODE                REPLICAS            IMAGE                                                              PORTS
d6hm9vrwl8gx        stackdemo_quotes-service    replicated          1/1                 maxsuelmarinho/golang-microservices-example:quotesservice-0.0.1    *:9090->8080/tcp
fq684muqeyel        stackdemo_viz               replicated          1/1                 dockersamples/visualizer:latest                                    *:8000->8080/tcp
k1wffx9o1583        stackdemo_account-service   replicated          1/1                 maxsuelmarinho/golang-microservices-example:accountservice-0.0.1   *:8080->8080/tcp
```

**Remove stack**

```
docker stack rm stackdemo
```


**Visualizers**

```
docker service create --name=viz --publish=8000:8080/tcp --constraint=node.role==manager --mount=type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock dockersamples/visualizer
```

```
docker service create --constraint=node.role==manager --replicas 1 --name dvizz -p 6969:6969 --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --network my_network eriklupander/dvizz
```


## Account Service API

**Get Account**

```
> curl http://localhost:8080/accounts/10000 | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    32  100    32    0     0   1882      0 --:--:-- --:--:-- --:--:--  1882
{
  "id": "10000",
  "name": "Person_0"
}
```
