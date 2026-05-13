from pathlib import Path

import pandas as pd


REQUIRED_COLUMNS = {"date", "sales", "promo_flag", "weekday"}


def load_sales(path: str | Path) -> pd.DataFrame:
    frame = pd.read_csv(path, parse_dates=["date"])
    missing = REQUIRED_COLUMNS.difference(frame.columns)
    if missing:
        raise ValueError(f"missing required columns: {sorted(missing)}")
    if frame["sales"].isna().any():
        raise ValueError("sales contains missing values")
    return frame.sort_values("date").reset_index(drop=True)
