package com.marinho.microservices.turbineamqpplugin;

import com.netflix.turbine.Turbine;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.netflix.eureka.EnableEurekaClient;
import org.springframework.cloud.netflix.turbine.stream.EnableTurbineStream;
import org.springframework.context.ConfigurableApplicationContext;

//@EnableEurekaClient
@EnableTurbineStream
@SpringBootApplication
public class TurbineAmqpPluginApplication {

    private static final Logger LOG = LoggerFactory.getLogger(TurbineAmqpPluginApplication.class);

    public static void main(String[] args) {
        ConfigurableApplicationContext ctx = SpringApplication.run(TurbineAmqpPluginApplication.class, args);

        LOG.info("Connected to RabbitMQ at: {}", ctx.getEnvironment().getProperty("spring.rabbitmq.host"));
        Turbine.startServerSentEventServer(8282, ctx.getBean(AmqpStreamDiscovery.class));
    }

}