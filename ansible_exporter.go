package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "time"
  log "github.com/sirupsen/logrus"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)


type AnsibleStats struct {
    Changed     int `json:"changed"`
    Failures    int `json:"failures"`
    Ignored     int `json:"ignored"`
    Ok          int `json:"ok"`
    Rescued     int `json:"rescued"`
    Skipped     int `json:"skipped"`
    Unreachable int `json:"unreachable"`
}


type AnsibleJsonRun struct {
	CustomStats struct {
	} `json:"custom_stats"`
	GlobalCustomStats struct {
	} `json:"global_custom_stats"`
	Plays []struct {
		Play struct {
			Duration struct {
				End   time.Time `json:"end"`
				Start time.Time `json:"start"`
			} `json:"duration"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"play"`
		Tasks []struct {
            Hosts map[string]interface{} `json:"hosts"`
		} `json:"tasks"`
	} `json:"plays"`
	Stats map[string]AnsibleStats `json:"stats"`
}

var (
        ansibleRunStatGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "ansible_play_stats",
             Help:      "Ansible stats status per playbook per host",
        },
        []string{
            // Which user has requested the operation?
            "playbook",
            // Of what type is the operation?
            "host",
            "status",
        })
        ansiblePlayDurationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "ansible_play_duration_sec",
             Help:      "Duration of the playbook run",
        },
        []string{
            // Which user has requested the operation?
            "playbook",
            // Of what type is the operation?
            "host",
        })
        ansiblePlayStartGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "ansible_play_date_start_tsunix",
             Help:      "Unix timestamp start",
        },
        []string{
            // Which user has requested the operation?
            "playbook",
            // Of what type is the operation?
            "host",
        })
        ansiblePlayEndGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "ansible_play_end_start_tsunix",
             Help:      "Unix timestamp end",
        },
        []string{
            // Which user has requested the operation?
            "playbook",
            // Of what type is the operation?
            "host",
        })
)

func main() {

  processPlaybookJson := func(w http.ResponseWriter, r *http.Request) {

    switch r.Method {
    case "POST":

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            panic(err)
        }

        jsonData := []byte(body)

        var ansible_run_data AnsibleJsonRun
        if json.Unmarshal(jsonData, &ansible_run_data) != nil {
            log.Println(err)
        }

        var playbook = ansible_run_data.Plays[0].Play.Name
        var date_start_unix_out = ansible_run_data.Plays[0].Play.Duration.Start.Unix()
        var date_end_unix_out = ansible_run_data.Plays[0].Play.Duration.End.Unix()

        for key := range ansible_run_data.Stats {
          ansiblePlayStartGauge.WithLabelValues(playbook, key).Set(float64(date_start_unix_out))
          ansiblePlayEndGauge.WithLabelValues(playbook, key).Set(float64(date_end_unix_out))
          ansiblePlayDurationGauge.WithLabelValues(playbook, key).Set(float64(date_end_unix_out - date_start_unix_out))

          ansibleRunStatGauge.WithLabelValues(playbook, key, "Changed").Set(float64(ansible_run_data.Stats[key].Changed))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Failures").Set(float64(ansible_run_data.Stats[key].Failures))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Ignored").Set(float64(ansible_run_data.Stats[key].Ignored))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Ok").Set(float64(ansible_run_data.Stats[key].Ok))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Rescued").Set(float64(ansible_run_data.Stats[key].Rescued))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Skipped").Set(float64(ansible_run_data.Stats[key].Skipped))
          ansibleRunStatGauge.WithLabelValues(playbook, key, "Unreachable").Set(float64(ansible_run_data.Stats[key].Unreachable))
        }

        fmt.Fprintf(w, "POST RECEIVED AND PROCESSED")
    default:
        fmt.Fprintf(w, "Sorry, only POST methods are supported.")
    }
  }

  //This section will start the HTTP server and expose
  //any metrics on the /metrics endpoint.
  http.Handle("/metrics", promhttp.Handler())
  http.HandleFunc("/ansible_injest", processPlaybookJson)
  log.Info("Beginning to serve on port :9515")
  log.Fatal(http.ListenAndServe(":9515", nil))
}
