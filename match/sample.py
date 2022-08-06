import json
import logging
import pandas as pd

AES_BLOCKSIZE = 16

logger = logging.getLogger('match')


def load_sample(args, line):
    try:
        data = json.loads(line)
    except:
        logger.error('Cannot parse JSON data')
        exit(1)

    if type(list(list(list(data.values())[0].values())[0].values())[0]) == dict:
        data = (
            (payload, hop, length, previous_length, count)
            for payload, hops in data.items()
            for hop, lengths in hops.items()
            for length, previous_lengths in lengths.items()
            for previous_length, count in previous_lengths.items()
        )
    else:
        data = (
            (payload, hop, length, 0, count)
            for payload, hops in data.items()
            for hop, lengths in hops.items()
            for length, count in lengths.items()
        )

    sample = pd.DataFrame(data, columns=['Payload', 'Hop', 'Length', 'PreviousLength', 'Count'])
    sample['Hop'] = pd.to_numeric(sample['Hop'])
    sample['Length'] = pd.to_numeric(sample['Length'])
    sample['PreviousLength'] = pd.to_numeric(sample['PreviousLength'])

    if 'PreviousLength' not in sample.columns:
        sample['PreviousLength'] = 0

    if not args.memory:
        sample['PreviousLength'] = 0
        sample = sample.groupby(['Payload', 'Hop', 'Length', 'PreviousLength']).sum().reset_index()

    if not args.only_gzip:
        sample.Length = sample.Length + (AES_BLOCKSIZE - sample.Length % AES_BLOCKSIZE)
        sample.PreviousLength = sample.PreviousLength + (AES_BLOCKSIZE - sample.PreviousLength % AES_BLOCKSIZE)
        sample = sample.groupby(['Payload', 'Hop', 'Length', 'PreviousLength']).sum().reset_index()

    # At h=0, there is never a valid previous length.
    sample.loc[sample['Hop'] == 0, 'PreviousLength'] = 0

    return sample
