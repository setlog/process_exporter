package flags

import "flag"

func Parse(arguments []string) (namespace, binary, argName string, port int, interval float64) {
	flagSet := flag.NewFlagSet("flags", flag.ExitOnError)
	namespacePtr := flagSet.String("namespace", "my", "Prometheus metric namespace.")
	binaryPtr := flagSet.String("binary", "", "Filter which processes to watch by binary name. This is limited to 15 bytes because of the kernel.")
	nameFlagPtr := flagSet.String("nameflag", "name", "Set Prometheus \"name\"-label value to value of this command line argument of the monitored processes.")
	portPtr := flagSet.Int("port", 80, "Port on which to listen to requests to /metrics.")
	intervalPtr := flagSet.Float64("interval", 10, "Interval between refreshes of metrics, in seconds. Should not be too large to prevent CPU reading from getting noisy.")
	err := flagSet.Parse(arguments)
	if err != nil {
		panic(err)
	}
	return *namespacePtr, *binaryPtr, *nameFlagPtr, *portPtr, *intervalPtr
}
