package elasticsync

import (
	"encoding/json"

	"github.com/MathWebSearch/mwsapi/result"

	"github.com/MathWebSearch/mwsapi/utils"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// insertSegmentHarvests inserts a segment
func (proc *Process) insertSegmentHarvests(segment string) error {
	bulk := make(chan interface{})
	errChan := make(chan error)

	// start processing async
	go func() {
		e := utils.ProcessLinePairs(segment, true, func(_, contentLine string) (err error) {
			// unmarshal the content
			var content *result.ElasticElement
			err = json.Unmarshal([]byte(contentLine), &content)
			if err != nil {
				return
			}

			content.Segment = segment

			// store the content and put it into the db
			bulk <- content

			return
		})

		// close both of the channel
		close(bulk)

		errChan <- e
		close(errChan)
	}()

	// run the insert and get the errors
	bulkError := elasticutils.CreateBulk(proc.conn.Client, proc.conn.Config.HarvestIndex, proc.conn.Config.HarvestType, bulk)
	parseError := <-errChan

	// return the parser error
	if parseError != nil {
		return parseError
	}

	// or the bulk error
	return bulkError
}
