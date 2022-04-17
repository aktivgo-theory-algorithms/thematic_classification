package errors

import "sync"

func HandleError(wg *sync.WaitGroup, wgErrors chan error) error {
	for {
		select {
		default:
			wg.Wait()
			return nil
		case err := <-wgErrors:
			return err
		}
	}
}
