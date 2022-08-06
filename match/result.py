from dataclasses import dataclass


@dataclass(frozen=True)
class Result:
    sample_size: int
    max_hop: int
    payload: str
    rank: int
