import logging

logger = logging.getLogger('match')


def read_file_lines(path):
    with open(path) as f:
        try:
            lines = f.readlines()
        except:
            logger.error(f'Cannot load {path}')
            exit(1)

    lines = [line.rstrip() for line in lines]

    return lines
