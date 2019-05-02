#!/bin/bash

if [ -z "$HEAP_OPTS" ]; then
    HEAP_OPTS="-Xms128M -Xmx256M"
fi

if [ -z "$JVM_PERFORMANCE_OPTS" ]; then
    JVM_PERFORMANCE_OPTS="-server -XX:+UseG1GC -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35 -XX:+DisableExplicitGC -Djava.awt.headless=true"
fi

cd /app/ && java ${HEAP_OPTS} ${JVM_PERFORMANCE_OPTS} -Djava.security.egd=file:/dev/./urandom -jar /app/app.jar
