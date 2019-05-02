package com.marinho.microservices.turbineamqpplugin;

import com.netflix.turbine.discovery.StreamAction;
import com.netflix.turbine.discovery.StreamDiscovery;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;
import rx.Observable;

import java.net.URI;
import java.net.URISyntaxException;

@Component
public class AmqpStreamDiscovery implements StreamDiscovery {

    private static final Logger logger = LoggerFactory.getLogger(AmqpStreamDiscovery.class);

    final static String HOSTNAME = "{HOSTNAME}";

    private final String uriTemplate;
    private final Environment environment;

    AmqpStreamDiscovery(Environment environment) {
        this.environment = environment;
        this.uriTemplate = environment.getProperty("turbineamqpplugin.uri-template");
    }

    @Override
    public Observable<StreamAction> getInstanceList() {
        return new AmqpInstanceDiscovery(environment)
                .getInstanceEvents()
                .map(e -> {
                    URI uri;
                    String uriString = uriTemplate.replace(HOSTNAME, e.getHostname() + ":" + e.getPort());
                    logger.debug("uriString: " + uriString);
                    try {
                        uri = new URI(uriString);
                    } catch (URISyntaxException ex) {
                        throw new RuntimeException("Invalid URI: " + uriString, ex);
                    }

                    if (e.getStatus() == AmqpInstance.Status.UP) {
                        logger.debug("StreamAction ADD");
                        return StreamAction.create(StreamAction.ActionType.ADD, uri);
                    } else {
                        logger.debug("StreamAction REMOVE");
                        return StreamAction.create(StreamAction.ActionType.REMOVE, uri);
                    }
                });
    }
}