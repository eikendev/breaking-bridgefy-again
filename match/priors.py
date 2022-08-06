import pandas as pd


def load_priors(args):
    df = pd.read_csv(args.priors, sep=' ', header=None)
    df.rename(columns={0: 'count', 1: 'payload'}, inplace=True)
    df['count'] = pd.to_numeric(df['count'])
    df.set_index('payload', inplace=True)
    df['count'] /= df['count'].sum()

    return df['count'].to_dict()
