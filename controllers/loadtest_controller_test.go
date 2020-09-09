package controllers

import (
	"context"
	"fmt"

	grpcv1 "github.com/grpc/test-infra/api/v1"
	"github.com/grpc/test-infra/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test Environment", func() {
	It("supports creation of load tests", func() {
		err := k8sClient.Create(context.Background(), newLoadTest())
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = Describe("Pod Creation", func() {
	var loadtest *grpcv1.LoadTest
	var defs *config.Defaults

	BeforeEach(func() {
		loadtest = newLoadTest()

		defs = &config.Defaults{
			DriverPool:  "drivers",
			WorkerPool:  "workers-8core",
			DriverPort:  10000,
			ServerPort:  10010,
			CloneImage:  "gcr.io/grpc-fake-project/test-infra/clone",
			ReadyImage:  "gcr.io/grpc-fake-project/test-infra/ready",
			DriverImage: "gcr.io/grpc-fake-project/test-infra/driver",
			Languages: []config.LanguageDefault{
				{
					Language:   "cxx",
					BuildImage: "l.gcr.io/google/bazel:latest",
					RunImage:   "gcr.io/grpc-fake-project/test-infra/cxx",
				},
				{
					Language:   "go",
					BuildImage: "golang:1.14",
					RunImage:   "gcr.io/grpc-fake-project/test-infra/go",
				},
				{
					Language:   "java",
					BuildImage: "java:jdk8",
					RunImage:   "gcr.io/grpc-fake-project/test-infra/java",
				},
			},
		}
	})

	Describe("newClientPod", func() {
		var component *grpcv1.Component

		BeforeEach(func() {
			component = &loadtest.Spec.Clients[0].Component
		})

		It("sets namespace to match loadtest", func() {
			namespace := "foobar"
			loadtest.Namespace = namespace

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Namespace).To(Equal(namespace))
		})

		It("sets component-name label", func() {
			name := "foo-bar-buzz"
			component.Name = &name

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.ComponentNameLabel]).To(Equal(name))
		})

		It("sets loadtest-role label to client", func() {
			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.RoleLabel]).To(Equal(config.ClientRole))
		})

		It("sets loadtest label", func() {
			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.LoadTestLabel]).To(Equal(loadtest.Name))
		})

		It("sets node selector for appropriate pool", func() {
			customPool := "custom-pool"
			component.Pool = &customPool

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.NodeSelector["pool"]).To(Equal(customPool))
		})

		It("sets clone init container", func() {
			cloneImage := "docker.pkg.github.com/grpc/test-infra/fake-image"
			component.Clone.Image = &cloneImage

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newCloneContainer(component.Clone)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets build init container", func() {
			buildImage := "docker.pkg.github.com/grpc/test-infra/fake-image"

			build := new(grpcv1.Build)
			build.Image = &buildImage
			build.Command = []string{"bazel"}
			build.Args = []string{"build", "//target"}
			component.Build = build

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newBuildContainer(component.Build)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets run container", func() {
			image := "golang:1.14"
			run := grpcv1.Run{
				Image:   &image,
				Command: []string{"go"},
				Args:    []string{"run", "main.go"},
			}
			component.Run = run

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newRunContainer(run)
			addDriverPort(&expectedContainer, defs.DriverPort)
			Expect(pod.Spec.Containers).To(ContainElement(expectedContainer))
		})

		It("disables retries", func() {
			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.RestartPolicy).To(Equal(corev1.RestartPolicyNever))
		})

		It("exposes a driver port", func() {
			pod, err := newClientPod(defs, loadtest, component)
			port := newContainerPort("driver", 10000)
			Expect(err).To(BeNil())
			Expect(pod.Spec.Containers[0].Ports).To(ContainElement(port))
		})

		It("sets driver port flag in run container args", func() {
			component.Run.Args = nil

			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			// TODO: Remove container lookup by index
			container := &pod.Spec.Containers[0]

			portFlag := fmt.Sprintf("--driver_port=%d", defs.DriverPort)
			Expect(container.Args).To(ContainElement(portFlag))
		})

		It("sets workspace volume", func() {
			pod, err := newClientPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			volume := newWorkspaceVolume()
			Expect(pod.Spec.Volumes).To(ContainElement(volume))
		})
	})

	Describe("newDriverPod", func() {
		var component *grpcv1.Component

		BeforeEach(func() {
			component = &loadtest.Spec.Driver.Component
		})

		It("sets namespace to match loadtest", func() {
			namespace := "foobar"
			loadtest.Namespace = namespace

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Namespace).To(Equal(namespace))
		})

		It("adds a scenario volume", func() {
			scenario := "example"
			loadtest.Spec.Scenarios[0] = grpcv1.Scenario{Name: scenario}

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedVolume := newScenarioVolume(scenario)
			Expect(expectedVolume).To(BeElementOf(pod.Spec.Volumes))
		})

		It("adds a scenario volume mount", func() {
			scenario := "example"
			loadtest.Spec.Scenarios[0] = grpcv1.Scenario{Name: scenario}

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			rc := &pod.Spec.Containers[0]
			expectedMount := newScenarioVolumeMount(scenario)
			Expect(expectedMount).To(BeElementOf(rc.VolumeMounts))
		})

		It("sets scenario file environment variable", func() {
			scenario := "example-scenario"
			loadtest.Spec.Scenarios[0] = grpcv1.Scenario{Name: scenario}

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			rc := &pod.Spec.Containers[0]
			expectedEnv := newScenarioFileEnvVar(scenario)
			Expect(expectedEnv).To(BeElementOf(rc.Env))
		})

		It("mounts GCP secrets", func() {
			// TODO: Add tests for mounting of GCP secrets
			Skip("complete this task when adding GCP secrets to pkg/defaults")
		})

		It("sets loadtest-role label to driver", func() {
			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.RoleLabel]).To(Equal(config.DriverRole))
		})

		It("sets loadtest label", func() {
			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.LoadTestLabel]).To(Equal(loadtest.Name))
		})

		It("sets component-name label", func() {
			name := "foo-bar-buzz"
			component.Name = &name

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.ComponentNameLabel]).To(Equal(name))
		})

		It("sets node selector for appropriate pool", func() {
			customPool := "custom-pool"
			component.Pool = &customPool

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.NodeSelector["pool"]).To(Equal(customPool))
		})

		It("sets clone init container", func() {
			cloneImage := "docker.pkg.github.com/grpc/test-infra/fake-image"
			component.Clone = new(grpcv1.Clone)
			component.Clone.Image = &cloneImage

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newCloneContainer(component.Clone)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets build init container", func() {
			buildImage := "docker.pkg.github.com/grpc/test-infra/fake-image"

			build := new(grpcv1.Build)
			build.Image = &buildImage
			build.Command = []string{"bazel"}
			build.Args = []string{"build", "//target"}
			component.Build = build

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newBuildContainer(component.Build)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets ready init container", func() {
			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newReadyContainer(defs, loadtest)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets run container", func() {
			scenario := "example"
			loadtest.Spec.Scenarios[0] = grpcv1.Scenario{Name: scenario}

			image := "golang:1.14"
			run := grpcv1.Run{
				Image:   &image,
				Command: []string{"go"},
				Args:    []string{"run", "main.go"},
			}
			component.Run = run

			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			rc := newRunContainer(run)
			addReadyInitContainer(defs, loadtest, &pod.Spec, &rc)

			rc.VolumeMounts = append(rc.VolumeMounts, newScenarioVolumeMount(scenario))
			rc.Env = append(rc.Env, newScenarioFileEnvVar(scenario))

			Expect(pod.Spec.Containers).To(ContainElement(rc))
		})

		It("disables retries", func() {
			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.RestartPolicy).To(Equal(corev1.RestartPolicyNever))
		})

		It("sets workspace volume", func() {
			pod, err := newDriverPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			volume := newWorkspaceVolume()
			Expect(pod.Spec.Volumes).To(ContainElement(volume))
		})
	})

	Describe("newServerPod", func() {
		var component *grpcv1.Component

		BeforeEach(func() {
			component = &loadtest.Spec.Servers[0].Component
		})

		It("sets namespace to match loadtest", func() {
			namespace := "foobar"
			loadtest.Namespace = namespace

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Namespace).To(Equal(namespace))
		})

		It("sets loadtest-role label to server", func() {
			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.RoleLabel]).To(Equal(config.ServerRole))
		})

		It("exposes a driver port", func() {
			pod, err := newServerPod(defs, loadtest, component)
			port := newContainerPort("driver", 10000)
			Expect(err).To(BeNil())
			Expect(pod.Spec.Containers[0].Ports).To(ContainElement(port))
		})

		It("sets driver port flag in run container args", func() {
			component.Run.Args = nil

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			// TODO: Remove container lookup by index
			container := &pod.Spec.Containers[0]

			portFlag := fmt.Sprintf("--driver_port=%d", defs.DriverPort)
			Expect(container.Args).To(ContainElement(portFlag))
		})

		It("exposes a server port", func() {
			pod, err := newServerPod(defs, loadtest, component)
			port := newContainerPort("server", 10010)
			Expect(err).To(BeNil())
			Expect(pod.Spec.Containers[0].Ports).To(ContainElement(port))
		})

		It("sets server port flag in run container args", func() {
			component.Run.Args = nil

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			// TODO: Remove container lookup by index
			container := &pod.Spec.Containers[0]

			portFlag := fmt.Sprintf("--server_port=%d", defs.ServerPort)
			Expect(container.Args).To(ContainElement(portFlag))
		})

		It("sets loadtest label", func() {
			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.LoadTestLabel]).To(Equal(loadtest.Name))
		})

		It("sets component-name label", func() {
			name := "foo-bar-buzz"
			component.Name = &name

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Labels[config.ComponentNameLabel]).To(Equal(name))
		})

		It("sets node selector for appropriate pool", func() {
			customPool := "custom-pool"
			component.Pool = &customPool

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.NodeSelector["pool"]).To(Equal(customPool))
		})

		It("sets clone init container", func() {
			cloneImage := "docker.pkg.github.com/grpc/test-infra/fake-image"
			component.Clone.Image = &cloneImage

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newCloneContainer(component.Clone)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets build init container", func() {
			buildImage := "docker.pkg.github.com/grpc/test-infra/fake-image"

			build := new(grpcv1.Build)
			build.Image = &buildImage
			build.Command = []string{"bazel"}
			build.Args = []string{"build", "//target"}
			component.Build = build

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newBuildContainer(component.Build)
			Expect(pod.Spec.InitContainers).To(ContainElement(expectedContainer))
		})

		It("sets run container", func() {
			image := "golang:1.14"
			run := grpcv1.Run{
				Image:   &image,
				Command: []string{"go"},
				Args:    []string{"run", "main.go"},
			}
			component.Run = run

			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			expectedContainer := newRunContainer(run)
			addDriverPort(&expectedContainer, defs.DriverPort)
			addServerPort(&expectedContainer, defs.ServerPort)
			Expect(pod.Spec.Containers).To(ContainElement(expectedContainer))
		})

		It("disables retries", func() {
			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod.Spec.RestartPolicy).To(Equal(corev1.RestartPolicyNever))
		})

		It("sets workspace volume", func() {
			pod, err := newServerPod(defs, loadtest, component)
			Expect(err).ToNot(HaveOccurred())

			volume := newWorkspaceVolume()
			Expect(pod.Spec.Volumes).To(ContainElement(volume))
		})
	})

	Describe("newCloneContainer", func() {
		var clone *grpcv1.Clone

		BeforeEach(func() {
			image := "docker.pkg.github.com/grpc/test-infra/clone"
			repo := "https://github.com/grpc/test-infra.git"
			gitRef := "master"

			clone = &grpcv1.Clone{
				Image:  &image,
				Repo:   &repo,
				GitRef: &gitRef,
			}
		})

		It("sets the name of the container", func() {
			container := newCloneContainer(clone)
			Expect(container.Name).To(Equal(cloneInitContainer))
		})

		It("sets workspace volume mount", func() {
			container := newCloneContainer(clone)
			volumeMount := newWorkspaceVolumeMount()
			Expect(container.VolumeMounts).To(ContainElement(volumeMount))
		})

		It("returns empty container given nil pointer", func() {
			clone = nil
			container := newCloneContainer(clone)
			Expect(container).To(Equal(corev1.Container{}))
		})

		It("sets clone image", func() {
			customImage := "debian:buster"
			clone.Image = &customImage

			container := newCloneContainer(clone)
			Expect(container.Image).To(Equal(customImage))
		})

		It("sets repo environment variable", func() {
			repo := "https://github.com/grpc/grpc.git"
			clone.Repo = &repo

			container := newCloneContainer(clone)
			Expect(container.Env).To(ContainElement(corev1.EnvVar{
				Name:  CloneRepoEnv,
				Value: repo,
			}))
		})

		It("sets git-ref environment variable", func() {
			gitRef := "master"
			clone.GitRef = &gitRef

			container := newCloneContainer(clone)
			Expect(container.Env).To(ContainElement(corev1.EnvVar{
				Name:  CloneGitRefEnv,
				Value: gitRef,
			}))
		})
	})

	Describe("newBuildContainer", func() {
		var build *grpcv1.Build

		BeforeEach(func() {
			image := "docker.pkg.github.com/grpc/test-infra/rust"

			build = &grpcv1.Build{
				Image:   &image,
				Command: nil,
				Args:    nil,
				Env:     nil,
			}
		})

		It("sets the name of the container", func() {
			container := newBuildContainer(build)
			Expect(container.Name).To(Equal(buildInitContainer))
		})

		It("sets workspace volume mount", func() {
			container := newBuildContainer(build)
			volumeMount := newWorkspaceVolumeMount()
			Expect(container.VolumeMounts).To(ContainElement(volumeMount))
		})

		It("sets workspace as working directory", func() {
			container := newBuildContainer(build)
			Expect(container.WorkingDir).To(Equal(workspaceMountPath))
		})

		It("returns empty container given nil pointer", func() {
			build = nil
			container := newBuildContainer(build)
			Expect(container).To(Equal(corev1.Container{}))
		})

		It("sets image", func() {
			customImage := "golang:latest"
			build.Image = &customImage

			container := newBuildContainer(build)
			Expect(container.Image).To(Equal(customImage))
		})

		It("sets command", func() {
			command := []string{"bazel"}
			build.Command = command

			container := newBuildContainer(build)
			Expect(container.Command).To(Equal(command))
		})

		It("sets args", func() {
			args := []string{"build", "//target"}
			build.Command = []string{"bazel"}
			build.Args = args

			container := newBuildContainer(build)
			Expect(container.Args).To(Equal(args))
		})

		It("sets environment variables", func() {
			env := []corev1.EnvVar{
				{Name: "EXPERIMENT", Value: "1"},
				{Name: "PROD", Value: "0"},
			}

			build.Env = env

			container := newBuildContainer(build)
			Expect(env[0]).To(BeElementOf(container.Env))
			Expect(env[1]).To(BeElementOf(container.Env))
		})
	})

	Describe("newRunContainer", func() {
		var run grpcv1.Run

		BeforeEach(func() {
			image := "docker.pkg.github.com/grpc/test-infra/fake-image"
			command := []string{"qps_worker"}

			run = grpcv1.Run{
				Image:   &image,
				Command: command,
			}
		})

		It("sets the name of the container", func() {
			container := newRunContainer(run)
			Expect(container.Name).To(Equal(runContainer))
		})

		It("sets workspace volume mount", func() {
			container := newRunContainer(run)
			volumeMount := newWorkspaceVolumeMount()
			Expect(container.VolumeMounts).To(ContainElement(volumeMount))
		})

		It("sets workspace as working directory", func() {
			container := newRunContainer(run)
			Expect(container.WorkingDir).To(Equal(workspaceMountPath))
		})

		It("sets image", func() {
			image := "golang:1.14"
			run.Image = &image

			container := newRunContainer(run)
			Expect(container.Image).To(Equal(image))
		})

		It("sets command", func() {
			command := []string{"go"}
			run.Command = command

			container := newRunContainer(run)
			Expect(container.Command).To(Equal(command))
		})

		It("sets args", func() {
			command := []string{"go"}
			args := []string{"run", "main.go"}
			run.Command = command
			run.Args = args

			container := newRunContainer(run)
			Expect(container.Args).To(Equal(args))
		})

		It("sets environment variables", func() {
			env := []corev1.EnvVar{
				{Name: "ENABLE_DEBUG", Value: "1"},
				{Name: "VERBOSE", Value: "1"},
			}

			run.Env = env

			container := newRunContainer(run)
			Expect(env[0]).To(BeElementOf(container.Env))
			Expect(env[1]).To(BeElementOf(container.Env))
		})
	})

	Describe("newWorkspaceVolumeMount", func() {
		It("grants read and write access", func() {
			volumeMount := newWorkspaceVolumeMount()
			Expect(volumeMount.ReadOnly).To(BeFalse())
		})
	})
})

