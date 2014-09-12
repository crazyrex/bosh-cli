package fakes

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bmdepl "github.com/cloudfoundry/bosh-micro-cli/deployment"
	bmrel "github.com/cloudfoundry/bosh-micro-cli/release"
	bmtestutils "github.com/cloudfoundry/bosh-micro-cli/testutils"
)

type compileInput struct {
	Release    bmrel.Release
	Deployment bmdepl.Deployment
}

type compileOutput struct {
	err error
}

type FakeReleaseCompiler struct {
	CompileInputs   []compileInput
	compileBehavior map[string]compileOutput
}

func NewFakeReleaseCompiler() *FakeReleaseCompiler {
	return &FakeReleaseCompiler{
		CompileInputs:   []compileInput{},
		compileBehavior: map[string]compileOutput{},
	}
}

func (f *FakeReleaseCompiler) Compile(release bmrel.Release, deployment bmdepl.Deployment) error {
	input := compileInput{Release: release, Deployment: deployment}
	f.CompileInputs = append(f.CompileInputs, input)
	inputString, marshalErr := bmtestutils.MarshalToString(input)
	if marshalErr != nil {
		return bosherr.WrapError(marshalErr, "Marshaling Find input")
	}

	output, found := f.compileBehavior[inputString]

	if found {
		return output.err
	}

	return fmt.Errorf("Unsupported Input: Compile('%#v', '%#v'). Available inputs are '%#v'", release, deployment, f.compileBehavior)
}

func (f *FakeReleaseCompiler) SetCompileBehavior(release bmrel.Release, deployment bmdepl.Deployment, err error) error {
	input := compileInput{Release: release, Deployment: deployment}
	inputString, marshalErr := bmtestutils.MarshalToString(input)
	if marshalErr != nil {
		return bosherr.WrapError(marshalErr, "Marshaling Find input")
	}

	f.compileBehavior[inputString] = compileOutput{err: err}
	return nil
}
