import logging

from argparse import ArgumentParser

logger = logging.getLogger('match')


def parse_arguments():
    parser = ArgumentParser(
        prog='match',
        description='Matches an attack sample against a simulation sample.'
    )

    parser.add_argument('-q', '--quiet', action='store_true',
                        help='Print errors and warnings only.')
    parser.add_argument('--memory', action='store_true',
                        help='Assume the sample files have memory.')
    parser.add_argument('--runners', type=int, default=8,
                        help='How many runner processes to spawn for processing.')
    parser.add_argument('--max-hop', type=int, default=None,
                        help='Consider packets up to this hop.')
    parser.add_argument('--strict', action='store_true',
                        help='Give strict penalty when a match is unlikely.')
    parser.add_argument('--only-gzip', action='store_true',
                        help='Disregard the encryption layer.')
    parser.add_argument('--method', type=str, default='laplace', choices=['laplace', 'good-turing'],
                        help='The frequency smoothing method to use.')
    parser.add_argument('--priors', type=str, metavar='priors-file',
                        help='The priors file. This will use MAP instead of MLE.')
    parser.add_argument('simulation_file', type=str, metavar='simulation-file',
                        help='The simulation sample.')
    parser.add_argument('attack_files', type=str, metavar='attack-file', nargs='*',
                        help='The attack sample(s).')

    args = parser.parse_args()

    return args
