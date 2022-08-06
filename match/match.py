import sys
import logging
import math

from collections import OrderedDict
from itertools import chain
from pathlib import Path

from joblib import Parallel, delayed

from .arguments import parse_arguments
from .distribution import get_distribution, get_memory_distribution, derive_distributions
from .logging import setup_logger
from .priors import load_priors
from .result import Result
from .files import read_file_lines
from .sample import load_sample

logger = logging.getLogger('match')


def _get_log_likelihood(args, dist, length, payload, priors):
    p = 0 if dist is None else dist.prob(length)

    if not args.strict and p == 0:
        p = 0.000001

    if priors is not None:
        p *= priors[payload]

    return math.log(p) if p > 0 else float('-inf')


def match(args, dists, sample, payloads, priors, max_hop=None):
    if max_hop is not None:
        sample = sample[sample.Hop <= max_hop]

    logger.info(f'Matching distrubutions with max_hop={max_hop}')

    likelihoods = {payload: 0 for payload in payloads}

    ds_single, ds_memory = dists

    for _, entry in sample.iterrows():
        for payload in payloads:
            if args.memory and entry.Hop > 0:
                dist = get_memory_distribution(ds_memory, entry.Hop, payload, entry.PreviousLength)
            else:
                dist = get_distribution(ds_single, entry.Hop, payload)

            logl = _get_log_likelihood(args, dist, entry.Length, payload, priors)
            likelihoods[payload] += entry.Count * logl

    return likelihoods


def make_result(likelihoods, sample_size, max_hop, atk_payload):
    sorted_likelihoods = OrderedDict(sorted(likelihoods.items(), key=lambda x: x[1], reverse=True))

    for payload, likelihood in sorted_likelihoods.items():
        logger.info(f"'{payload}'\t{likelihood}")

    sorted_keys = list(sorted_likelihoods.keys())
    payload_idx = sorted_keys.index(atk_payload) + 1

    return Result(sample_size, max_hop, atk_payload, payload_idx)


def process_attack(args, data, max_hop, dists, sim_payloads, priors):
    s_atk = load_sample(args, data)
    atk_payloads = s_atk.Payload.unique().tolist()
    del s_atk['Payload']

    if len(atk_payloads) != 1:
        logger.error('Attack sample must have exactly one payload; skipping...')
        return None

    real_payload = atk_payloads[0]
    logger.info(f"Attack payload is '{real_payload}'")

    atk_sample_size = s_atk.groupby(['Hop']).sum().head(1).Count.tolist()
    if len(atk_sample_size) == 0:
        logger.error('Attack sample cannot be empty; skipping...')
        return None

    atk_sample_size = atk_sample_size[0]
    if atk_sample_size == 0:
        logger.error('Attack sample cannot be empty; skipping...')
        return None

    atk_max_hop = s_atk.Hop.max()
    current_max_hop = min(max_hop, atk_max_hop)

    if atk_max_hop < max_hop:
        logger.error(f'Attack sample only supports max_hop<={atk_max_hop}; skipping...')
        return None

    payloads = set(sim_payloads + atk_payloads)
    likelihoods = match(args, dists, s_atk, payloads, priors, current_max_hop)

    result = make_result(likelihoods, atk_sample_size, current_max_hop, real_payload)

    if args.runners == 1:
        print(f"{result.sample_size}\t{result.max_hop}\t{result.payload}\t{result.rank}")

    return result


def process_attack_file(args, path, max_hop, dists, sim_payloads, priors):
    lines = read_file_lines(path)
    results = list()

    for line in lines:
        results.append(process_attack(args, line, max_hop, dists, sim_payloads, priors))

    return results


def main():
    setup_logger(logger)

    args = parse_arguments()

    if args.quiet:
        logger.setLevel(logging.WARNING)

    priors = load_priors(args) if args.priors is not None else None

    sim_filepath = Path(args.simulation_file)
    s_sim = load_sample(args, read_file_lines(sim_filepath)[0])

    sim_max_hop = s_sim.Hop.max()
    max_hop = sim_max_hop if args.max_hop is None else min(args.max_hop, sim_max_hop)

    if args.max_hop is not None and sim_max_hop < args.max_hop:
        logger.fatal(f'Simulation sample only supports max_hop<={sim_max_hop}')
        exit(1)

    dists = derive_distributions(args, s_sim, max_hop)
    sim_payloads = s_sim.Payload.unique().tolist()

    logger.info(f'Simulation file has {len(sim_payloads)} payloads')

    iter = args.attack_files if len(args.attack_files) > 0 else map(str.rstrip, sys.stdin)

    if args.runners == 1:
        logger.info('Using a single process; not scheduling jobs...')
        results = [process_attack_file(args, attack_file, max_hop, dists, sim_payloads, priors) for attack_file in iter]
    else:
        results = Parallel(n_jobs=args.runners)(
            delayed(process_attack_file)(args, attack_file, max_hop, dists, sim_payloads, priors)
            for attack_file in iter
        )

    results = chain.from_iterable(results)

    if args.runners > 1:
        for r in filter(lambda x: x is not None, results):
            print(f"{r.sample_size}\t{r.max_hop}\t{r.payload}\t{r.rank}")

    if any(filter(lambda x: x is None, results)):
        exit(1)
