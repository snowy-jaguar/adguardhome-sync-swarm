package sync

import (
	"github.com/bakito/adguardhome-sync/pkg/client"
	"github.com/bakito/adguardhome-sync/pkg/client/model"
	"github.com/bakito/adguardhome-sync/pkg/utils"
	"go.uber.org/zap"
)

var (
	actionProfileInfo = func(ac *actionContext) (bool, error) {
		if pro, err := ac.client.ProfileInfo(); err != nil {
			return false, err
		} else if merged := pro.ShouldSyncFor(ac.origin.profileInfo, ac.cfg.Features.Theme); merged != nil {
			return true, ac.client.SetProfileInfo(merged)
		}
		return false, nil
	}
	actionProtection = func(ac *actionContext) (bool, error) {
		if ac.origin.status.ProtectionEnabled != ac.replicaStatus.ProtectionEnabled {
			return true, ac.client.ToggleProtection(ac.origin.status.ProtectionEnabled)
		}
		return false, nil
	}
	actionParental = func(ac *actionContext) (bool, error) {
		if rp, err := ac.client.Parental(); err != nil {
			return false, err
		} else if ac.origin.parental != rp {
			return true, ac.client.ToggleParental(ac.origin.parental)
		}
		return false, nil
	}
	actionSafeSearchConfig = func(ac *actionContext) (bool, error) {
		if ssc, err := ac.client.SafeSearchConfig(); err != nil {
			return false, err
		} else if !ac.origin.safeSearch.Equals(ssc) {
			return true, ac.client.SetSafeSearchConfig(ac.origin.safeSearch)
		}
		return false, nil
	}
	actionSafeBrowsing = func(ac *actionContext) (bool, error) {
		if rs, err := ac.client.SafeBrowsing(); err != nil {
			return false, err
		} else if ac.origin.safeBrowsing != rs {
			if err = ac.client.ToggleSafeBrowsing(ac.origin.safeBrowsing); err != nil {
				return false, err
			}
		}
		return false, nil
	}
	actionQueryLogConfig = func(ac *actionContext) (bool, error) {
		qlc, err := ac.client.QueryLogConfig()
		if err != nil {
			return false, err
		}
		if !ac.origin.queryLogConfig.Equals(qlc) {
			return true, ac.client.SetQueryLogConfig(ac.origin.queryLogConfig)
		}
		return false, nil
	}
	actionStatsConfig = func(ac *actionContext) (bool, error) {
		sc, err := ac.client.StatsConfig()
		if err != nil {
			return false, err
		}
		if !sc.Equals(ac.origin.statsConfig) {
			return true, ac.client.SetStatsConfig(ac.origin.statsConfig)
		}
		return false, nil
	}
	actionDNSRewrites = func(ac *actionContext) (bool, error) {
		replicaRewrites, err := ac.client.RewriteList()
		if err != nil {
			return false, err
		}

		a, r, d := replicaRewrites.Merge(ac.origin.rewrites)

		if err = ac.client.DeleteRewriteEntries(r...); err != nil {
			return false, err
		}
		if err = ac.client.AddRewriteEntries(a...); err != nil {
			return false, err
		}

		for _, dupl := range d {
			ac.rl.With("domain", dupl.Domain, "answer", dupl.Answer).Warn("Skipping duplicated rewrite from source")
		}
		return false, nil
	}
	actionFilters = func(ac *actionContext) (bool, error) {
		rf, err := ac.client.Filtering()
		if err != nil {
			return false, err
		}

		if err = syncFilterType(ac.rl, ac.origin.filters.Filters, rf.Filters, false, ac.client, ac.cfg.ContinueOnError); err != nil {
			return false, err
		}
		if err = syncFilterType(ac.rl, ac.origin.filters.WhitelistFilters, rf.WhitelistFilters, true, ac.client, ac.cfg.ContinueOnError); err != nil {
			return false, err
		}

		changed := false
		if utils.PtrToString(ac.origin.filters.UserRules) != utils.PtrToString(rf.UserRules) {
			err := ac.client.SetCustomRules(ac.origin.filters.UserRules)
			if err != nil {
				return false, err
			}
			changed = true
		}

		if !utils.PtrEquals(ac.origin.filters.Enabled, rf.Enabled) ||
			!utils.PtrEquals(ac.origin.filters.Interval, rf.Interval) {
			err = ac.client.ToggleFiltering(*ac.origin.filters.Enabled, *ac.origin.filters.Interval)
			if err != nil {
				return false, err
			}
			changed = true
		}
		return changed, nil
	}

	actionBlockedServicesSchedule = func(ac *actionContext) (bool, error) {
		rbss, err := ac.client.BlockedServicesSchedule()
		if err != nil {
			return false, err
		}

		if !ac.origin.blockedServicesSchedule.Equals(rbss) {
			return true, ac.client.SetBlockedServicesSchedule(ac.origin.blockedServicesSchedule)
		}
		return false, nil
	}
	actionClientSettings = func(ac *actionContext) (bool, error) {
		rc, err := ac.client.Clients()
		if err != nil {
			return false, err
		}

		a, u, r := rc.Merge(ac.origin.clients)

		for _, client := range r {
			if err := ac.client.DeleteClient(client); err != nil {
				ac.rl.With("client-name", client.Name, "error", err).Error("error deleting client setting")
				if !ac.cfg.ContinueOnError {
					return false, err
				}
			}
		}

		for _, client := range a {
			if err := ac.client.AddClient(client); err != nil {
				ac.rl.With("client-name", client.Name, "error", err).Error("error adding client setting")
				if !ac.cfg.ContinueOnError {
					return false, err
				}
			}
		}

		for _, client := range u {
			if err := ac.client.UpdateClient(client); err != nil {
				ac.rl.With("client-name", client.Name, "error", err).Error("error updating client setting")
				if !ac.cfg.ContinueOnError {
					return false, err
				}
			}
		}

		return false, nil
	}

	actionDNSAccessLists = func(ac *actionContext) (bool, error) {
		al, err := ac.client.AccessList()
		if err != nil {
			return false, err
		}
		if !al.Equals(ac.origin.accessList) {
			return true, ac.client.SetAccessList(ac.origin.accessList)
		}
		return false, nil
	}
	actionDNSServerConfig = func(ac *actionContext) (bool, error) {
		dc, err := ac.client.DNSConfig()
		if err != nil {
			return false, err
		}

		// dc.Sanitize(ac.rl)

		if !dc.Equals(ac.origin.dnsConfig) {
			if err = ac.client.SetDNSConfig(ac.origin.dnsConfig); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
	actionDHCPServerConfig = func(ac *actionContext) (bool, error) {
		if ac.origin.dhcpServerConfig.HasConfig() {
			sc, err := ac.client.DhcpConfig()
			if err != nil {
				return false, err
			}
			origClone := ac.origin.dhcpServerConfig.Clone()
			if ac.replica.InterfaceName != "" {
				// overwrite interface name
				origClone.InterfaceName = utils.Ptr(ac.replica.InterfaceName)
			}
			if ac.replica.DHCPServerEnabled != nil {
				// overwrite dhcp enabled
				origClone.Enabled = ac.replica.DHCPServerEnabled
			}

			if !sc.CleanAndEquals(origClone) {
				return true, ac.client.SetDhcpConfig(origClone)
			}
		}
		return false, nil
	}
	actionDHCPStaticLeases = func(ac *actionContext) (bool, error) {
		sc, err := ac.client.DhcpConfig()
		if err != nil {
			return false, err
		}

		a, r := model.MergeDhcpStaticLeases(sc.StaticLeases, ac.origin.dhcpServerConfig.StaticLeases)

		for _, lease := range r {
			if err := ac.client.DeleteDHCPStaticLease(lease); err != nil {
				ac.rl.With("hostname", lease.Hostname, "error", err).Error("error deleting dhcp static lease")
				if !ac.cfg.ContinueOnError {
					return false, err
				}
			}
		}

		for _, lease := range a {
			if err := ac.client.AddDHCPStaticLease(lease); err != nil {
				ac.rl.With("hostname", lease.Hostname, "error", err).Error("error adding dhcp static lease")
				if !ac.cfg.ContinueOnError {
					return false, err
				}
			}
		}
		return false, nil
	}
)

func syncFilterType(
	rl *zap.SugaredLogger,
	of *[]model.Filter,
	rFilters *[]model.Filter,
	whitelist bool,
	replica client.Client,
	continueOnError bool,
) error {
	fa, fu, fd := model.MergeFilters(rFilters, of)

	for _, f := range fd {
		if err := replica.DeleteFilter(whitelist, f); err != nil {
			rl.With("filter", f.Name, "url", f.Url, "whitelist", whitelist, "error", err).Error("error deleting filter")
			if !continueOnError {
				return err
			}
		}
	}

	for _, f := range fa {
		if err := replica.AddFilter(whitelist, f); err != nil {
			rl.With("filter", f.Name, "url", f.Url, "whitelist", whitelist, "error", err).Error("error adding filter")
			if !continueOnError {
				return err
			}
		}
	}

	for _, f := range fu {
		if err := replica.UpdateFilter(whitelist, f); err != nil {
			rl.With("filter", f.Name, "url", f.Url, "whitelist", whitelist, "error", err).Error("error updating filter")
			if !continueOnError {
				return err
			}
		}
	}

	if len(fa) > 0 || len(fu) > 0 {
		if err := replica.RefreshFilters(whitelist); err != nil {
			return err
		}
	}
	return nil
}
