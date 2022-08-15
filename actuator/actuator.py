#!/usr/bin/python
# -*- coding: utf-8 -*-

import os, pika, json

url = os.environ.get('RABBIT_URL', 'amqp://guest:guest@localhost:5672/')
actuator_id = 'DummyTarget'

print(f'Connecting to {url}')
connection = pika.BlockingConnection(pika.URLParameters(url))
channel = connection.channel()
result = channel.queue_declare(queue='', exclusive=True)
queue_name = result.method.queue
channel.queue_bind(exchange='router_output', queue=queue_name, routing_key=actuator_id)

def callback(channel, method, properties, body):
    print(f'Received message: { body }')

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)

print(f'Listening on exchange: "router_output" for queue { queue_name } on key: { actuator_id }')

channel.start_consuming()