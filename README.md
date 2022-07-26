# FanController

Fan Motor: https://www.amazon.de/dp/B0922N8MCR

Docker RabbitMQ: ```docker run -d -p 5672:5672 --hostname fan-rabbit --name fan-rabbit rabbitmq:3```

## Architecture Diagram:
Principle: using messages queues and exchanges to send data around.

Using a smart broker to route messages to the intended receivers (using RabbitMQ's Exchanges and Topics)

![Architecture Diagram](./documentation/FanController.png)

The Transformation logic is stored in a separate config file: it lets users define the threshold value, comparator and result values as well as the targets for a given transformation. This flexibility moves the config outside the system and makes having an external config system possible.

## Links:
* RabbitMQ: [https://www.rabbitmq.com](https://www.rabbitmq.com)