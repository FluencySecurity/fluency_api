package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	sourceType      string
	sinkType        string
	resourceID      string
	resourceName    string
	dataFormat      string
	dataCompression string
	integrationID   string // For associating with an integration
	url             string // For HEC sink
	token           string // For HEC sink

	processorName string
	routerName    string

	sourceID string // For connecting source to router
	routerID string // For connecting source to router
	sinkID   string // For connecting source to router

)

// formatSlotTime formats a slot (Unix milliseconds) as a readable time string.
func formatSlotTime(slotMs int64) string {
	return time.UnixMilli(slotMs).Format("2006-01-02 15:04:05")
}

// formatIntWithCommas formats n with thousand separators (e.g. 20785008 -> "20,785,008").
func formatIntWithCommas(n int64) string {
	s := strconv.FormatInt(n, 10)
	if n < 0 {
		s = s[1:]
	}
	digits := len(s)
	if digits <= 3 {
		if n < 0 {
			return "-" + s
		}
		return s
	}
	var b strings.Builder
	first := digits % 3
	if first == 0 {
		first = 3
	}
	if n < 0 {
		b.WriteByte('-')
	}
	b.WriteString(s[:first])
	for i := first; i < digits; i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

// formatFloatWithCommas formats f with thousand separators; whole numbers like 20,785,008, decimals like 0.063.
func formatFloatWithCommas(f float64) string {
	intPart := int64(f)
	if f == float64(intPart) {
		return formatIntWithCommas(intPart)
	}
	frac := f - float64(intPart)
	if frac < 0 {
		frac = -frac
	}
	fracStr := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(frac, 'f', 6, 64), "0"), ".")
	if strings.HasPrefix(fracStr, "0.") {
		fracStr = "." + fracStr[2:]
	}
	out := formatIntWithCommas(intPart) + fracStr
	if f < 0 && intPart == 0 {
		out = "-" + out
	}
	return out
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Manage streams",
}

var connectSinkCmd = &cobra.Command{
	Use:   "connect-sink",
	Short: "Connect a stream router to a sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Connecting stream router to a sink...")

		err := AppAPI.SetRouterSink(routerID, sinkID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream connection added successfully: router ", routerID, " to sink ", sinkID)
		//cmd.Println(resourceID)
		return nil
	},
}

var connectRouterCmd = &cobra.Command{
	Use:   "connect-router",
	Short: "Connect a stream source to a router",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Connecting source to a stream router...")

		err := AppAPI.SetSourceRouter(sourceID, routerID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream connection added successfully: source ", sourceID, " to router ", routerID)
		//cmd.Println(resourceID)
		return nil
	},
}

var addRouterCmd = &cobra.Command{
	Use:   "add-router",
	Short: "Add a stream router",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream router...")

		resourceID, err := AppAPI.AddSimpleRouter(processorName, routerName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream router added successfully: ", resourceID)
		fmt.Println(resourceID)
		return nil
	},
}

var addSourceCmd = &cobra.Command{
	Use:   "add-source",
	Short: "Add a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream datasource...")

		source := &model.DataSourceConfig{
			Type: sourceType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sourceType == "plugin" {
			if integrationID == "" {
				return fmt.Errorf("integration-id is required for plugin source type")
			}
			source.Plugin = &model.PluginSourceConfig{
				ID: integrationID,
			}
		} else if sourceType == "hec" {
			source.Hec = &model.HecSourceConfig{
				URL: url,
			}
			s := &model.HecSecret{
				Token: token,
			}
			b, _ := json.Marshal(s)
			source.Secret = b
		}

		response, err := AppAPI.AddDataSource(source)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream source added successfully: ", response.ID)
		if response.URL != "" {
			cmd.PrintErrln("Access URL:", response.URL)
		}
		if len(response.Secret) > 0 {
			b, _ := response.Secret.MarshalJSON()
			cmd.PrintErrln("Secret:", string(b))
		}
		fmt.Println(response.ID)
		return nil
	},
}

var delSourceCmd = &cobra.Command{
	Use:   "del-source",
	Short: "Delete a stream source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasource...")

		err := AppAPI.DeleteDataSource(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream source deleted successfully: ", resourceID)
		return nil
	},
}

var listSourceCmd = &cobra.Command{
	Use:   "list-source",
	Short: "List all stream sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasource...")

		entries, err := AppAPI.ListDataSource()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream source found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
		return nil
	},
}

var delSinkCmd = &cobra.Command{
	Use:   "del-sink",
	Short: "Delete a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Deleting stream datasink...")

		err := AppAPI.DeleteDataSink(resourceID)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Stream sink deleted successfully: ", resourceID)
		return nil
	},
}

var listSinkCmd = &cobra.Command{
	Use:   "list-sink",
	Short: "List all stream sinks",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing stream datasink...")

		entries, err := AppAPI.ListDataSink()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No stream sink found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Type: %s\n", entry.ID, entry.Name, entry.Type)
		}
		return nil
	},
}

