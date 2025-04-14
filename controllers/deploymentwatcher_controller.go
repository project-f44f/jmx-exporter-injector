/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"regexp"
	"strings"

	jmxv1 "github.com/project-f44f/jmx-exporter-injector/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentWatcherReconciler reconciles a DeploymentWatcher object
type DeploymentWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploymentWatcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DeploymentWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var deploy appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deploy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 获取所有 JmxExporter 对象
	var jmxList jmxv1.JmxExporterList
	if err := r.List(ctx, &jmxList); err != nil {
		logger.Error(err, "无法列出 JmxExporter")
		return ctrl.Result{}, err
	}
	for _, jmx := range jmxList.Items {
		val := jmx.Spec.DeploymentAnnotationValue

		if deploy.Annotations != nil {
			if deployVal, ok := deploy.Annotations["jmx.f44f.com/jmx-exporter"]; ok && deployVal == val {
				logger.Info("匹配到 JmxExporter 策略",
					"deployment", deploy.Name,
					"namespace", deploy.Namespace,
					"annotationValue", val,
					"jmxexporter", jmx.Name,
				)
				container := &deploy.Spec.Template.Spec.Containers[0]
				// 如果需要注入
				if jmx.Spec.EnableInjection {
					// 添加挂载路径
					mountPathExists := false
					for _, vm := range container.VolumeMounts {
						if vm.MountPath == "/khaos/" {
							mountPathExists = true
							logger.Info("挂载路径 已存在")
							break
						}
					}
					if !mountPathExists {
						container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
							Name:      "preset-file",
							MountPath: "/khaos/",
						})
					}
					// 添加 initContainer
					initExists := false
					for _, ic := range deploy.Spec.Template.Spec.InitContainers {
						if ic.Name == "init-preset-file" {
							initExists = true
							logger.Info("initcontainer 已存在")
							break
						}
					}
					if !initExists {
						deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, corev1.Container{
							Name:            "init-preset-file",
							Image:           "registry2-qingdao.cosmoplat.com/64_paas/init_external:latest",
							ImagePullPolicy: corev1.PullAlways,
							Command:         []string{"sh", "-c", "cp -r /usr/local/external/jmx /khaos"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("200m"),
									corev1.ResourceMemory: resource.MustParse("1G"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1G"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "preset-file",
									MountPath: "/khaos/",
								},
							},
						})
					}
					// 添加 volume
					volumeExists := false
					for _, vol := range deploy.Spec.Template.Spec.Volumes {
						if vol.Name == "preset-file" {
							volumeExists = true
							logger.Info("volume 已存在")
							break
						}
					}
					if !volumeExists {
						deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, corev1.Volume{
							Name: "preset-file",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						})
					}
				}
				// 添加环境变量
				var port string
				if val, ok := deploy.Annotations["jmx.f44f.com/jmx-exporter-port"]; ok {
					port = val
				} else {
					port = "8088"
				}
				logger.Info("检测到port为" + port)
				// 定义 JAVA_TOOL_OPTIONS 环境变量的值
				name := "JAVA_TOOL_OPTIONS"
				value := "-javaagent:/khaos/jmx/jmx_prometheus_javaagent-1.0.1.jar=" + port + ":/khaos/jmx/prometheus-jmx-config.yaml"
				// 遍历容器的环境变量，检查是否存在 JAVA_TOOL_OPTIONS
				envExists := false
				for i, env := range container.Env {
					if env.Name == name {
						envExists = true
						envValue := env.Value
						// 检查值是否完全一致
						if envValue == value {
							logger.Info("环境变量的值和目标值完全一致")
							// 如果环境变量的值和目标值完全一致，什么都不做
							break
						}
						// 检查是否只有 port 不一样
						// 如果存在 -javaagent，并且只有 port 部分不同，则替换 port
						if strings.Contains(envValue, "-javaagent:/khaos/jmx/jmx_prometheus_javaagent-1.0.1.jar") {
							// 使用正则表达式替换掉旧的 port 部分
							logger.Info("port 不一样")
							pattern := `-javaagent:/khaos/jmx/jmx_prometheus_javaagent-1\.0\.1\.jar=[^ ]*/khaos/jmx/prometheus-jmx-config\.yaml`
							re := regexp.MustCompile(pattern)
							updatedValue := re.ReplaceAllString(envValue, value)
							container.Env[i].Value = updatedValue
						} else {
							logger.Info("完全不一样，拼接")
							// 如果环境变量值完全不同，直接在后面拼接新值
							container.Env[i].Value = envValue + " " + value
						}

						// 退出循环
						break
					}
				}
				// 如果环境变量不存在，添加新的环境变量
				if !envExists {
					container.Env = append(container.Env, corev1.EnvVar{
						Name:  name,
						Value: value,
					})
				}
				if err := r.Update(ctx, &deploy); err != nil {
					logger.Error(err, "更新 Deployment 失败")
					return ctrl.Result{}, err
				}
				logger.Info("Deployment 修改成功", "name", deploy.Name)
				break
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
