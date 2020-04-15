package ami

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kl "k8s.io/apimachinery/pkg/labels"
)

func (k *kubeImpl) Log(ns, service, tailLines, sinceSeconds string) ([]byte, error) {
	deploy, err := k.cli.app.Deployments(ns).Get(service, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if deploy == nil {
		return nil, fmt.Errorf("service doesn't exist")
	}
	ls := kl.Set{}
	selector := deploy.Spec.Selector.MatchLabels
	err = copier.Copy(&ls, &selector)
	pods, err := k.cli.core.Pods(ns).List(metav1.ListOptions{
		LabelSelector: ls.String(),
	})
	if pods == nil || len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pod or more than one pod exists")
	}

	cli := k.cli.core.
		RESTClient().
		Get().
		Namespace(ns).
		Name(pods.Items[0].Name).
		Resource("pods").
		SubResource("log").
		Param("follow", k.conf.LogConfig.Follow).
		Param("previous", k.conf.LogConfig.Previous).
		Param("timestamps", k.conf.LogConfig.TimeStamps)

	if tailLines != "" {
		cli = cli.Param("tailLines", tailLines)
	}
	if sinceSeconds != "" {
		cli = cli.Param("sinceSeconds", sinceSeconds)
	}

	buf := new(bytes.Buffer)
	readCloser, err := cli.Stream()
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	if _, err = io.Copy(buf, readCloser); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
