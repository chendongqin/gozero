package main

import (
	"flag"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"gozore-mall/app/cmd/cronjob/internal/config"
	"gozore-mall/app/cmd/cronjob/internal/listen"
	"gozore-mall/app/cmd/cronjob/internal/svc"
	"os"
)

var fileName, _ = os.Getwd()
var configFile = flag.String("f", fileName+"etc/cronjob.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	// log、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	cronJob := cron.New(cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.DowOptional|cron.Descriptor)), cron.WithChain())
	jobList := listen.Jobs(ctx)
	for _, j := range jobList {
		_, _ = cronJob.AddJob(j.Spec(), j)
	}
	cronJob.Start()
	defer cronJob.Stop()
	select {}
}
