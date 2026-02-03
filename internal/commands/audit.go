package commands

import (
	"encoding/json"
	"time"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	auditQuery string
	auditFrom  int64
	auditTo    int64
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit service",
}

var auditDbStatusCmd = &cobra.Command{
	Use:   "db_status",
	Short: "Get DB status",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.DbStatus()
		if err != nil {
			return err
		}
		// etcd overall status and nodes (name: active)
		if resp.EtcdStatus != nil {
			cmd.Printf("etcd status: %s\n", resp.EtcdStatus.Status)
			for _, n := range resp.EtcdStatus.Nodes {
				if n != nil {
					active := "inactive"
					if n.Active {
						active = "active"
					}
					cmd.Printf("  %s: %s\n", n.Name, active)
				}
			}
		}
		// master status with leader
		if resp.Master != nil {
			cmd.Printf("master status: %s, leader: %s\n", resp.Master.Status, resp.Master.Leader)
		}
		// indexes: indexname: queue.length
		if len(resp.Indexes) > 0 {
			for _, idx := range resp.Indexes {
				if idx != nil {
					length := uint(0)
					if idx.Queue != nil {
						length = idx.Queue.Length
					}
					cmd.Printf("index: %s, length: %d\n", idx.IndexName, length)
				}
			}
		}
		return nil
	},
}

var auditListCmd = &cobra.Command{
	Use:   "list",
	Short: "List audit events",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := auditFrom, auditTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}
		resp, err := AppAPI.AuditSearch(auditQuery, from, to)
		if err != nil {
			return err
		}
		if resp.Hits == nil || len(resp.Hits.Hits) == 0 {
			cmd.PrintErrln("No hits found.")
			return nil
		}
		for _, hit := range resp.Hits.Hits {
			var src model.AuditEvent
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.Id, err)
				continue
			}
			createdAt := time.UnixMilli(src.CreatedOn).Format(time.RFC3339)
			cmd.Printf("Username: %s, Action: %s, Time: %s\n", src.Username, src.Action, createdAt)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditListCmd, auditDbStatusCmd)

	auditListCmd.Flags().StringVar(&auditQuery, "query", "", "Search query")
	auditListCmd.Flags().Int64Var(&auditFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	auditListCmd.Flags().Int64Var(&auditTo, "to", 0, "Range end (Unix ms); default: now")
}
