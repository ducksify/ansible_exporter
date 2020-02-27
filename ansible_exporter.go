package main

import (
  "net/http"
  "fmt"
  "io/ioutil"
  "encoding/json"
  log "github.com/sirupsen/logrus"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
  opsProcessed.Inc()
//         go func() {
//                 for {
//                         opsProcessed.Inc()
//                 }
//         }()
}

func gaugeTestChangeValue() {
  testGauge.WithLabelValues("myhost", "myPB").Add(4)
  testGauge.WithLabelValues("myhost", "myPB3").Add(8)
//         go func() {
//                 for {
//                         opsProcessed.Inc()
//                 }
//         }()
}

var (
        opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
                Name: "myapp_processed_ops_total",
                Help: "The total number of processed events",
        })
        testGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "my_gauge_test",
             Help:      "Test pour gauge",
        },
        []string{
            // Which user has requested the operation?
            "host",
            // Of what type is the operation?
            "playbook",
        })
        ansibleGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
             Name:      "ansible_stats",
             Help:      "Test pour gauge ansible_stats",
        },
        []string{
            // Which user has requested the operation?
            "playbook",
            // Of what type is the operation?
            "host",
            "status",
        })
)

func main() {

  recordMetrics()
  //Create a new instance of the foocollector and
  //register it with the prometheus client.
  // foo := newFooCollector()
  // prometheus.MustRegister(foo)

  h2 := func(w http.ResponseWriter, r *http.Request) {

    switch r.Method {
    case "POST":
        recordMetrics()



        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            panic(err)
        }
        // log.Println(string(body))



        var result map[string]interface{}
        json.Unmarshal([]byte(body), &result)


        global_custom_stats := result["global_custom_stats"].(map[string]interface{})

        var playbook = ""

        for key, value := range global_custom_stats {
          // Each value is an interface{} type, that is type asserted as a string
          playbook = value.(string)
          log.Println(string(key) + " " + value.(string))
        }

        stats := result["stats"].(map[string]interface{})

        log.Println(stats["k8s-0-edgeduck-prd"])

        for key, value := range stats {
          // Each value is an interface{} type, that is type asserted as a string
          log.Println(string(key) + " " + playbook + " : ")
          log.Println(stats[key])
          statsPerHost := stats[key].(map[string]interface{})
          log.Println(statsPerHost)

          for keyH, valueH := range statsPerHost {

            s := fmt.Sprintf("%f", valueH.(float64))
            log.Println(playbook + " ::: " + key + " :: " + keyH + " : " + s)
            ansibleGauge.WithLabelValues(playbook, key, keyH).Set(valueH.(float64))
          }

          if (value == "" ) {
            log.Println("hello")
          }
        }

        gaugeTestChangeValue()
        fmt.Fprintf(w, "POST RECEIVED")
    default:
        fmt.Fprintf(w, "Sorry, only POST methods are supported.")
    }
  }

  //This section will start the HTTP server and expose
  //any metrics on the /metrics endpoint.
  http.Handle("/metrics", promhttp.Handler())
  http.HandleFunc("/endpoint", h2)
  log.Info("Beginning to serve on port :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
