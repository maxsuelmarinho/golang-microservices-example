package com.marinho.microservices.turbineamqpplugin;

public class DiscoveryToken {

    private String state;
    private String address;

    public DiscoveryToken() {
    }

    public DiscoveryToken(String state, String address) {
        this.state = state;
        this.address = address;
    }

    public String getState() {
        return state;
    }

    public void setState(String state) {
        this.state = state;
    }

    public String getAddress() {
        return address;
    }

    public void setAddress(String address) {
        this.address = address;
    }
}
