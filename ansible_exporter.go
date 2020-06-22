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

        var result map[string]interface{}
        json.Unmarshal([]byte(body), &result)

        global_custom_stats := result["global_custom_stats"].(map[string]interface{})
        var playbook = ""
        for key, value := range global_custom_stats {
          // Each value is an interface{} type, that is type asserted as a string
          if(key != "") {
             playbook = value.(string)
          }
        }

        plays := result["plays"].([]interface {})

        var date_start_unix_out, date_end_unix_out int64 = 0, 0
        //duration_sec := 0

        for key := range plays {
            //log.Println(plays[key]["play"])
            keysPerHost := plays[key].(map[string]interface{})

            for keyH, valueH := range keysPerHost {
                if (keyH == "play") {
                    date_layout := "2006-01-02T15:04:05.000000Z"
                    date_start := valueH.(map[string]interface {})["duration"].(map[string]interface {})["start"].(string)
                    date_end := valueH.(map[string]interface {})["duration"].(map[string]interface {})["end"].(string)
                    date_start_unix, _ := time.Parse(date_layout, date_start)
                    date_end_unix, _ := time.Parse(date_layout, date_end)
                    date_start_unix_out = date_start_unix.Unix()
                    date_end_unix_out = date_end_unix.Unix()
                }
            }
        }


        stats := result["stats"].(map[string]interface{})

        //log.Println(stats["k8s-0-edgeduck-prd"])

        for key := range stats {
          // Each value is an interface{} type, that is type asserted as a string
          statsPerHost := stats[key].(map[string]interface{})
          // stats for duration and timestamp
          ansiblePlayStartGauge.WithLabelValues(playbook, key).Set(float64(date_start_unix_out))
          ansiblePlayEndGauge.WithLabelValues(playbook, key).Set(float64(date_end_unix_out))
          ansiblePlayDurationGauge.WithLabelValues(playbook, key).Set(float64(date_end_unix_out - date_start_unix_out))

          // stats with status
          for keyH, valueH := range statsPerHost {
            ansibleRunStatGauge.WithLabelValues(playbook, key, keyH).Set(valueH.(float64))
          }
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
