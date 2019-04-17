package sync

import (
	"encoding/json"

	"github.com/MathWebSearch/mwsapi/tema"

	"github.com/MathWebSearch/mwsapi/elasticutils"
)

// insertSegmentHarvests inserts a segment
func (proc *Process) insertSegmentHarvests(segment string) error {
	bulk := make(chan map[string]interface{})
	errChan := make(chan error)

	// start processing async
	go func() {
		e := processLinePairs(segment, true, func(_, contentLine string) (err error) {
			// unmarshal the content
			var content *tema.HarvestElement
			err = json.Unmarshal([]byte(contentLine), &content)
			if err != nil {
				return
			}

			content.Segment = segment

			var raw map[string]interface{}
			elasticutils.Repack(content, &raw)

			// store the content and put it into the db
			bulk <- raw

			return
		})

		// close both of the channel
		close(bulk)

		errChan <- e
		close(errChan)
	}()

	// run the insert and get the errors
	bulkError := elasticutils.CreateBulk(proc.connection.Client, proc.connection.Config.HarvestIndex, proc.connection.Config.HarvestType, bulk)
	parseError := <-errChan

	// return the parser error
	if parseError != nil {
		return parseError
	}

	// or the bulk error
	return bulkError
}
