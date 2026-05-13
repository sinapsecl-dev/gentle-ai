from dataclasses import dataclass

import pandas as pd
from sklearn.dummy import DummyRegressor
from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import mean_absolute_error


FEATURES = ["promo_flag", "weekday"]


@dataclass(frozen=True)
class ForecastResult:
    baseline_mae: float
    candidate_mae: float


def time_split(frame: pd.DataFrame, holdout: int = 3) -> tuple[pd.DataFrame, pd.DataFrame]:
    if len(frame) <= holdout:
        raise ValueError("not enough rows for holdout split")
    return frame.iloc[:-holdout].copy(), frame.iloc[-holdout:].copy()


def evaluate_forecast(frame: pd.DataFrame) -> ForecastResult:
    train, test = time_split(frame)
    x_train = train[FEATURES]
    y_train = train["sales"]
    x_test = test[FEATURES]
    y_test = test["sales"]

    baseline = DummyRegressor(strategy="mean")
    baseline.fit(x_train, y_train)
    baseline_pred = baseline.predict(x_test)

    candidate = RandomForestRegressor(n_estimators=20, random_state=42)
    candidate.fit(x_train, y_train)
    candidate_pred = candidate.predict(x_test)

    return ForecastResult(
        baseline_mae=float(mean_absolute_error(y_test, baseline_pred)),
        candidate_mae=float(mean_absolute_error(y_test, candidate_pred)),
    )
