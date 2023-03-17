package cronjob

type CronJob interface {
	Run()
	Spec() string
}
