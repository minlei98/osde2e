package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// KrknAIYAML represents the structure of krkn-ai.yaml
type KrknAIYAML struct {
	KubeconfigFilePath string                 `yaml:"kubeconfig_file_path"`
	Parameters         map[string]interface{} `yaml:"parameters,omitempty"`
	Generations        int                    `yaml:"generations"`
	PopulationSize     int                    `yaml:"population_size"`
	WaitDuration       int                    `yaml:"wait_duration"`
	MutationRate       float64                `yaml:"mutation_rate,omitempty"`
	ScenarioMutationRate float64              `yaml:"scenario_mutation_rate,omitempty"`
	CrossoverRate      float64                `yaml:"crossover_rate,omitempty"`
	CompositionRate    float64                `yaml:"composition_rate,omitempty"`
	PopulationInjectionRate float64           `yaml:"population_injection_rate,omitempty"`
	PopulationInjectionSize int                `yaml:"population_injection_size,omitempty"`
	FitnessFunction    FitnessFunction        `yaml:"fitness_function"`
	HealthChecks       HealthChecks           `yaml:"health_checks"`
	Scenario           Scenario               `yaml:"scenario"`
	ClusterComponents  ClusterComponents      `yaml:"cluster_components"`
}

// FitnessFunction represents the fitness function configuration
type FitnessFunction struct {
	Query                          string        `yaml:"query"`
	Type                           string        `yaml:"type"`
	IncludeKrknFailure             bool          `yaml:"include_krkn_failure"`
	IncludeHealthCheckFailure      bool          `yaml:"include_health_check_failure"`
	IncludeHealthCheckResponseTime bool          `yaml:"include_health_check_response_time"`
	Items                          []interface{} `yaml:"items"`
}

// HealthChecks represents the health checks configuration
type HealthChecks struct {
	StopWatcherOnFailure bool             `yaml:"stop_watcher_on_failure"`
	Applications         []HealthCheckApp `yaml:"applications"`
}

// HealthCheckApp represents a single health check application
type HealthCheckApp struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	StatusCode int    `yaml:"status_code"`
	Timeout    int    `yaml:"timeout"`
	Interval   int    `yaml:"interval"`
}

// Scenario represents chaos scenario toggles
type Scenario struct {
	ApplicationOutages ScenarioToggle `yaml:"application_outages"`
	PodScenarios       ScenarioToggle `yaml:"pod_scenarios"`
	ContainerScenarios ScenarioToggle `yaml:"container_scenarios"`
	NodeCPUHog         ScenarioToggle `yaml:"node_cpu_hog"`
	NodeMemoryHog      ScenarioToggle `yaml:"node_memory_hog"`
	NodeIOHog          ScenarioToggle `yaml:"node_io_hog,omitempty"`
	TimeScenarios      ScenarioToggle `yaml:"time_scenarios"`
	NetworkScenarios   ScenarioToggle `yaml:"network_scenarios"`
	DNSOutage          ScenarioToggle `yaml:"dns_outage"`
	SynFlood           ScenarioToggle `yaml:"syn_flood,omitempty"`
}

// ScenarioToggle represents a scenario enable/disable toggle
type ScenarioToggle struct {
	Enable bool `yaml:"enable"`
}

// ClusterComponents represents discovered cluster components
type ClusterComponents struct {
	Namespaces []interface{} `yaml:"namespaces"`
	Nodes      []interface{} `yaml:"nodes,omitempty"`
}

