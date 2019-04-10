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
```