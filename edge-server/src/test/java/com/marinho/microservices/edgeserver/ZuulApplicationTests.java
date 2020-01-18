package com.marinho.microservices.edgeserver;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

@RunWith(SpringRunner.class)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT, properties = {
        "spring.cloud.bus.enabled=false",
        "spring.sleuth.stream.enabled=false"
})
public class ZuulApplicationTests {

    @Test
    public void contextLoads() {}
}