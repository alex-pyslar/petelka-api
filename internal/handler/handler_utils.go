package handler

import (
	"strconv"

	"github.com/alex-pyslar/online-store/internal/logger"
)

// getIDFromVars is a utility function to extract ID from URL variables.
func getIDFromVars(vars map[string]string, key string, log *logger.Logger, requestID string) (int, error) {
	idStr := vars[key]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Invalid ID format for key %s, request_id: %s: %v", key, requestID, err)
		return 0, err
	}
	return id, nil
}
