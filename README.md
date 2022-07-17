# FanController

Fan Motor: https://www.amazon.de/dp/B0922N8MCR

Docker RabbitMQ: ```docker run -d -p 5672:5672 --hostname fan-rabbit --name fan-rabbit rabbitmq:3```

## Architecture Diagram:
Principle: using messages queues and exchanges to send data around.

Using a smart broker to route messages to the intended receivers (using rabbitMQ's Exchanges and Topics)

![Architecture Diagram](./documentation/FanController.png)

## Links:
* RabbitMQ: [https://www.rabbitmq.com](https://www.rabbitmq.com)