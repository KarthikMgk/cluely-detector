package main

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"time"

	"watchdog/overlaydetector"
)

type reply struct {
	Overlay bool   `json:"overlay"`
	Reason  string `json:"reason,omitempty"`
	Ts      int64  `json:"ts"`
}

// func main() {
// 	for {
// 		ov, why := overlaydetector.Scan()

// 		out, _ := json.Marshal(reply{ov, why, time.Now().Unix()})

// 		binary.Write(os.Stdout, binary.LittleEndian, uint32(len(out)))
// 		os.Stdout.Write(out)            // ← write raw JSON ONLY (no newline)

// 		time.Sleep(500 * time.Millisecond)
// 	}
// }

// func main() {
// 	for {
// 		ov, why := overlaydetector.Scan()

// 		// marshal ONE json blob
// 		payload, _ := json.Marshal(reply{ov, why, time.Now().Unix()})

// 		// ---- Chrome framing: 4-byte length + raw JSON ----
// 		binary.Write(os.Stdout, binary.LittleEndian, uint32(len(payload)))
// 		os.Stdout.Write(payload)        // <-- write the bytes ONLY

// 		time.Sleep(500 * time.Millisecond)
// 	}
// }


func main() {
    for {
        ov, why := overlaydetector.Scan()
        if !ov {                               // nothing suspicious → skip
            time.Sleep(500 * time.Millisecond)
            continue
        }

        payload, _ := json.Marshal(reply{ov, why, time.Now().Unix()})
        binary.Write(os.Stdout, binary.LittleEndian, uint32(len(payload)))
        os.Stdout.Write(payload)               // send only the positive hit

        time.Sleep(500 * time.Millisecond)     // keep loop cadence
    }
}
