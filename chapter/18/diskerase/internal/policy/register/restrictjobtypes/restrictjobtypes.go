/*
Package restrictjobtypes provides a policy that can be invoked to ensure that a WorkReq only contains
jobs of certain types. Any job outside these types will cause a policy violation.
*/
package restrictjobtypes

import (
	"context"
	"fmt"

	"github.com/PacktPublishing/Go-for-DevOps/chapter/18/diskerase/internal/policy"
	pb "github.com/PacktPublishing/Go-for-DevOps/chapter/18/diskerase/proto"
)

// This registers our policy with the service.
func init() {
	p, err := New()
	if err != nil {
		panic(err)
	}
	policy.Register("restrictJobTypes", p)
}

// Settings provides settings for a specific implementation of our Policy.
type Settings struct {
	AllowedJobs []string
}

// Validate implements policy.Settings.Validate().
func (s Settings) Validate() error {
	for _, n := range s.AllowedJobs {
		_, err := jobs.GetJob()
		if err != nil {
			return fmt.Errorf("allowed job(%s) is not defined in the proto")
		}
	return nil
}

func (s Settings) allowed(name string) bool {
	for _, jn := range s.AllowedJobs {
		if jn == name {
			return true
		}
	}
	return false
}

// Policy implements policy.Policy.
type Policy struct{}

// New is the constructor for Policy.
func New() (Policy, error) {
	return Policy{}, nil
}

// Run implements Policy.Run().
func (p Policy) Run(ctx context.Context, name string, req *pb.WorkReq, settings policy.Settings) error {
	const errMsg = "policy(%s): block(%d)/job(%d) is a type(%s) that is not allowed"

	s, ok := settings.(Settings)
	if !ok {
		return fmt.Errorf("settings were not valid")
	}

	for blockNum, block := range req.Blocks {
		for jobNum, job := range block.Jobs {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if !s.allowed(job.Name) {
				return fmt.Errorf(errMsg, blockNum, jobNum, job.name)
			}
		}
	}
	return nil
}