// UpdateKrknAIYAMLWithJenkinsParams updates the discovered krkn-ai.yaml with Jenkins parameters
// This function merges user-provided Jenkins parameters with the auto-discovered cluster configuration
func (e *Executor) UpdateKrknAIYAMLWithJenkinsParams(discoveredYAMLPath string) error {
	if e.cfg.KrknAIConfig == nil {
		return fmt.Errorf("KrknAIConfig is nil")
	}

	e.logger.Info("Updating krkn-ai.yaml with Jenkins parameters", "file", discoveredYAMLPath)

	// Read discovered YAML
	yamlData, err := os.ReadFile(discoveredYAMLPath)
	if err != nil {
		return fmt.Errorf("reading discovered yaml: %w", err)
	}

	var krknConfig KrknAIYAML
	if err := yaml.Unmarshal(yamlData, &krknConfig); err != nil {
		return fmt.Errorf("unmarshaling yaml: %w", err)
	}

	// Update Genetic Algorithm Parameters
	if e.cfg.KrknAIConfig.Generations != "" {
		if gen, err := strconv.Atoi(e.cfg.KrknAIConfig.Generations); err == nil {
			e.logger.Info("Updating generations", "from", krknConfig.Generations, "to", gen)
			krknConfig.Generations = gen
		} else {
			e.logger.Error(err, "invalid generations value", "value", e.cfg.KrknAIConfig.Generations)
		}
	}

	if e.cfg.KrknAIConfig.PopulationSize != "" {
		if pop, err := strconv.Atoi(e.cfg.KrknAIConfig.PopulationSize); err == nil {
			e.logger.Info("Updating population_size", "from", krknConfig.PopulationSize, "to", pop)
			krknConfig.PopulationSize = pop
		} else {
			e.logger.Error(err, "invalid population_size value", "value", e.cfg.KrknAIConfig.PopulationSize)
		}
	}

	if e.cfg.KrknAIConfig.WaitDuration != "" {
		if wait, err := strconv.Atoi(e.cfg.KrknAIConfig.WaitDuration); err == nil {
			e.logger.Info("Updating wait_duration", "from", krknConfig.WaitDuration, "to", wait)
			krknConfig.WaitDuration = wait
		} else {
			e.logger.Error(err, "invalid wait_duration value", "value", e.cfg.KrknAIConfig.WaitDuration)
		}
	}

	if e.cfg.KrknAIConfig.CompositionRate != "" {
		if rate, err := strconv.ParseFloat(e.cfg.KrknAIConfig.CompositionRate, 64); err == nil {
			e.logger.Info("Updating composition_rate", "from", krknConfig.CompositionRate, "to", rate)
			krknConfig.CompositionRate = rate
		} else {
			e.logger.Error(err, "invalid composition_rate value", "value", e.cfg.KrknAIConfig.CompositionRate)
		}
	}

	// Update Scenario Toggles
	e.updateScenarioToggle("pod_scenarios", e.cfg.KrknAIConfig.EnablePodScenarios, &krknConfig.Scenario.PodScenarios)
	e.updateScenarioToggle("container_scenarios", e.cfg.KrknAIConfig.EnableContainerScenarios, &krknConfig.Scenario.ContainerScenarios)
	e.updateScenarioToggle("node_cpu_hog", e.cfg.KrknAIConfig.EnableNodeCPUHog, &krknConfig.Scenario.NodeCPUHog)
	e.updateScenarioToggle("node_memory_hog", e.cfg.KrknAIConfig.EnableNodeMemoryHog, &krknConfig.Scenario.NodeMemoryHog)
	e.updateScenarioToggle("node_io_hog", e.cfg.KrknAIConfig.EnableNodeIOHog, &krknConfig.Scenario.NodeIOHog)
	e.updateScenarioToggle("network_scenarios", e.cfg.KrknAIConfig.EnableNetworkScenarios, &krknConfig.Scenario.NetworkScenarios)
	e.updateScenarioToggle("dns_outage", e.cfg.KrknAIConfig.EnableDNSOutage, &krknConfig.Scenario.DNSOutage)
	e.updateScenarioToggle("time_scenarios", e.cfg.KrknAIConfig.EnableTimeScenarios, &krknConfig.Scenario.TimeScenarios)

	// Update Fitness Function Query
	if e.cfg.KrknAIConfig.FitnessFunctionQuery != "" {
		e.logger.Info("Updating fitness_function.query", "to", e.cfg.KrknAIConfig.FitnessFunctionQuery)
		krknConfig.FitnessFunction.Query = e.cfg.KrknAIConfig.FitnessFunctionQuery
	}

	// Update Health Checks URL
	if e.cfg.KrknAIConfig.HealthChecksURL != "" {
		e.logger.Info("Updating health_checks URL", "to", e.cfg.KrknAIConfig.HealthChecksURL)
		if len(krknConfig.HealthChecks.Applications) > 0 {
			// Update first health check application
			oldURL := krknConfig.HealthChecks.Applications[0].URL
			krknConfig.HealthChecks.Applications[0].URL = e.cfg.KrknAIConfig.HealthChecksURL
			e.logger.Info("Updated health check URL", "from", oldURL, "to", e.cfg.KrknAIConfig.HealthChecksURL)
		} else {
			// Create default health check application if none exists
			krknConfig.HealthChecks.Applications = []HealthCheckApp{
				{
					Name:       "cluster-health",
					URL:        e.cfg.KrknAIConfig.HealthChecksURL,
					StatusCode: 200,
					Timeout:    4,
					Interval:   2,
				},
			}
			e.logger.Info("Created new health check application", "url", e.cfg.KrknAIConfig.HealthChecksURL)
		}
	}

	// Update Host parameter
	if e.cfg.KrknAIConfig.Host != "" {
		e.logger.Info("Updating HOST parameter", "to", e.cfg.KrknAIConfig.Host)
		if krknConfig.Parameters == nil {
			krknConfig.Parameters = make(map[string]interface{})
		}
		krknConfig.Parameters["HOST"] = e.cfg.KrknAIConfig.Host
	}

	// Write updated YAML back
	updatedYAML, err := yaml.Marshal(&krknConfig)
	if err != nil {
		return fmt.Errorf("marshaling updated yaml: %w", err)
	}

	// Overwrite the original file with updated configuration
	if err := os.WriteFile(discoveredYAMLPath, updatedYAML, 0644); err != nil {
		return fmt.Errorf("writing updated yaml: %w", err)
	}

	// Also save a backup copy
	backupPath := filepath.Join(filepath.Dir(discoveredYAMLPath), "krkn-ai-updated.yaml")
	if err := os.WriteFile(backupPath, updatedYAML, 0644); err != nil {
		e.logger.Error(err, "failed to write backup yaml", "path", backupPath)
		// Don't return error, backup is optional
	} else {
		e.logger.Info("Created backup of updated config", "path", backupPath)
	}

	e.logger.Info("Successfully updated krkn-ai.yaml with Jenkins parameters", "file", discoveredYAMLPath)

	return nil
}

