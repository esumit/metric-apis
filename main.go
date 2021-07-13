package main


import (
	"context"
	"github.com/esumit/metric-apis/pkg/config"
	"github.com/esumit/metric-apis/pkg/data"
	"github.com/esumit/metric-apis/pkg/httprqrs"
	"github.com/esumit/metric-apis/pkg/metric"
	"github.com/esumit/metric-apis/pkg/mw"
	negronilogrus "github.com/esumit/metric-apis/pkg/third-party/negroni-logrus"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	msc  config.MetricServerConfig
	mtc  config.MetricCollectionConfig
)
var c = make(chan os.Signal, 1)

var rootCmd = &cobra.Command{
	Use:   "metric-apis",
	Short: "Run metric-apis as a microservice",
	Run:   MetricApiService,
}



func MetricApiService(cmd *cobra.Command, args []string) {
	l:= data.NewListMetric()
	md:= metric.NewMetricApiDataManager(l,mtc.CollectionTime)
	mh:=metric.NewMetricApiRqHandler(md)
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(httprqrs.NotFoundHandler)
	
	metricAPIs := r.PathPrefix("/metrics").Subrouter()
	metricAPIs.HandleFunc("/{key}",mw.HttpRqRsMiddleware(mh.Save)).Methods("POST")
	metricAPIs.HandleFunc("/{key}/sum",mw.HttpRqRsMiddleware(mh.Get)).Methods("GET")
	
	h := cors.AllowAll().Handler(r)
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(log.StandardLogger(), "web"))
	n.UseHandler(h)
	
	srv := &http.Server{
		Addr:         msc.IPAddress + ":" + msc.Port,
		WriteTimeout: time.Second * time.Duration(msc.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(msc.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(msc.IdleTimeout),
		Handler:      n,
	}
	
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	
	signal.Notify(c, os.Interrupt)
	
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
	log.Infoln("metric service server shutdown ...")
	os.Exit(0)
}

func initConfig() {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	msc.Port = os.Getenv("SERVER_PORT")
	msc.IPAddress =os.Getenv("SERVER_IP_ADDRESS")
	msc.WriteTimeout, _ = strconv.Atoi(os.Getenv("HTTP_WRITE_TIMEOUT"))
	msc.ReadTimeout, _ = strconv.Atoi(os.Getenv("HTTP_READ_TIMEOUT"))
	msc.IdleTimeout, _ = strconv.Atoi(os.Getenv("HTTP_IDLE_TIMEOUT"))
	mtc.CollectionTime, _ = strconv.Atoi(os.Getenv("COLLECTION_TIMEOUT"))
	
	log.Println("Config Applied:")
	log.Println("Port: ", msc.Port)
	log.Println("IPAddress: ", msc.IPAddress)
	log.Println("HTTP WriteTimeout: ", msc.WriteTimeout)
	log.Println("HTTP ReadTimeout: ", msc.ReadTimeout)
	log.Println("HTTP IdleTimeout: ", msc.IdleTimeout)
	log.Println("Collection Timeout: ", mtc.CollectionTime)
	
	log.Println("All configs loaded")
}

func init() {
	cobra.OnInitialize(initConfig)
}

func MetricApis() {
	if err := rootCmd.Execute(); err != nil {
		log.Infoln(err)
		os.Exit(1)
	}
}

func main() {
	MetricApis()
}



