# FanController

Fan Motor: https://www.amazon.de/dp/B0922N8MCR

Docker RabbitMQ: ```docker run -d -p 5672:5672 --hostname fan-rabbit --name fan-rabbit rabbitmq:3```

## Architecture Diagram:
Principle: using messages queues and exchanges to send data around.

Using a smart broker to route messages to the intended receivers (using RabbitMQ's Exchanges and Topics)

![Architecture Diagram](./documentation/FanController.png)

The Transformation logic is stored in a separate config file: it lets users define the transformation function by using a custom AST syntax as well as the targets for a given transformation. This flexibility moves the config outside the system and makes having an external config system possible.
The AST syntax supports Comparators, Math Operations, using the sent value as well as static values.

```json
{
    "type": "COMPARATOR_NODE",
    "comparator": "GREATER_EQUAL",
    "left": {
        "type": "MATH_NODE",
        "operator": "SUBTRACT",
        "left": {
            "type": "INPUT_NODE"
        },
        "right": {
            "type": "VALUE_NODE",
            "value": 18.0
        }
    },
    "right": {
        "type": "VALUE_NODE",
        "value": 5.0
    },
    "resultTrue": {
        "type": "VALUE_NODE",
        "value": 1.0
    },
    "resultFalse": {
        "type": "VALUE_NODE",
        "value": 0.0
    }
}
```

As [Go](https://go.dev/) does not allow for [inheritance](https://www.geeksforgeeks.org/inheritance-in-golang/), using [Polymorphism](https://en.wikipedia.org/wiki/Polymorphism_(computer_science)) is not as straightforward as in [Object-Oriented](https://www.techtarget.com/searchapparchitecture/definition/object-oriented-programming-OOP) languages such as [Java](https://www.java.com/en/). We currently need to pass a type field in each node of the function expression denoting the type to be able to emulate polymorphism in deserialiying differen Node types that each have different evaluation logic.

## Links:
* RabbitMQ: [https://www.rabbitmq.com](https://www.rabbitmq.com)