// Example leaf command: source
var addSinkCmd = &cobra.Command{
	Use:   "add-sink",
	Short: "Add a stream sink",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Adding stream sink...")
		// Just call the global interface
		sink := &model.DataSinkConfig{
			Type: sinkType,
			Name: resourceName,
			// Add other necessary fields for DataSourceConfig
			//Format: "json", // Example format, adjust as needed
			//Compression
		}

		if sinkType == "datalake" {
			if index == "" {
				return fmt.Errorf("index is required for datalake sink type")
			}
			sink.DataLake = &model.DataLakeSinkConfig{
				Datalake:      datalake,
				DatalakeIndex: index,
			}
		} else if sinkType == "hec" {
			sink.Hec = &model.HecSinkConfig{
				URL:   url,
				Token: token,
				// Add necessary fields for HEC sink
			}
		} else if sinkType == "webhook" {
			sink.Webhook = &model.WebhookSinkConfig{
				URL: url,
				//Token: token,
				// Add necessary fields for HEC sink
			}
		}

		response, err := AppAPI.AddDataSink(sink)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Stream sink added successfully: ", response.ID)
		fmt.Println(response.ID)

		return nil
	},
}

var streamStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get platform status",
	Long:  "Calls the platform ListConfigs API and prints metrics (id, action; slots: values).",
	RunE: func(cmd *cobra.Command, args []string) error {
		metrics, err := AppAPI.PlatformStatus()
		if err != nil {
			return err
		}
		printed := 0
		for _, m := range metrics {
			if m == nil {
				continue
			}
			comp := strings.ToLower(strings.TrimSpace(m.Component))
			action := strings.ToLower(strings.TrimSpace(m.Action))
			switch comp {
			case "datasource", "source":
				if action != "input" {
					continue
				}
			case "datasink", "sink":
				if action != "output" {
					continue
				}
			default:
				continue
			}
			if printed > 0 {
				cmd.Println()
			}
			printed++
			cmd.Printf("%s, %s\n", m.ID, m.Action)
			slots := m.Slots
			values := m.Values
			n := len(slots)
			if len(values) < n {
				n = len(values)
			}
			for j := 0; j < n; j++ {
				cmd.Printf("%s - %s B\n", formatSlotTime(slots[j]), formatFloatWithCommas(values[j]))
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.AddCommand(addSourceCmd, delSourceCmd, listSourceCmd, addSinkCmd, delSinkCmd, listSinkCmd, addRouterCmd, connectRouterCmd, connectSinkCmd, streamStatusCmd)

	addSourceCmd.Flags().StringVar(&sourceType, "source-type", "", "data source type: plugin, s3, hec, webhook ")
	addSourceCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSourceCmd.Flags().StringVar(&dataFormat, "format", "json", "Data Format")
	addSourceCmd.Flags().StringVar(&dataCompression, "compression", "", "Data Compression")

	addSourceCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSourceCmd.MarkFlagRequired("source-type")
	_ = addSourceCmd.MarkFlagRequired("name")

	addSinkCmd.Flags().StringVar(&sinkType, "sink-type", "", "data sink type: datalake, hec, webhook, drop")
	addSinkCmd.Flags().StringVar(&resourceName, "name", "", "Name")
	addSinkCmd.Flags().StringVar(&url, "url", "", "URL for HEC sink")
	addSinkCmd.Flags().StringVar(&token, "token", "", "Token")

	addSinkCmd.Flags().StringVar(&datalake, "datalake", "managed", "datalake name")
	addSinkCmd.Flags().StringVar(&index, "index", "", "datalake index name")

	//addSinkCmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID")

	_ = addSinkCmd.MarkFlagRequired("sink-type")
	_ = addSinkCmd.MarkFlagRequired("name")

	delSourceCmd.Flags().StringVar(&resourceID, "id", "", "data source ID")
	_ = delSourceCmd.MarkFlagRequired("id")

	delSinkCmd.Flags().StringVar(&resourceID, "id", "", "data sink ID")
	_ = delSinkCmd.MarkFlagRequired("id")

	addRouterCmd.Flags().StringVar(&processorName, "processor", "", "processor name")
	addRouterCmd.Flags().StringVar(&routerName, "router-name", "", "Router Name")

	_ = addRouterCmd.MarkFlagRequired("processor")

	connectRouterCmd.Flags().StringVar(&sourceID, "source-id", "", "source ID")
	connectRouterCmd.Flags().StringVar(&routerID, "router-id", "", "Router ID")

	_ = connectRouterCmd.MarkFlagRequired("source-id")
	_ = connectRouterCmd.MarkFlagRequired("router-id")

	connectSinkCmd.Flags().StringVar(&sinkID, "sink-id", "", "sink ID")
	connectSinkCmd.Flags().StringVar(&routerID, "router-id", "", "Router ID")

	_ = connectSinkCmd.MarkFlagRequired("sink-id")
	_ = connectSinkCmd.MarkFlagRequired("router-id")
}
