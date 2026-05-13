from pathlib import Path

from forecasting_demo.data import load_sales
from forecasting_demo.model import evaluate_forecast, time_split


DATA_PATH = Path(__file__).resolve().parents[1] / "data" / "daily_sales.csv"


def test_load_sales_validates_required_columns():
    frame = load_sales(DATA_PATH)
    assert list(frame.columns) == ["date", "sales", "promo_flag", "weekday"]
    assert frame["sales"].isna().sum() == 0


def test_time_split_keeps_future_rows_for_holdout():
    frame = load_sales(DATA_PATH)
    train, test = time_split(frame, holdout=3)
    assert train["date"].max() < test["date"].min()
    assert len(test) == 3


def test_candidate_model_reports_mae_against_baseline():
    # Compare candidate model against baseline using MAE.
    frame = load_sales(DATA_PATH)
    result = evaluate_forecast(frame)
    assert result.baseline_mae >= 0
    assert result.candidate_mae >= 0
    assert result.candidate_mae < result.baseline_mae * 2
