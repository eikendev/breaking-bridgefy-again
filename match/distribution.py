import logging
import nltk
import warnings

logger = logging.getLogger('match')


def _make_dist(args, lengths, counts):
    counts = dict(zip(lengths, counts))
    fdist = nltk.FreqDist(counts)

    if args.method == 'good-turing':
        warnings.filterwarnings("ignore", module="nltk")
        dist = nltk.SimpleGoodTuringProbDist(fdist)
    elif args.method == 'laplace':
        dist = nltk.LaplaceProbDist(fdist)
    else:
        logger.error(f'Method {args.method} is not supported')
        exit(1)

    return dist


# https://www.nltk.org/api/nltk.html
def _derive_distributions(args, sample, max_hop):
    """ Precalculate the distributions for each (hop,length) """

    sample = sample.copy()
    del sample['PreviousLength']

    payloads = sample.Payload.unique().tolist()

    distributions = list()

    for hop in range(0, max_hop + 1):
        payload_map = dict()
        distributions.append(payload_map)

        hop_sample = sample[sample.Hop == hop]

        for payload in payloads:
            counts = hop_sample[hop_sample.Payload == payload]

            assert len(counts) > 0

            payload_map[payload] = _make_dist(args, counts.Length, counts.Count)

    return distributions


def _derive_memory(args, sample, max_hop):
    payloads = sample.Payload.unique().tolist()

    memory = list()

    for hop in range(0, max_hop + 1):
        payload_map = dict()
        memory.append(payload_map)

        hop_sample = sample[sample.Hop == hop]

        for payload in payloads:
            previous_length_map = dict()
            payload_map[payload] = previous_length_map

            payload_sample = hop_sample[hop_sample.Payload == payload]
            previous_lengths = payload_sample.PreviousLength.unique().tolist()

            for previous_length in previous_lengths:
                counts = payload_sample[payload_sample.PreviousLength == previous_length]

                assert len(counts) > 0

                previous_length_map[previous_length] = _make_dist(args, counts.Length, counts.Count)

    return memory


def derive_distributions(args, sample, max_hop=None):
    if args.method not in ['laplace', 'good-turing']:
        raise ValueError('method for smoothing not supported')

    if max_hop is None:
        max_hop = sample.Hop.max()

    logger.info(f'Deriving distributions with max_hop={max_hop}')

    distributions = _derive_distributions(args, sample, max_hop)

    if args.memory:
        memory = _derive_memory(args, sample, max_hop)
        return (distributions, memory)
    else:
        return (distributions, None)


def get_distribution(dists, hop, payload):
    if payload not in dists[hop]:
        msg = f'No frequency distribution available for hop={hop} payload={payload}'
        logger.warning(msg)
        return None

    return dists[hop][payload]


def get_memory_distribution(dists, hop, payload, previous_length):
    if payload not in dists[hop] or previous_length not in dists[hop][payload]:
        msg = f'No frequency distribution available for hop={hop} payload={payload} previous_length={previous_length}'
        logger.warning(msg)
        return None

    return dists[hop][payload][previous_length]
