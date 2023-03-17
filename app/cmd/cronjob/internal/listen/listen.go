package listen

import (
	"gozore-mall/app/cmd/cronjob/internal/svc"
	"gozore-mall/service/cronjob"
)

func Jobs(svcCtx *svc.ServiceContext) []cronjob.CronJob {
	var jobList []cronjob.CronJob

	//jobList = append(jobList, task.Testjob(svcCtx))

	return jobList
}
