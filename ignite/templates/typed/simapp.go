package typed

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
)

func ModuleSimulationMsgModify(
	replacer placeholder.Replacer,
	content,
	moduleName string,
	typeName multiformatname.Name,
	msgs ...string,
) string {
	if len(msgs) == 0 {
		msgs = append(msgs, "")
	}
	for _, msg := range msgs {
		// simulation constants
		templateConst := `opWeightMsg%[2]v%[3]v = "op_weight_msg_%[4]v"
	// TODO: Determine the simulation weight value
	defaultWeightMsg%[2]v%[3]v int = 100

	%[1]v`
		replacementConst := fmt.Sprintf(templateConst, PlaceholderSimappConst, msg, typeName.UpperCamel, typeName.Snake)
		content = replacer.Replace(content, PlaceholderSimappConst, replacementConst)

		// simulation operations
		templateOp := `var weightMsg%[2]v%[3]v int
	simState.AppParams.GetOrGenerate(opWeightMsg%[2]v%[3]v, &weightMsg%[2]v%[3]v, nil,
		func(_ *rand.Rand) {
			weightMsg%[2]v%[3]v = defaultWeightMsg%[2]v%[3]v
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsg%[2]v%[3]v,
		%[4]vsimulation.SimulateMsg%[2]v%[3]v(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	%[1]v`
		replacementOp := fmt.Sprintf(templateOp, PlaceholderSimappOperation, msg, typeName.UpperCamel, moduleName)
		content = replacer.Replace(content, PlaceholderSimappOperation, replacementOp)

		if strings.Contains(content, PlaceholderSimappOperationMsg) {
			templateOpMsg := `simulation.NewWeightedProposalMsg(
	opWeightMsg%[2]v%[3]v,
	defaultWeightMsg%[2]v%[3]v,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		%[4]vsimulation.SimulateMsg%[2]v%[3]v(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
%[1]v`
			replacementOpMsg := fmt.Sprintf(templateOpMsg, PlaceholderSimappOperationMsg, msg, typeName.UpperCamel, moduleName)
			content = replacer.Replace(content, PlaceholderSimappOperationMsg, replacementOpMsg)
		}
	}
	return content
}