// updateScenarioToggle is a helper function to update scenario enable/disable flags
func (e *Executor) updateScenarioToggle(name string, value string, toggle *ScenarioToggle) {
	if value != "" {
		if enable, err := strconv.ParseBool(value); err == nil {
			oldValue := toggle.Enable
			toggle.Enable = enable
			e.logger.Info("Updated scenario toggle", "scenario", name, "from", oldValue, "to", enable)
		} else {
			e.logger.Error(err, "invalid boolean value for scenario", "scenario", name, "value", value)
		}
	}
}

// ValidateKrknAIConfig validates the KrknAI configuration parameters
func ValidateKrknAIConfig(cfg *KrknAIConfig) error {
	if cfg == nil {
		return fmt.Errorf("KrknAIConfig is nil")
	}

	// Validate mode
	if cfg.Mode != "discover" && cfg.Mode != "run" {
		return fmt.Errorf("invalid mode: %s (must be 'discover' or 'run')", cfg.Mode)
	}

	// Validate numeric parameters if provided
	if cfg.Generations != "" {
		if _, err := strconv.Atoi(cfg.Generations); err != nil {
			return fmt.Errorf("invalid generations value: %s", cfg.Generations)
		}
	}

	if cfg.PopulationSize != "" {
		if _, err := strconv.Atoi(cfg.PopulationSize); err != nil {
			return fmt.Errorf("invalid population_size value: %s", cfg.PopulationSize)
		}
	}

	if cfg.WaitDuration != "" {
		if _, err := strconv.Atoi(cfg.WaitDuration); err != nil {
			return fmt.Errorf("invalid wait_duration value: %s", cfg.WaitDuration)
		}
	}

	// Validate boolean parameters if provided
	boolParams := map[string]string{
		"enable_pod_scenarios":       cfg.EnablePodScenarios,
		"enable_container_scenarios": cfg.EnableContainerScenarios,
		"enable_node_cpu_hog":        cfg.EnableNodeCPUHog,
		"enable_node_memory_hog":     cfg.EnableNodeMemoryHog,
		"enable_node_io_hog":         cfg.EnableNodeIOHog,
		"enable_network_scenarios":   cfg.EnableNetworkScenarios,
		"enable_dns_outage":          cfg.EnableDNSOutage,
		"enable_time_scenarios":      cfg.EnableTimeScenarios,
	}

	for name, value := range boolParams {
		if value != "" {
			if _, err := strconv.ParseBool(value); err != nil {
				return fmt.Errorf("invalid boolean value for %s: %s", name, value)
			}
		}
	}

	return nil
}
