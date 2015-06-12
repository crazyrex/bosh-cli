package installation_test

import (
	"errors"

	. "github.com/cloudfoundry/bosh-init/installation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"github.com/cloudfoundry/bosh-init/installation/blobextract/fakeblobextract"
	mock_install "github.com/cloudfoundry/bosh-init/installation/mocks"
	mock_registry "github.com/cloudfoundry/bosh-init/registry/mocks"

	biinstallmanifest "github.com/cloudfoundry/bosh-init/installation/manifest"
	bireljob "github.com/cloudfoundry/bosh-init/release/job"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	biproperty "github.com/cloudfoundry/bosh-utils/property"

	fakebiui "github.com/cloudfoundry/bosh-init/ui/fakes"
)

var _ = Describe("Installer", func() {
	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var (
		installationManifest      biinstallmanifest.Manifest
		mockJobRenderer           *mock_install.MockJobRenderer
		mockJobResolver           *mock_install.MockJobResolver
		mockPackageCompiler       *mock_install.MockPackageCompiler
		fakeExtractor             *fakeblobextract.FakeExtractor
		mockRegistryServerManager *mock_registry.MockServerManager

		logger boshlog.Logger

		packagesPath string
		installer    Installer
		target       Target
		installedJob InstalledJob
	)

	BeforeEach(func() {
		logger = boshlog.NewWriterLogger(boshlog.LevelDebug, GinkgoWriter, GinkgoWriter)

		mockJobRenderer = mock_install.NewMockJobRenderer(mockCtrl)
		mockJobResolver = mock_install.NewMockJobResolver(mockCtrl)
		mockPackageCompiler = mock_install.NewMockPackageCompiler(mockCtrl)
		fakeExtractor = &fakeblobextract.FakeExtractor{}
		mockRegistryServerManager = mock_registry.NewMockServerManager(mockCtrl)

		packagesPath = "/path/to/installed/packages"
		target = NewTarget("fake-installation-path")
		installationManifest = biinstallmanifest.Manifest{
			Name:       "fake-installation-name",
			Properties: biproperty.Map{},
		}
		renderedCPIJob := NewRenderedJobRef("cpi", "fake-release-job-fingerprint", "fake-rendered-job-blobstore-id", "fake-rendered-job-blobstore-id")
		installedJob = NewInstalledJob(renderedCPIJob, "/extracted-release-path/cpi")
	})

	JustBeforeEach(func() {
		installer = NewInstaller(
			target,
			mockJobRenderer,
			mockJobResolver,
			mockPackageCompiler,
			packagesPath,
			fakeExtractor,
			mockRegistryServerManager,
			logger,
		)
	})

	Describe("Install", func() {
		var (
			fakeStage *fakebiui.FakeStage

			expectedResolveJobsFrom     *gomock.Call
			expectedPackageCompilerFrom *gomock.Call
			expectedRenderAndUploadFrom *gomock.Call
		)

		BeforeEach(func() {
			fakeStage = fakebiui.NewFakeStage()
		})

		JustBeforeEach(func() {
			ref := CompiledPackageRef{
				Name:        "fake-release-package-name",
				Version:     "fake-release-package-fingerprint",
				BlobstoreID: "fake-compiled-package-blobstore-id",
				SHA1:        "fake-compiled-package-blobstore-id",
			}
			compiledPackages := []CompiledPackageRef{ref}

			releaseJobs := []bireljob.Job{}
			renderedJobRefs := []RenderedJobRef{installedJob.RenderedJobRef}
			expectedResolveJobsFrom = mockJobResolver.EXPECT().From(installationManifest).Return(releaseJobs, nil).AnyTimes()
			expectedPackageCompilerFrom = mockPackageCompiler.EXPECT().For(releaseJobs, fakeStage).Return(compiledPackages, nil).AnyTimes()
			expectedRenderAndUploadFrom = mockJobRenderer.EXPECT().RenderAndUploadFrom(installationManifest, releaseJobs, fakeStage).Return(renderedJobRefs, nil).AnyTimes()
		})

		It("compiles and installs the jobs' packages", func() {
			_, err := installer.Install(installationManifest, fakeStage)
			Expect(err).NotTo(HaveOccurred())
		})

		It("installs the rendered jobs", func() {
			_, err := installer.Install(installationManifest, fakeStage)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the installation", func() {
			installation, err := installer.Install(installationManifest, fakeStage)
			Expect(err).NotTo(HaveOccurred())
			Expect(installation.Target().JobsPath()).To(Equal(target.JobsPath()))
		})
	})

	Describe("Cleanup", func() {
		var installation Installation

		BeforeEach(func() {
			installation = NewInstallation(
				target,
				installedJob,
				installationManifest,
				mockRegistryServerManager,
			)
		})

		It("cleans up installed jobs", func() {
			err := installer.Cleanup(installation)
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeExtractor.CleanupCallCount()).To(Equal(1))

			blobstoreID := fakeExtractor.CleanupArgsForCall(0)
			Expect(blobstoreID).To(Equal(installedJob.BlobstoreID))
		})

		It("returns errors when cleaning up installed jobs", func() {
			fakeExtractor.CleanupReturns(errors.New("nope"))

			err := installer.Cleanup(installation)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("nope"))
		})
	})
})
