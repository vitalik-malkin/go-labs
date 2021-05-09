package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	smpl_driver "github.com/vitalik-malkin/go-labs/backoop/internal/driver/simple"
	plan "github.com/vitalik-malkin/go-labs/backoop/pkg/plan"
	plan_json "github.com/vitalik-malkin/go-labs/backoop/pkg/plan/json"
)

var (
	makeCmd = &cobra.Command{
		Use:          "make",
		Short:        "Make backup",
		SilenceUsage: true,
		RunE:         makeCmdExec,
	}
	makeCmdConfigVar = makeCmdConfig{}
)

func makeCmdExec(cmd *cobra.Command, args []string) error {
	ctx, cfg := cmd.Context(), makeCmdConfigVar

	// load plan.
	var loadPlanF = func() (plan.Plan, error) {
		planFile, err := os.Open(cfg.PlanFile)
		if err != nil {
			return nil, fmt.Errorf("error while opening plan file; err: %w", err)
		}
		defer planFile.Close()
		plan, err := plan_json.Load(ctx, planFile)
		if err != nil {
			return nil, fmt.Errorf("error while loading plan; err: %w", err)
		}
		return plan, nil
	}

	plan, err := loadPlanF()
	if err != nil {
		return err
	}

	driver := smpl_driver.New()
	_, err = driver.Exec(ctx, plan, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Version: %s (%s, %s)\n", "v1.0", runtime.GOOS, runtime.GOARCH)
	return nil
}
