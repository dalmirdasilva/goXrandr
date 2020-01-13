package display

const (
  axisX Axis = 0
  axisY Axis = 1
)

type Position struct {
  X int `json:"x"`
  Y int `json:"y"`
}

type Mode struct {
  Width  int `json:"width"`
  Height int `json:"height"`
}

type Axis int

type Display struct {
  Output       string   `json:"output"`
  Connected    bool     `json:"connected"`
  Inverted     bool     `json:"inverted"`
  Scale        int      `json:"scale"`
  Frequency    int      `json:"frequency"`
  InvertedAxis Axis     `json:"invertedAxis"`
  Pos          Position `json:"pos"`
  Mode         Mode     `json:"mode"`
}

type Arrangement []Display

type Preference struct {
  When        []string    `json:"when"`
  Arrangement Arrangement `json:"arrangement"`
}
