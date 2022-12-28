package mainhandler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/kubescape/operator/utils"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (actionHandler *ActionHandler) updateRegistryScanCronJob() error {
	jobParams := actionHandler.command.GetCronJobParams()
	if jobParams == nil {
		glog.Infof("updateRegistryScanCronJob: failed to get jobParams")
		return fmt.Errorf("failed to get failed to get jobParams")
	}

	jobTemplateObj, err := actionHandler.k8sAPI.KubernetesClient.BatchV1().CronJobs(utils.Namespace).Get(context.Background(), jobParams.JobName, metav1.GetOptions{})
	if err != nil {
		glog.Infof("updateRegistryScanCronJob: failed to get cronjob: %s", jobParams.JobName)
		return err
	}

	jobTemplateObj.Spec.Schedule = getCronTabSchedule(actionHandler.command)
	if jobTemplateObj.Spec.JobTemplate.Spec.Template.Annotations == nil {
		jobTemplateObj.Spec.JobTemplate.Spec.Template.Annotations = make(map[string]string)
	}

	jobTemplateObj.Spec.JobTemplate.Spec.Template.Annotations[armotypes.CronJobTemplateAnnotationUpdateJobIDDeprecated] = actionHandler.command.JobTracking.JobID // deprecated
	jobTemplateObj.Spec.JobTemplate.Spec.Template.Annotations[armotypes.CronJobTemplateAnnotationUpdateJobID] = actionHandler.command.JobTracking.JobID

	_, err = actionHandler.k8sAPI.KubernetesClient.BatchV1().CronJobs(utils.Namespace).Update(context.Background(), jobTemplateObj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	glog.Infof("updateRegistryScanCronJob: cronjob: %v updated successfully", jobParams.JobName)
	return nil

}

func (actionHandler *ActionHandler) setRegistryScanCronJob(sessionObj *utils.SessionObj) error {
	// 1 - If we have credentials, create secret with it.
	// 2 - Create configmap with command to trigger operator. Command includes secret name (if there were credentials).
	// 3 - Create cronjob which will send request to operator to trigger scan using the configmap data.

	if getCronTabSchedule(sessionObj.Command) == "" {
		return fmt.Errorf("schedule cannot be empty")
	}

	registryScan, err := actionHandler.loadRegistryScan(sessionObj)
	if err != nil {
		glog.Errorf("in parseRegistryCommand: error: %v", err.Error())
		sessionObj.Reporter.SetDetails("loadRegistryScan")
		return fmt.Errorf("scanRegistries failed with err %v", err)
	}

	// name is registryScanConfigmap name + random string - configmap and cronjob
	nameSuffix := rand.NewSource(time.Now().UnixNano()).Int63()
	name := fixK8sCronJobNameLimit(fmt.Sprintf("%s-%d", registryScanConfigmap, nameSuffix))
	if registryScan.registryInfo.AuthMethod.Type != "public" {
		err = registryScan.createTriggerRequestSecret(actionHandler.k8sAPI, name, registryScan.registryInfo.RegistryName)
		if err != nil {
			glog.Infof("setRegistryScanCronJob: error creating configmap : %s", err.Error())
			return err
		}
	}

	// create configmap with POST data to trigger websocket
	err = registryScan.createTriggerRequestConfigMap(actionHandler.k8sAPI, name, registryScan.registryInfo.RegistryName, sessionObj.Command)
	if err != nil {
		glog.Infof("setRegistryScanCronJob: error creating configmap : %s", err.Error())
		return err
	}

	err = registryScan.createTriggerRequestCronJob(actionHandler.k8sAPI, name, registryScan.registryInfo.RegistryName, sessionObj.Command)
	if err != nil {
		glog.Infof("setRegistryScanCronJob: error creating cronjob : %s", err.Error())
		return err
	}

	glog.Infof("setRegistryScanCronJob: cronjob: %s created successfully", name)
	return err
}

func (actionHandler *ActionHandler) deleteRegistryScanCronJob() error {
	jobParams := actionHandler.command.GetCronJobParams()
	if jobParams == nil {
		glog.Infof("updateRegistryScanCronJob: failed to get jobParams")
		return fmt.Errorf("failed to get jobParams")
	}

	return actionHandler.deleteCronjob(jobParams.JobName)
}
