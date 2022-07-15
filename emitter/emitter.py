#!/usr/bin/python
# -*- coding: utf-8 -*-

import re, os, time, pika, json

url = os.environ.get('RABBIT_URL', 'amqp://guest:guest@localhost:5672/')
emitter_id = 'TemperatureSensor'

def read_temperature_sensor(path):
    value = 'U'
    try:
        f = open(path, 'r')
        line = f.readline()
        if re.match(r'([0-9a-f]{2} ){9}: crc=[0-9a-f]{2} YES', line):
            line = f.readline()
            m = re.match(r'([0-9a-f]{2} ){9}t=([+-]?[0-9]+)', line)
            if m:
                value = str(float(m.group(2)) / 1000.0)
        f.close()
    except IOError as e:
        print(time.strftime('%x %X'), 'Error reading', path, ' : ', e)
    return float(value) - 1

path = '/sys/devices/w1_bus_master1/28-000006dd9b7c/w1_slave'

temperature = read_temperature_sensor(path)

payload = {'emitterId': emitter_id, 'value': str(temperature)}

print(f'Connecting to {url}')
connection = pika.BlockingConnection(pika.URLParameters(url))
channel = connection.channel()
channel.basic_publish(exchange='router_input', routing_key='', body=json.dumps(payload))
print(f'[x] Sent "{payload}"')