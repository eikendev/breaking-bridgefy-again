import logging
import sys


def setup_logger(logger):
    fmt = '%(asctime)s [%(levelname)s] %(message)s'
    datefmt = '%Y-%m-%d %H:%M:%S'
    level = logging.INFO

    formatter = logging.Formatter(fmt=fmt, datefmt=datefmt)
    stream_handler = logging.StreamHandler(sys.stderr)
    stream_handler.setFormatter(formatter)

    logger.setLevel(level)
    logger.addHandler(stream_handler)
