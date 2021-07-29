package v1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	v2 "github.com/darkowlzz/crd-api-conversion-demo/api/v2"
)

func (src *CronJob) ConvertTo(dstRaw conversion.Hub) error {
	cronjoblog.Info("CronJob.v1 ConvertTo called", "resource", src.Name)
	dst := dstRaw.(*v2.CronJob)

	dst.ObjectMeta = src.ObjectMeta

	return nil
}

func (dst *CronJob) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v2.CronJob)

	cronjoblog.Info("CronJob.v1 ConvertFrom called", "resource", src.Name)

	dst.ObjectMeta = src.ObjectMeta

	return nil
}
