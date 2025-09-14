package sm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	smmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sm/models"
	"github.com/mitchellh/mapstructure"
)

const (
	sessionsURL          = "/api/sessions"
	sessionURL           = "/api/sessions/%s"
	sessionActivitiesURL = "/api/sessions/%s/activities"
)

// ArkSMPage represents a page of ArkSMSession items.
type ArkSMPage = common.ArkPage[smmodels.ArkSMSession]

// ArkSMActivitiesPage represents a page of ArkSMSessionActivity items.
type ArkSMActivitiesPage = common.ArkPage[smmodels.ArkSMSessionActivity]

// ArkSMService is the implementation of the ArkSMService interface.
type ArkSMService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
	env     commonmodels.AwsEnv
}

// NewArkSMService creates a new instance of ArkSMService.
func NewArkSMService(authenticators ...auth.ArkAuth) (*ArkSMService, error) {
	SMService := &ArkSMService{}
	var SMServiceInterface services.ArkService = SMService
	baseService, err := services.NewArkBaseService(SMServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "sessionmonitoring", ".", "", SMService.refreshSMAuth)
	if err != nil {
		return nil, err
	}
	SMService.client = client
	SMService.ispAuth = ispAuth
	SMService.ArkBaseService = baseService
	return SMService, nil
}

func (s *ArkSMService) refreshSMAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// searchParamsFromFilter private function converts an ArkSMSessionsFilter to a map of search parameters
func (s *ArkSMService) searchParamsFromFilter(sessionsFilter *smmodels.ArkSMSessionsFilter) map[string]string {
	return map[string]string{
		"search": sessionsFilter.Search,
	}
}

// callListSessions private function that retrieves a list of sessions, parameters can be passed to filter the results.
func (s *ArkSMService) callListSessions(params map[string]string) (*smmodels.ArkSMSessions, error) {
	if params == nil {
		params = make(map[string]string)
	}
	response, err := s.client.Get(context.Background(), sessionsURL, params)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list sessions - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	sessionsJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var sessions smmodels.ArkSMSessions
	err = mapstructure.Decode(sessionsJSON, &sessions)
	if err != nil {
		return nil, err
	}
	return &sessions, nil
}