var _ = Describe("checkMissingPods", func() {

	var currentLoadTest *grpcv1.LoadTest
	var allRunningPods *corev1.PodList
	var actualReturn *LoadTestMissing
	var expectedReturn *LoadTestMissing

	BeforeEach(func() {
		currentLoadTest = newLoadTestWithMultipleClientsAndServers()
		newListRef := corev1.PodList{Items: []corev1.Pod{}}
		allRunningPods = &newListRef
		expectedReturn = &LoadTestMissing{Clients: []grpcv1.Client{}, Servers: []grpcv1.Server{}}
	})

	Describe("no pods from the current load test is running", func() {
		BeforeEach(func() {
			for i := 0; i < len(currentLoadTest.Spec.Clients); i++ {
				expectedReturn.Clients = append(expectedReturn.Clients, currentLoadTest.Spec.Clients[i])
			}
			for i := 0; i < len(currentLoadTest.Spec.Servers); i++ {
				expectedReturn.Servers = append(expectedReturn.Servers, currentLoadTest.Spec.Servers[i])
			}
			expectedReturn.Driver = currentLoadTest.Spec.Driver
		})

		Context("no pod is running", func() {
			It("returns the full pod list from the current load test", func() {
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})

		Context("irrelevant pods are running", func() {
			It("returns the full pod list from the current load test", func() {
				allRunningPods.Items = append(allRunningPods.Items, createPodListWithIrrelevantPod().Items...)
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})
	})

	Describe("some of pods from the current load test is running", func() {

		BeforeEach(func() {
			allRunningPods.Items = append(allRunningPods.Items,
				corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "random-name",
						Labels: map[string]string{
							config.LoadTestLabel:      "test-loadtest-multiple-clients-and-servers",
							config.RoleLabel:          "server",
							config.ComponentNameLabel: "server-1",
						},
					},
				},
				corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "random-name",
						Labels: map[string]string{
							config.LoadTestLabel:      "test-loadtest-multiple-clients-and-servers",
							config.RoleLabel:          "client",
							config.ComponentNameLabel: "client-2",
						},
					},
				},
				corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "random-name",
						Labels: map[string]string{
							config.LoadTestLabel:      "test-loadtest-multiple-clients-and-servers",
							config.RoleLabel:          "driver",
							config.ComponentNameLabel: "driver-1",
						},
					},
				},
			)
			for i := 0; i < len(currentLoadTest.Spec.Clients); i++ {
				if *currentLoadTest.Spec.Clients[i].Name != "client-2" {
					expectedReturn.Clients = append(expectedReturn.Clients, currentLoadTest.Spec.Clients[i])
				}
			}

			for i := 0; i < len(currentLoadTest.Spec.Servers); i++ {
				if *currentLoadTest.Spec.Servers[i].Name != "server-1" {
					expectedReturn.Servers = append(expectedReturn.Servers, currentLoadTest.Spec.Servers[i])
				}
			}
		})

		Context("only pods from the current load test are running", func() {
			It("returns the list of pods missing from collection of running pods", func() {
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})

		Context("there are also irrelevant pods running", func() {
			It("returns the list of pods missing from collection of running pods", func() {
				allRunningPods.Items = append(allRunningPods.Items, createPodListWithIrrelevantPod().Items...)
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})
	})

	Describe("all of pods from the current load test is running", func() {

		BeforeEach(func() {
			allRunningPods = populatePodListWithCurrentLoadTestPod(currentLoadTest)
		})

		Context("only pods from the current load test are running", func() {
			It("returns a empty list", func() {
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})

		Context("there are also irrelevant pods running", func() {

			It("returns an empty list", func() {
				allRunningPods.Items = append(allRunningPods.Items, createPodListWithIrrelevantPod().Items...)
				actualReturn = checkMissingPods(currentLoadTest, allRunningPods)
				Expect(actualReturn.Clients).To(ConsistOf(expectedReturn.Clients))
				Expect(actualReturn.Servers).To(ConsistOf(expectedReturn.Servers))
				Expect(actualReturn.Driver).To(Equal(expectedReturn.Driver))
			})
		})
	})
})
