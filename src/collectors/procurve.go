package collectors

import (
"time"
"regexp"
"fmt"
"contra/src/utils"
)

func CollectProcurve() string {
    fmt.Printf("Collect Works - Procurve\n")

    // set up ssh connection - obviously not the right place for this
    connection, err := utils.SshClient("changeme", "wrongplace", "10.0.0.0:22")

    if err != nil {
        panic(err)
    }

    // Output we expect to receive
    receive := map[int]string {
        1 : "continue", // 1 : should always match the initial connection string
        2 : ".*#",
        3 : ".*#",
        4 : ".*#",
    }

    // Commands we will send in response to output above
    send := map[int]string {
        1 : "a\n",
        2 : "no page\n",
        3 : "show running-config\n",
    }

    // Build batcher
    batch := utils.BuildBatcher(send, receive)

    // call GatherExpect to collect the configs
    result, err := utils.GatherExpect(batch, time.Second*10, connection)
    if err != nil {
        panic(err)
    }

    // Strip shell commands, grab only the xml file
    // this regex assumes all procurve configs begin with 'hostname', and end with 'password manager'
    // Should probably find a better match...
    config := regexp.MustCompile(`hostname[\s\S]*?manager`)
    // search the last element of result for the regex above
    match := config.FindStringSubmatch(result[3].Output)

    utils.WriteFile(match[0], "procurve.txt")
    return match[0]
}