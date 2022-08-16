#!/usr/bin/python
# -*- coding: utf-8 -*-

import os, pika, json, pigpio, time

pwm = pigpio.pi()
FAN_PIN = 2

class Fan:
    def __init__(self, pin: int):
        self.speed = 0
        self.pin = pin
        pwm.set_PWM_frequency(self.pin, 100)
        #print(pwm.get_PWM_frequency(self.pin))
        self.set_speed(0)
    def set_speed(self, speed: float = 1.0):
        assert speed >= 0 and speed <= 1.0, 'Invalid fan speed'
        pwm.set_PWM_dutycycle(self.pin, int(speed * 255.0))
    def stop(self):
        self.set_speed(0)

fan = Fan(FAN_PIN)

url = os.environ.get('RABBIT_URL', 'amqp://guest:guest@localhost:5672/')
actuator_id = 'DummyTarget'

print(f'Connecting to {url}')
connection = pika.BlockingConnection(pika.URLParameters(url))
channel = connection.channel()
result = channel.queue_declare(queue='', exclusive=True)
queue_name = result.method.queue
channel.queue_bind(exchange='router_output', queue=queue_name, routing_key=actuator_id)

def callback(channel, method, properties, body):
    data = json.loads(body)
    value = data.get('value', 0.0)
    print(f'Received command: set fan speed to: { value }')
    fan.set_speed(value)

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)

print(f'Listening on exchange: "router_output" for queue { queue_name } on key: { actuator_id }')

if __name__ == "__main__":
    try:
        fan.set_speed(0)
        channel.start_consuming()
    except:
        fan.stop()
        pwm.stop()
        channel.stop_consuming()