package healthcheck

import (
	"time"
	"math/rand"
	"github.com/rcrowley/go-metrics"
	"gitlab.alipay-inc.com/afe/mosn/pkg/types"
	"gitlab.alipay-inc.com/afe/mosn/pkg/api/v2"
)

// types.HealthChecker
type healthChecker struct {
	serviceName         string
	healthCheckCbs      []types.HealthCheckCb
	cluster             types.Cluster
	healthCheckSessions map[types.Host]types.HealthCheckSession

	timeout        time.Duration
	interval       time.Duration
	intervalJitter time.Duration

	localProcessHealthy uint64

	healthyThreshold   uint32
	unhealthyThreshold uint32

	stats *healthCheckStats
}

func NewHealthCheck(config v2.HealthCheck) *healthChecker {
	hc := &healthChecker{
		healthCheckSessions: make(map[types.Host]types.HealthCheckSession),
		timeout:             config.Timeout,
		interval:            config.Interval,
		intervalJitter:      config.IntervalJitter,
		healthyThreshold:    config.HealthyThreshold,
		unhealthyThreshold:  config.UnhealthyThreshold,
	}

	return hc
}

func (c *healthChecker) Start() {
	for _, hostSet := range c.cluster.PrioritySet().HostSetsByPriority() {
		c.addHosts(hostSet.Hosts())
	}
}

func (c *healthChecker) Stop() {
	// todo
}

func (c *healthChecker) AddHostCheckCompleteCb(cb types.HealthCheckCb) {
	c.healthCheckCbs = append(c.healthCheckCbs, cb)
}

func (c *healthChecker) newSession(host types.Host) types.HealthCheckSession {
	return nil
}

func (c *healthChecker) addHosts(hosts []types.Host) {
	for _, host := range hosts {
		c.healthCheckSessions[host] = c.newSession(host)
		c.healthCheckSessions[host].Start()
	}
}

func (c *healthChecker) decHealthy() {
	c.localProcessHealthy--
	c.refreshHealthyStat()
}

func (c *healthChecker) incHealthy() {
	c.localProcessHealthy++
	c.refreshHealthyStat()
}

func (c *healthChecker) refreshHealthyStat() {
	c.stats.healthy.Update(int64(c.localProcessHealthy))
}

func (c *healthChecker) getStats() *healthCheckStats {
	return c.stats
}

func (c *healthChecker) getInterval() time.Duration {
	baseInterval := c.interval

	if c.intervalJitter > 0 {
		jitter := int(rand.Float32() * float32(c.intervalJitter))
		baseInterval += time.Duration(jitter)
	}

	if baseInterval < 0 {
		baseInterval = 0
	}

	maxUint := ^uint(0)
	if uint(baseInterval) > maxUint {
		baseInterval = time.Duration(maxUint)
	}

	return baseInterval
}

func (c *healthChecker) runCallbacks(host types.Host, changed bool) {
	c.refreshHealthyStat()

	for _, cb := range c.healthCheckCbs {
		cb(host, changed)
	}
}

type healthCheckSession struct {
	healthChecker *healthChecker
	intervalTimer *timer
	timeoutTimer  *timer
	numHealthy    uint32
	numUnHealthy  uint32
	host          types.Host
}

func NewHealthCheckSession(hc *healthChecker, host types.Host) *healthCheckSession {
	hcs := &healthCheckSession{
		healthChecker: hc,
		host:          host,
	}

	hcs.intervalTimer = newTimer(hcs.onIntervalBase)
	hcs.timeoutTimer = newTimer(hcs.onTimeoutBase)

	if !host.ContainHealthFlag(types.FAILED_ACTIVE_HC) {
		hcs.healthChecker.decHealthy()
	}

	return hcs
}

func (s *healthCheckSession) Start() {
	s.onIntervalBase()
}

func (s *healthCheckSession) Stop() {}

func (s *healthCheckSession) handleSuccess() {
	s.numUnHealthy = 0

	stateChanged := false
	if s.host.ContainHealthFlag(types.FAILED_ACTIVE_HC) {
		s.numHealthy++

		if s.numHealthy == s.healthChecker.healthyThreshold {
			s.healthChecker.incHealthy()
			stateChanged = true
		}
	}

	s.healthChecker.stats.success.Inc(1)
	s.healthChecker.runCallbacks(s.host, stateChanged)

	s.timeoutTimer.stop()
	s.intervalTimer.start(s.healthChecker.getInterval())
}

func (s *healthCheckSession) SetUnhealthy(fType types.FailureType) {
	s.numHealthy = 0

	stateChanged := false
	if !s.host.ContainHealthFlag(types.FAILED_ACTIVE_HC) {
		s.numUnHealthy++

		if s.numUnHealthy == s.healthChecker.unhealthyThreshold {
			s.host.SetHealthFlag(types.FAILED_ACTIVE_HC)
			s.healthChecker.decHealthy()
			stateChanged = true
		}
	}

	s.healthChecker.stats.failure.Inc(1)

	switch fType {
	case types.FailureNetwork:
		s.healthChecker.stats.networkFailure.Inc(1)
	case types.FailurePassive:
		s.healthChecker.stats.passiveFailure.Inc(1)
	}

	s.healthChecker.runCallbacks(s.host, stateChanged)
}

func (s *healthCheckSession) handleFailure(fType types.FailureType) {
	s.SetUnhealthy(fType)

	s.timeoutTimer.stop()
	s.intervalTimer.start(s.healthChecker.getInterval())
}

func (s *healthCheckSession) onIntervalBase() {
	s.onInterval()

	s.timeoutTimer.start(s.healthChecker.getInterval())
	s.healthChecker.stats.attempt.Inc(1)
}

func (s *healthCheckSession) onInterval() {}

func (s *healthCheckSession) onTimeoutBase() {
	s.onTimeout()
	s.SetUnhealthy(types.FailureNetwork)
}

func (s *healthCheckSession) onTimeout() {}

type healthCheckStats struct {
	attempt        metrics.Counter
	success        metrics.Counter
	failure        metrics.Counter
	passiveFailure metrics.Counter
	networkFailure metrics.Counter
	verifyCluster  metrics.Counter
	healthy        metrics.Gauge
}