package com.marinho.microservices.turbineamqpplugin;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

public class AmqpInstance {

    public enum Status {
        UP, DOWN
    }

    private final String cluster;
    private final Status status;
    private final String hostname;
    private final int port;
    private final Map<String, Object> attributes;

    public AmqpInstance(String cluster, Status status, String hostname, int port) {
        this.cluster = cluster;
        this.status = status;
        this.hostname = hostname;
        this.port = port;
        this.attributes = new HashMap<>();
    }

    public String getCluster() {
        return cluster;
    }

    public Status getStatus() {
        return status;
    }

    public String getHostname() {
        return hostname;
    }

    public int getPort() {
        return port;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        AmqpInstance that = (AmqpInstance) o;
        return port == that.port &&
                Objects.equals(cluster, that.cluster) &&
                status == that.status &&
                Objects.equals(hostname, that.hostname) &&
                Objects.equals(attributes, that.attributes);
    }

    @Override
    public int hashCode() {
        return Objects.hash(cluster, status, hostname, port, attributes);
    }

    @Override
    public String toString() {
        return "AmqpInstance{" +
                "cluster='" + cluster + '\'' +
                ", status=" + status +
                ", hostname='" + hostname + '\'' +
                ", port=" + port +
                ", attributes=" + attributes +
                '}';
    }
}
