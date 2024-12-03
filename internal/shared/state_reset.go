package shared

import "context"

func (s State) ResetDB() error {
	err := s.DB.ResetGoalieStats(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.ResetSkaterGameStats(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.ResetShots(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.ResetPenSummaries(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.ResetScoringSummaries(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.ResetGameResults(context.Background())
	if err != nil {
		return err
	}
	
	return nil
}	