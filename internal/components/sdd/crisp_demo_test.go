package sdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRISPForecastingDemoFixtureExists(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "crisp-forecasting-demo")
	paths := []string{
		"pyproject.toml",
		filepath.Join("data", "daily_sales.csv"),
		filepath.Join("src", "forecasting_demo", "__init__.py"),
		filepath.Join("src", "forecasting_demo", "data.py"),
		filepath.Join("src", "forecasting_demo", "model.py"),
		filepath.Join("tests", "test_forecasting.py"),
	}
	for _, rel := range paths {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("missing demo fixture file %s: %v", rel, err)
		}
	}
}

func TestCRISPForecastingDemoDocumentsMetrics(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "crisp-forecasting-demo")
	data, err := os.ReadFile(filepath.Join(root, "tests", "test_forecasting.py"))
	if err != nil {
		t.Fatalf("read test_forecasting.py: %v", err)
	}
	text := string(data)
	for _, want := range []string{"baseline", "MAE", "candidate", "assert"} {
		if !strings.Contains(text, want) {
			t.Fatalf("demo tests missing %q", want)
		}
	}
}