// callListSessionActivities private function that retrieves a list of activities for a session, parameters can be passed to filter the results.
func (s *ArkSMService) callListSessionActivities(sessionID string, params map[string]string) (*smmodels.ArkSMSessionActivities, error) {
	if params == nil {
		params = make(map[string]string)
	}
	response, err := s.client.Get(context.Background(), fmt.Sprintf(sessionActivitiesURL, sessionID), params)
	if err != nil {
		s.Logger.Error("failed to list session activities: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list session activities - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	sessionActivitiesJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var sessionActivities smmodels.ArkSMSessionActivities
	err = mapstructure.Decode(sessionActivitiesJSON, &sessionActivities)
	if err != nil {
		return nil, err
	}
	return &sessionActivities, nil
}

// listPagedSessions private function that retrieves a list of sessions, parameters can be passed to filter the results.
func (s *ArkSMService) listPagedSessions(params map[string]string) (<-chan *ArkSMPage, error) {
	results := make(chan *ArkSMPage)
	if params == nil {
		params = make(map[string]string)
	}
	offset := 0
	go func() {
		defer close(results)
		for {
			sessionsResponse, err := s.callListSessions(params)
			if err != nil {
				s.Logger.Error("failed to list sessions: %v", err)
				return
			}
			if sessionsResponse.ReturnedCount == 0 {
				break
			}
			sessions := make([]*smmodels.ArkSMSession, len(sessionsResponse.Sessions))
			for i := range sessionsResponse.Sessions {
				sessions[i] = &sessionsResponse.Sessions[i]
			}
			results <- &ArkSMPage{Items: sessions}
			offset += sessionsResponse.ReturnedCount
			params["offset"] = strconv.Itoa(offset)
		}
	}()
	return results, nil
}

// listActivities private function that retrieves the activities by session ID
// parameters can be passed to filter the results.
func (s *ArkSMService) listPagedSessionActivities(sessionID string) (<-chan *ArkSMActivitiesPage, error) {
	results := make(chan *ArkSMActivitiesPage)
	params := make(map[string]string)
	offset := 0
	go func() {
		defer close(results)
		for {
			sessionActivitiesResponse, err := s.callListSessionActivities(sessionID, params)
			if err != nil {
				s.Logger.Error("failed to list session activities: %v", err)
				return
			}
			if sessionActivitiesResponse.ReturnedCount == 0 {
				break
			}
			activities := make([]*smmodels.ArkSMSessionActivity, len(sessionActivitiesResponse.Activities))
			for i := range sessionActivitiesResponse.Activities {
				activities[i] = &sessionActivitiesResponse.Activities[i]
			}
			results <- &ArkSMActivitiesPage{Items: activities}
			offset += sessionActivitiesResponse.ReturnedCount
			params["offset"] = strconv.Itoa(offset)

		}
	}()
	return results, nil
}

// ListSessions retrieves a list of sessions
func (s *ArkSMService) ListSessions() (<-chan *ArkSMPage, error) {
	return s.listPagedSessions(nil)
}

// CountSessions retrieves the count of sessions on the last 24 hours
func (s *ArkSMService) CountSessions() (int, error) {
	sessions, err := s.callListSessions(nil)
	if err != nil {
		s.Logger.Error("failed to count sessions: %v", err)
		return 0, err
	}
	return sessions.ReturnedCount, err
}

// ListSessionsBy retrieves a list of sessions and applies an optional filter.
func (s *ArkSMService) ListSessionsBy(filter *smmodels.ArkSMSessionsFilter) (<-chan *ArkSMPage, error) {
	return s.listPagedSessions(s.searchParamsFromFilter(filter))
}

// CountSessionsBy retrieves the count of sessions on the last 24 hours and applies an optional filter.
func (s *ArkSMService) CountSessionsBy(filter *smmodels.ArkSMSessionsFilter) (int, error) {
	sessions, err := s.callListSessions(s.searchParamsFromFilter(filter))
	if err != nil {
		s.Logger.Error("failed to count sessions: %v", err)
		return 0, err
	}
	return sessions.FilteredCount, err
}

// ListSessionActivities retrieves the activities of a session by its ID
func (s *ArkSMService) ListSessionActivities(sessionActivities *smmodels.ArkSIASMGetSessionActivities) (<-chan *ArkSMActivitiesPage, error) {
	return s.listPagedSessionActivities(sessionActivities.SessionID)
}

// CountSessionActivities retrieves the count all session activities by session id
func (s *ArkSMService) CountSessionActivities(activities *smmodels.ArkSIASMGetSessionActivities) (int, error) {
	sessionActivities, err := s.callListSessionActivities(activities.SessionID, nil)
	if err != nil {
		s.Logger.Error("failed counting session activities: %v", err)
		return 0, err
	}
	return sessionActivities.ReturnedCount, err
}

// ListSessionActivitiesBy retrieves the activities of a session by its ID and applies an optional filter.
func (s *ArkSMService) ListSessionActivitiesBy(filter *smmodels.ArkSMSessionActivitiesFilter) (<-chan *ArkSMActivitiesPage, error) {
	pagedSessionActivities, err := s.listPagedSessionActivities(filter.SessionID)
	if err != nil {
		s.Logger.Error("failed to list session activities: %v", err)
		return nil, err
	}
	out := make(chan *ArkSMActivitiesPage)

	go func() {
		defer close(out)

		for page := range pagedSessionActivities {
			filteredItems := make([]*smmodels.ArkSMSessionActivity, 0, len(page.Items))

			for _, activity := range page.Items {
				if filter.CommandContains == "" || strings.Contains(activity.Command, filter.CommandContains) {
					filteredItems = append(filteredItems, activity)
				}
			}

			out <- &ArkSMActivitiesPage{
				Items: filteredItems,
			}
		}
	}()

	return out, nil
}

// CountSessionActivitiesBy retrieves the count all session activities by session id and applies an optional filter.
func (s *ArkSMService) CountSessionActivitiesBy(filter *smmodels.ArkSMSessionActivitiesFilter) (int, error) {
	pagedSessionActivities, err := s.ListSessionActivitiesBy(filter)
	if err != nil {
		s.Logger.Error("failed counting session activities: %v", err)
		return 0, err
	}
	count := 0
	for page := range pagedSessionActivities {
		count += len(page.Items)
	}
	return count, err
}

// SessionsStats retrieves the session statistics for the SM service.
func (s *ArkSMService) SessionsStats() (*smmodels.ArkSMSessionsStats, error) {
	s.Logger.Info("Calculating sessions stats for the last 30 days")
	startTimeFrom := time.Now().AddDate(0, 0, -30).UTC().Format("2006-01-02T15:04:05Z")

	filter := smmodels.ArkSMSessionsFilter{
		Search: fmt.Sprintf("startTime ge %s", startTimeFrom),
	}
	pages, err := s.ListSessionsBy(&filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	var sessions []*smmodels.ArkSMSession
	for page := range pages {
		sessions = append(sessions, page.Items...)
	}

	stats := &smmodels.ArkSMSessionsStats{}
	stats.SessionsCount = len(sessions)
	stats.SessionsFailureCount = 0
	stats.SessionsCountPerApplicationCode = make(map[string]int)
	stats.SessionsCountPerPlatform = make(map[string]int)
	stats.SessionsCountPerProtocol = make(map[string]int)
	stats.SessionsCountPerStatus = make(map[smmodels.ArkSMSessionStatus]int)

	for _, session := range sessions {
		if session.SessionStatus == smmodels.Failed {
			stats.SessionsFailureCount++
		}
		stats.SessionsCountPerApplicationCode[session.ApplicationCode]++
		stats.SessionsCountPerPlatform[session.Platform]++
		stats.SessionsCountPerProtocol[session.Protocol]++
		stats.SessionsCountPerStatus[session.SessionStatus]++
	}

	return stats, nil
}

// Session retrieves a session by its ID
func (s *ArkSMService) Session(getSession *smmodels.ArkSIASMGetSession) (*smmodels.ArkSMSession, error) {
	s.Logger.Info("Getting session [%s]", getSession.SessionID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(sessionURL, getSession.SessionID), nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get session - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	sessionJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var session smmodels.ArkSMSession
	err = mapstructure.Decode(sessionJSON, &session)
	return &session, nil
}

// ServiceConfig returns the service configuration for the ArkSMservice.
func (s *ArkSMService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
