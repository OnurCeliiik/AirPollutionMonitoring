package anomaly

import (
	"math"
	"sort"

	"github.com/user/airpollution/internal/models"
)

// WHO limits for common air pollutants (in μg/m³)
const (
	PM25Limit = 15.0  // Annual mean
	PM10Limit = 45.0  // Annual mean
	NO2Limit  = 25.0  // Annual mean
	O3Limit   = 100.0 // 8-hour mean
)

// Detector is responsible for detecting anomalies in air quality data
type Detector struct {
	historicalData map[string][]models.AirQualityData // Map of parameter to historical data
}

// NewDetector creates a new anomaly detector
func NewDetector() *Detector {
	return &Detector{
		historicalData: make(map[string][]models.AirQualityData),
	}
}

// AddHistoricalData adds historical data for a parameter
func (d *Detector) AddHistoricalData(data []models.AirQualityData) {
	for _, item := range data {
		if _, ok := d.historicalData[item.Parameter]; !ok {
			d.historicalData[item.Parameter] = []models.AirQualityData{}
		}
		d.historicalData[item.Parameter] = append(d.historicalData[item.Parameter], item)
	}
}

// Detect checks for anomalies in the given data point
func (d *Detector) Detect(data *models.AirQualityData, recentData []models.AirQualityData) (*models.Anomaly, error) {
	// Check for threshold exceedance
	if anomaly := d.checkThresholdExceeded(data); anomaly != nil {
		return anomaly, nil
	}

	// Check for statistical outlier (Z-score)
	if anomaly := d.checkStatisticalOutlier(data, recentData); anomaly != nil {
		return anomaly, nil
	}

	// Check for spike detection
	if anomaly := d.checkSpikeDetection(data, recentData); anomaly != nil {
		return anomaly, nil
	}

	// Check for geographic inconsistency
	if anomaly := d.checkGeographicInconsistency(data, recentData); anomaly != nil {
		return anomaly, nil
	}

	return nil, nil
}

// checkThresholdExceeded checks if the value exceeds WHO limits
func (d *Detector) checkThresholdExceeded(data *models.AirQualityData) *models.Anomaly {
	var limit float64

	switch data.Parameter {
	case "PM2.5":
		limit = PM25Limit
	case "PM10":
		limit = PM10Limit
	case "NO2":
		limit = NO2Limit
	case "O3":
		limit = O3Limit
	default:
		return nil // No known threshold for this parameter
	}

	if data.Value > limit {
		return models.NewAnomalyFromData(
			string(models.ThresholdExceeded),
			data,
		)
	}

	return nil
}

// checkStatisticalOutlier uses Z-score to detect outliers
func (d *Detector) checkStatisticalOutlier(data *models.AirQualityData, recentData []models.AirQualityData) *models.Anomaly {
	if len(recentData) < 10 {
		return nil // Not enough data for statistical analysis
	}

	// Extract values
	values := make([]float64, len(recentData))
	for i, d := range recentData {
		values[i] = d.Value
	}

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(len(values)))

	// Calculate Z-score
	zScore := (data.Value - mean) / stdDev

	// If Z-score is greater than 3, it's an outlier (99.7% confidence)
	if math.Abs(zScore) > 3.0 {
		return models.NewAnomalyFromData(
			string(models.StatisticalOutlier),
			data,
		)
	}

	return nil
}

// checkSpikeDetection checks for sudden spikes in values
func (d *Detector) checkSpikeDetection(data *models.AirQualityData, recentData []models.AirQualityData) *models.Anomaly {
	if len(recentData) < 5 {
		return nil // Not enough data for spike detection
	}

	// Sort by timestamp to ensure chronological order
	sort.Slice(recentData, func(i, j int) bool {
		return recentData[i].Timestamp.Before(recentData[j].Timestamp)
	})

	// Calculate average of recent values (last 24 hours)
	var sum float64
	for _, d := range recentData {
		sum += d.Value
	}
	avg := sum / float64(len(recentData))

	// Check if current value is 50% higher than the average
	if data.Value > avg*1.5 {
		return models.NewAnomalyFromData(
			string(models.SpikeDetected),
			data,
		)
	}

	return nil
}

// checkGeographicInconsistency checks if value is inconsistent with nearby readings
func (d *Detector) checkGeographicInconsistency(data *models.AirQualityData, recentData []models.AirQualityData) *models.Anomaly {
	if len(recentData) < 3 {
		return nil // Not enough data for geographic consistency check
	}

	// Find nearby readings within the last hour (approximately 25km radius)
	const (
		latDelta = 0.25 // Roughly 25km
		lonDelta = 0.25
	)

	// Filter nearby readings
	var nearbyReadings []models.AirQualityData
	for _, d := range recentData {
		latDiff := math.Abs(data.Latitude - d.Latitude)
		lonDiff := math.Abs(data.Longitude - d.Longitude)

		if latDiff <= latDelta && lonDiff <= lonDelta {
			nearbyReadings = append(nearbyReadings, d)
		}
	}

	if len(nearbyReadings) < 3 {
		return nil // Not enough nearby readings
	}

	// Calculate median of nearby readings
	values := make([]float64, len(nearbyReadings))
	for i, d := range nearbyReadings {
		values[i] = d.Value
	}
	sort.Float64s(values)

	var median float64
	if len(values)%2 == 0 {
		median = (values[len(values)/2-1] + values[len(values)/2]) / 2
	} else {
		median = values[len(values)/2]
	}

	// Check if current value is significantly different (>3x) from median
	if data.Value > median*3 || (median > 0 && data.Value*3 < median) {
		return models.NewAnomalyFromData(
			string(models.GeographicInconsistency),
			data,
		)
	}

	return nil
}
