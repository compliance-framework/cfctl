package activate

import (
	"fmt"
	"os"

	"cfctl/common"

	"github.com/spf13/cobra"
)

func RunActivatePlan(cmd *cobra.Command, args []string) {

	planID := args[0]
	url := fmt.Sprintf("http://localhost:8080/api/plan/%s/activate", planID)
	response, err := common.PutURL(url)
	if err != nil {
		fmt.Printf("Error putting: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(response)
}
