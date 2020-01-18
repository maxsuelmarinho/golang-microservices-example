package com.marinho.microservices.turbineamqpplugin;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.netflix.turbine.discovery.StreamAction;
import com.rabbitmq.client.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.env.Environment;
import org.springframework.integration.amqp.dsl.Amqp;
import rx.Observable;
import rx.subjects.PublishSubject;

import java.io.IOException;
import java.net.URISyntaxException;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.util.concurrent.TimeoutException;

public class AmqpInstanceDiscovery {

    private static final Logger logger  = LoggerFactory.getLogger(AmqpInstanceDiscovery.class);

    private String amqpBrokerUrl;
    private String discoveryQueue;
    private String clusterName;
    private String consumerTag;
    private PublishSubject<AmqpInstance> subject;
    private final Environment environment;

    AmqpInstanceDiscovery(final Environment environment) {
        this.environment = environment;
        configureFromEnv();

        subject = PublishSubject.create();
        try {
            ConnectionFactory factory = new ConnectionFactory();
            factory.setUri(amqpBrokerUrl);
            Connection conn = factory.newConnection();
            Channel channel = conn.createChannel();
            boolean autoAck = true;
            channel.basicConsume(discoveryQueue, autoAck, consumerTag, new DefaultConsumer(channel) {
               @Override
               public void handleDelivery(String consumerTag,
                                          Envelope envelope,
                                          AMQP.BasicProperties properties,
                                          byte[] body)
                       throws IOException {
                   DiscoveryToken token = new ObjectMapper().readValue(body, DiscoveryToken.class);
                   AmqpInstance.Status status = AmqpInstance.Status.valueOf(token.getState());
                   AmqpInstance amqpInstance = new AmqpInstance(clusterName, status, token.getAddress(), 8181);
                   logger.debug("Got token for instance: " + new String(body));
                   if (subject != null) {
                       subject.onNext(amqpInstance);
                   }
               }
            });
        } catch (NoSuchAlgorithmException | KeyManagementException | URISyntaxException | TimeoutException | IOException e) {
            logger.error("Problem connecting to RabbitMQ: " + e.getMessage(), e);
            logger.error("Sleeping for five seconds before terminating");
            try {
                Thread.sleep(5000L);
            } catch (InterruptedException ex) {
                ex.printStackTrace();
            }
            logger.error("System exit!");
            System.exit(0);
        }
    }

    private void configureFromEnv() {
        amqpBrokerUrl = String.format("amqp://%s:%s@%s:%s",
                environment.getProperty("spring.rabbitmq.username"),
                environment.getProperty("spring.rabbitmq.password"),
                environment.getProperty("spring.rabbitmq.host"),
                environment.getProperty("spring.rabbitmq.port"));

        discoveryQueue = environment.getProperty("turbineamqpplugin.discovery-queue");
        clusterName = environment.getProperty("turbineamqpplugin.cluster-name");
        consumerTag = environment.getProperty("turbineamqpplugin.consumer-tag");
    }

    Observable<AmqpInstance> getInstanceEvents() {
        subject.subscribe(System.out::println);
        return subject;
    }
}
